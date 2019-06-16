package v1

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"fmt"
	"net/http"

	"time"
)

type ComDeviceAmountHandler struct {
}

func NewComDeviceAmountHandler() ComDeviceAmountHandler {
	return ComDeviceAmountHandler{}
}

const (
	CompanyDeviceAmountAPI = "/com/dev/amount"
)

func (h ComDeviceAmountHandler) SetRouter(r *Router) {
	r.GET(CompanyDeviceAmountAPI, GetCompanyDeviceAmount, mds...)
	r.POST(CompanyDeviceAmountAPI, PostCompanyDeviceAmount, mds...)
}

type DeviceAmount struct {
	Company uint64 `json:"com_id"`
	Model   string `json:"dev_model"`
	Amount  uint64 `json:"amount"`
}

func GetCompanyDeviceAmount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Type != "Noovo" && u.Info().Type != "SI" {
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	// var cid uint64
	// com := r.URL.Query().Get("com")
	// if com == "" {
	// 	cid = u.Info().Company

	// } else {
	// 	id, err := strconv.Atoi(com)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		response.ErrInvalidQueryParameter("com", com).Send(w)
	// 		return
	// 	}

	// 	cid = uint64(id)
	// 	if cid != u.Info().Company {
	// 		if u.Info().Type != "Noovo" || u.Info().Company != 1 {
	// 			response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
	// 			return
	// 		}
	// 	}
	// }

	cid := u.Info().Company

	das, err := view.DBDeviceAmounts(u.DB(), cid)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(das).Send(w)

}

func PostCompanyDeviceAmount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Type != "Noovo" || u.Info().Company != 1 {
		response.ErrNotNoovo("post company device amount").Send(w)
		return
	}

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}

	var da DeviceAmount
	if err := json.NewDecoder(r.Body).Decode(&da); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	db := u.DB()
	rows, err := db.Query(fmt.Sprintf(view.SelectDeviceGroupSQL, "WHERE com=1 AND gp=1 AND model=$1 LIMIT $2"), da.Model, da.Amount)
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
		if err := models.DeleteDevicesGroup(tx, dg.SN); err != nil {
			defer tx.Rollback()
			response.ErrRaw(err).Send(w)
			return
		}
		if err := models.UpdateDeviceCompany(tx, dg.SN, da.Company); err != nil {
			defer tx.Rollback()
			response.ErrRaw(err).Send(w)
			return
		}
		n, err := models.CountDeviceStatusBySn(tx, dg.SN)
		if err != nil {
			defer tx.Rollback()
			response.ErrRaw(err).Send(w)
			return
		}
		now := null.NewTime(time.Now().UTC())
		if n == 0 {
			if err := models.InsertDeviceStatus(tx, models.DeviceStatus{SN: dg.SN, WarrantyDate: now}); err != nil {
				defer tx.Rollback()
				response.ErrRaw(err).Send(w)
				return
			}
		} else {
			ds, err := models.DeviceStatusBySn(tx, dg.SN)
			if err != nil {
				defer tx.Rollback()
				response.ErrRaw(err).Send(w)
				return
			}
			// if WarrantyDate is null, update it
			if !ds.WarrantyDate.Valid {
				ds.WarrantyDate = now
				if err := models.UpdateDeviceStatus(tx, ds); err != nil {
					defer tx.Rollback()
					response.ErrRaw(err).Send(w)
					return
				}
			}
		}

		dg.Device.Company = da.Company
		devs = append(devs, dg.Device)
	}
	if err := tx.Commit(); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(devs).Send(w)

}
