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

type ComInfoHandler struct {
	
}

func NewComInfoHandler() ComInfoHandler {
	return ComInfoHandler{}
}

const (
	CompanyInfoAPI         = "/com/info"
)

func (h ComInfoHandler) SetRouter(r *Router) {
	r.GET(CompanyInfoAPI, GetCompanyInfo, mds...)
}

func GetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	com := r.URL.Query().Get("com")
	if com == "" {
		response.ErrInvalidQueryParameter("com", com).Send(w)
		return
	}
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
	company, err := view.DBCompanyById(u.DB(), cid)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(company).Send(w)

}
