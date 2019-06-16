package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"net/http"
)

type ComDeviceHandler struct {
}

func NewComHandler() ComDeviceHandler {
	return ComDeviceHandler{}
}

const (
	CompanyDeviceAPI       = "/com/dev"
)

func (h ComDeviceHandler) SetRouter(r *Router) {
	r.DELETE(CompanyDeviceAPI, DeleteCompanyDevice, mds...)
}

func DeleteCompanyDevice(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	var sns []string
	if err := json.NewDecoder(r.Body).Decode(&sns); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if u.Info().Type != "Noovo" && u.Info().Type != "SI" {
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if invalid := util.CheckDevices(devs, sns); len(invalid) != 0 {
		response.ErrInsufficientPrivilegeOnDevices(invalid...).Send(w)
		return
	}

	for _, sn := range sns {

		gp, err := models.DbGroupBySn(u.DB(), sn)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		if gp == 1 {
			response.ErrDeleteDeviceFromNoovoDefault(sn).Send(w)
			return
		}

		if err := models.DbDeleteDeviceFromSI(u.DB(), sn); err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	response.NewSuccessResponse(nil).Send(w)
}
