package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"database/sql"
	"net/http"

	"time"
)

type DeviceStatusInfoHandler struct {
}

func NewDeviceStatusInfoHandler() DeviceStatusInfoHandler {
	return DeviceStatusInfoHandler{}
}

const (
	DeviceStatusInfoAPI = "/dev/status/info"
)

func (h DeviceStatusInfoHandler) SetRouter(r *Router) {
	r.GET(DeviceStatusInfoAPI, GetDeviceStatusInfo, mds...)
}

type DeviceStatusDetail struct {
	Basic   view.DeviceStatus `json:"basic"`
	Traffic Traffic           `json:"traffic"`
}

func GetDeviceStatusInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	sn := r.URL.Query().Get("sn")
	if sn == "" {
		response.ErrInvalidQueryParameter("sn", sn).Send(w)
		return
	}
	db := u.DB()
	dev, err := models.DbDeviceBySn(db, sn)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	switch u.Info().Type {
	case "Noovo":
	case "SI":
		if u.Info().Company != dev.Company {
			response.ErrInsufficientPrivilege().Send(w)
			return
		}
	case "group":
		gp, err := models.DbGroupByAccountID(db, u.Info().ID)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		n, err := models.DbCountDeviceGroupByDeviceGroup(db, sn, gp)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if n == 0 {
			response.ErrInsufficientPrivilege().Send(w)
			return
		}
	default:
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	ds, err := view.DBDeviceStatusBySn(db, sn)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	now := time.Now().UTC()
	// now := time.Date(1999, 3, 1, 0, 0, 0, 0, time.UTC)

	var tr Traffic
	if err := db.QueryRow(`SELECT total FROM cloud.get_device_monthly_traffic($1, $2, $3);`, sn, now, 0).Scan(&tr.Monthly); err != nil && err != sql.ErrNoRows {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if err := db.QueryRow(`SELECT total FROM cloud.get_device_weekly_traffic($1, $2, $3);`, sn, now, 0).Scan(&tr.Weekly); err != nil && err != sql.ErrNoRows {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if err := db.QueryRow(`SELECT total FROM cloud.get_device_daily_traffic($1, $2, $3);`, sn, now, 0).Scan(&tr.Daily); err != nil && err != sql.ErrNoRows {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	dsd := DeviceStatusDetail{Basic: ds, Traffic: tr}

	response.NewSuccessResponse(dsd).Send(w)

}
