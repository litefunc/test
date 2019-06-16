package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"
	"strconv"
)

type ComDefaultGroupHandler struct {

}

func NewComDefaultGroupHandler() ComDefaultGroupHandler {
	return ComDefaultGroupHandler{}
}

const (
	CompanyDefaultGroupAPI = "/com/default/group"
)

func (h ComDefaultGroupHandler) SetRouter(r *Router) {
	r.GET(CompanyDefaultGroupAPI, GetCompanyDefaultGroup, mds...)
}

func GetCompanyDefaultGroup(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	var cid uint64
	com := r.URL.Query().Get("com")
	if com == "" {
		cid = u.Info().Company

	} else {
		id, err := strconv.Atoi(com)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("com", com).Send(w)
			return
		}

		cid = uint64(id)
		if cid != u.Info().Company {
			if u.Info().Type != "Noovo" || u.Info().Company != 1 {
				response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
				return
			}
		}
	}

	g, err := u.DefaultGroupByCom(cid)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(g).Send(w)

}
