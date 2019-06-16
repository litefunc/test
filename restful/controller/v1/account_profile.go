package v1

import (
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"
)

type AccountProfileHandler struct {
}

func NewAccountProfileHandler() AccountHandler {
	return AccountHandler{}
}

const (
	AccountProfileAPI = "/account/profile"
)

func (h AccountProfileHandler) SetRouter(r *Router) {
	r.GET(AccountProfileAPI, GetAccountProfile, mds...)
}

func GetAccountProfile(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	response.NewSuccessResponse(u.Info()).Send(w)
}
