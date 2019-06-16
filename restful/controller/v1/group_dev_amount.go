package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"fmt"
	"net/http"
)

type GroupDeviceAmountHandler struct {
}

func NewGroupDeviceAmountHandler() GroupDeviceAmountHandler {
	return GroupDeviceAmountHandler{}
}

const (
	GroupDeviceAmountAPI = "/group/dev/amount"
)

func (h GroupDeviceAmountHandler) SetRouter(r *Router) {
	r.POST(GroupDeviceAmountAPI, PostGroupDeviceAmount, mds...)
}

type GroupDeviceAmount struct {
	Group  uint64 `json:"group_id"`
	Model  string `json:"dev_model"`
	Amount uint64 `json:"amount"`
}

func PostGroupDeviceAmount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}
	if u.Info().Type != "Noovo" && u.Info().Type != "SI" {
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	var da GroupDeviceAmount
	if err := json.NewDecoder(r.Body).Decode(&da); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	gps, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if !util.CheckGroup(gps, da.Group) {
		response.ErrInsufficientPrivilegeOnGroup(da.Group).Send(w)
		return
	}

	db := u.DB()
	var cid uint64
	if err := db.QueryRow(`SELECT com FROM groups WHERE id=$1`, da.Group).Scan(&cid); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	rows, err := db.Query(fmt.Sprintf(view.SelectDeviceGroupSQL, "WHERE com=$1 AND gp IS null AND model=$2 LIMIT $3"), cid, da.Model, da.Amount)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	dgs, err := view.ScanDeviceGroups(rows)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if has, want := len(dgs), int(da.Amount); has < want {
		response.ErrNotEnoughDevice(has, want).Send(w)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	var devs []models.Device
	for _, dg := range dgs {
		if err := models.InsertDeviceGroup(tx, models.DeviceGroup{SN: dg.SN, Group: da.Group}); err != nil {
			defer tx.Rollback()
			response.ErrRaw(err).Send(w)
			return
		}
		devs = append(devs, dg.Device)
	}
	if err := tx.Commit(); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(devs).Send(w)

}
