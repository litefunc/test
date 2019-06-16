package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
)

type DeviceStatusHandler struct {
}

func NewDeviceStatusHandler() DeviceStatusHandler {
	return DeviceStatusHandler{}
}

const (
	DeviceStatusAPI = "/dev/status"
)

func (h DeviceStatusHandler) SetRouter(r *Router) {
	r.GET(DeviceStatusAPI, GetDeviceStatus, mds...)
}

func GetDeviceStatus(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	com := r.URL.Query().Get("com")
	if com != "" {
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
		ds, err := view.DBDeviceStatusByCom(u.DB(), cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(ds).Send(w)
		return
	}

	gp := r.URL.Query().Get("gp")
	if gp != "" {
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
		ds, err := view.DBDeviceStatusByGroup(u.DB(), gid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(ds).Send(w)
		return
	}

	if u.Info().Type == "group" {
		gps, err := u.Groups()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		gids := models.Groups(gps).IDs()
		ds, err := view.DBDeviceStatusByGroups(u.DB(), gids)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(ds).Send(w)
		return
	}
	coms, err := u.Companies()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	cids := models.Companies(coms).IDs()

	ds, err := view.DBDeviceStatusByComs(u.DB(), cids)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(ds).Send(w)
}
