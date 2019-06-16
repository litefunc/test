package v1

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"cloud/protos/pb"
	"cloud/server/portal/models"
	"cloud/server/portal/proto"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"net/http"

	"time"
)

type DeviceGroupHandler struct {
}

func NewDeviceGroupHandler() DeviceGroupHandler {
	return DeviceGroupHandler{}
}

const (
	DeviceGroupAPI = "/dev/group"
)

func (h DeviceGroupHandler) SetRouter(r *Router) {
	r.PUT(DeviceGroupAPI, PutDeviceGroup, mds...)
}

type NewDevsGroup struct {
	Gp  uint64   `json:"new_group"`
	Sns []string `json:"sns"`
}

func PutDeviceGroup(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	var dg NewDevsGroup
	if err := json.NewDecoder(r.Body).Decode(&dg); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	if dg.Gp == 0 {
		response.ErrInvalidParameter("new_group", dg.Gp).Send(w)
		return
	}

	gps, err := u.Groups()
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if !util.CheckGroup(gps, dg.Gp) {
		response.ErrInsufficientPrivilegeOnGroup(dg.Gp).Send(w)
		return
	}

	devs, err := u.Devices()
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if invalid := util.CheckDevices(devs, dg.Sns); len(invalid) != 0 {
		response.ErrInsufficientPrivilegeOnDevices(invalid...).Send(w)
		return
	}

	affected, err := u.Ship(dg.Gp, dg.Sns...)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	com, err := models.CompanyIDByGroup(u.DB(), dg.Gp)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	for _, sn := range affected {
		if err := models.DbUpdateDeviceCompany(u.DB(), sn, com); err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	tx, err := u.DB().Begin()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if u.Info().Company == 1 && com != 1 {

		for _, sn := range affected {
			n, err := models.CountDeviceStatusBySn(tx, sn)
			if err != nil {
				defer tx.Rollback()
				response.ErrRaw(err).Send(w)
				return
			}
			now := null.NewTime(time.Now().UTC())
			if n == 0 {
				if err := models.InsertDeviceStatus(tx, models.DeviceStatus{SN: sn, WarrantyDate: now}); err != nil {
					defer tx.Rollback()
					response.ErrRaw(err).Send(w)
					return
				}
			} else {
				ds, err := models.DeviceStatusBySn(tx, sn)
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
		}

	} else {
		for _, sn := range affected {
			if err := models.UpdateDeviceStatusWarrantyDate(tx, sn, null.Time{}); err != nil {
				defer tx.Rollback()
				response.ErrRaw(err).Send(w)
				return
			}
		}
	}
	if err := tx.Commit(); err != nil {
		logger.Error(err)
		defer tx.Rollback()
	}

	ng := &pb.NewGroup{Id: dg.Gp, Sn: affected}
	// for _, ch := range u.ShipListener() {
	// 	go func(ch chan *pb.NewGroup) {
	// 		ch <- ng
	// 	}(ch)
	// }
	proto.GroupChanged(ng)

	response.NewSuccessResponse(dg).Send(w)

}
