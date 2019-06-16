package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"
)

type DeviceModelHandler struct {
}

func NewDeviceModelHandler() DeviceModelHandler {
	return DeviceModelHandler{}
}

const (
	DeviceModelAPI = "/dev/model"
)

func (h DeviceModelHandler) SetRouter(r *Router) {
	r.GET(DeviceModelAPI, GetDeviceModel, mds...)
}

func GetDeviceModel(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	db := u.DB()
	models, err := models.DbModelDemods(db)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(models).Send(w)

}
