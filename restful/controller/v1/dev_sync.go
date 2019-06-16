package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
)

type DeviceSyncHandler struct {
}

func NewDeviceSyncHandler() DeviceSyncHandler {
	return DeviceSyncHandler{}
}

const (
	DeviceSyncAPI = "/dev/sync"
)

func (h DeviceSyncHandler) SetRouter(r *Router) {
	r.GET(DeviceSyncAPI, GetDeviceSync, mds...)
}

func GetDeviceSync(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	var devs []models.Device
	var err error

	gp := r.URL.Query().Get("gp")
	com := r.URL.Query().Get("com")

	switch {

	case gp != "":
		c, err := strconv.ParseUint(gp, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", com).Send(w)
			return
		}
		gps, err := u.Groups()
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		gid := uint64(c)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		devs, err = models.DbDevicesByGroup(u.DB(), gid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	case com != "":
		c, err := strconv.ParseUint(com, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("com", com).Send(w)
			return
		}
		coms, err := u.Companies()
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		cid := uint64(c)
		if !util.CheckCompany(coms, cid) {
			response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
			return
		}
		devs, err = models.DbDevicesByCom(u.DB(), cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	default:
		devs, err = u.Devices()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	var sns []string
	for _, dev := range devs {
		sns = append(sns, dev.SN)
	}

	list, err := models.DbDeviceConfigSync(u.DB(), sns)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
