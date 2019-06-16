package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"net/http"

	"time"
)

type DeviceMsgHandler struct {
}

func NewDeviceMsgHandler() DeviceMsgHandler {
	return DeviceMsgHandler{}
}

const (
	DeviceMsgAPI = "/dev/msg"
)

func (h DeviceMsgHandler) SetRouter(r *Router) {
	r.POST(DeviceMsgAPI, PostDeviceMsg, mds...)
}

func PostDeviceMsg(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Company != 1 {
		response.ErrNotNoovoCreateDevice().Send(w)
		return
	}

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	var dm models.DeviceMsg
	if err := json.NewDecoder(r.Body).Decode(&dm); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	var missing []string
	if dm.SN == "" || dm.Msg == "" {
		missing = append(missing, "sn")
	}
	if dm.Msg == "" {
		missing = append(missing, "msg")
	}
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	dm.Time = time.Now().UTC()

	if err := models.DBInsertDeviceMsg(u.DB(), dm); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(dm).Send(w)
}
