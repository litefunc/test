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

type ConnectivityHandler struct {
}

func NewConnectivityHandler() ConnectivityHandler {
	return ConnectivityHandler{}
}

const (
	ConnectivityAPI = "/connectivity"
)

func (h ConnectivityHandler) SetRouter(r *Router) {
	r.GET(ConnectivityAPI, GetConnectivity, mds...)
}

func GetConnectivity(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if gp := r.URL.Query().Get("gp"); gp != "" {
		id, err := strconv.Atoi(gp)
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		gps, err := u.Groups()
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		gid := uint64(id)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}

		conns, err := u.FilterConnectivities(uint64(id))
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(conns).Send(w)
		return
	}

	if sn := r.URL.Query().Get("sn"); sn != "" {

		devs, err := u.Devices()
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		if !util.CheckDevice(devs, sn) {
			response.ErrInsufficientPrivilegeOnDevice(sn).Send(w)
			return
		}

		conns, err := models.ConnectivitiesBySn(u.DB(), sn)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(conns).Send(w)
		return
	}

	if name := r.URL.Query().Get("name"); name != "" {

		devs, err := u.Devices()
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		if !util.CheckDeviceByName(devs, name) {
			response.ErrInsufficientPrivilegeOnDevice(name).Send(w)
			return
		}

		conns, err := models.ConnectivitiesByName(u.DB(), name)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(conns).Send(w)
		return
	}

	conns, err := u.Connectivities()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(conns).Send(w)

}
