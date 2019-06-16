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

type AccountPwdHandler struct {
}

func NewAccountPwdHandler() AccountPwdHandler {
	return AccountPwdHandler{}
}

const (
	AccountPwdAPI = "/account/pwd"
)

func (h AccountPwdHandler) SetRouter(r *Router) {
	r.PUT(AccountPwdAPI, PutAccountPwd, mds...)
}

func PutAccountPwd(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	var new IDPwd
	if err := json.NewDecoder(r.Body).Decode(&new); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	// check parameters
	var missing []string
	if new.ID == 0 {
		missing = append(missing, "id")
	}
	if new.Pwd == "" {
		missing = append(missing, "pwd")
	}
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	// check permission
	bios, err := u.Accounts()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if !util.CheckAccount(bios, new.ID) {
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	if err := u.UpdatePassword(new.ID, new.Pwd); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}
