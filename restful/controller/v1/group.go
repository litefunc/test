package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/gcs"
	"cloud/server/portal/models"
	"cloud/server/portal/proto/client/ota"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/form"
	"cloud/server/portal/web/response"
	"fmt"
	"net/http"

	"strconv"
)

type GroupHandler struct {
	gcs *gcs.BucketClient
}

func NewGroupHandler(gcs *gcs.BucketClient) GroupHandler {
	return GroupHandler{gcs: gcs}
}

const (
	GroupAPI = "/group"
)

func (h GroupHandler) SetRouter(r *Router) {
	r.GET(GroupAPI, GetGroup, mds...)
	r.POST(GroupAPI, PostGroup, mds...)
	r.PUT(GroupAPI, PutGroup, mds...)
	r.DELETE(GroupAPI, h.DeleteGroup, mds...)
}

func GetGroup(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	gs, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if gp := r.URL.Query().Get("gp"); gp != "" {
		id, err := strconv.Atoi(gp)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", gp).Send(w)
			return
		}
		gid := uint64(id)
		if !util.CheckGroup(gs, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}

		g, err := models.DbGroupByID(u.DB(), gid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		response.NewSuccessResponse(g).Send(w)
		return
	}

	if com := r.URL.Query().Get("com"); com != "" {
		id, err := strconv.Atoi(com)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("com", com).Send(w)
			return
		}
		cid := uint64(id)

		coms, err := u.Companies()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		if !util.CheckCompany(coms, cid) {
			response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
			return
		}

		gs, err = models.DbGroupsByCom(u.DB(), cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		response.NewSuccessResponse(gs).Send(w)
		return
	}

	response.NewSuccessResponse(gs).Send(w)
}

func PostGroup(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	fv, missing := form.CheckForm(r, "name", "email", "pwd")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	var cid uint64
	if com := r.FormValue("com_id"); com != "" {
		c, err := strconv.ParseUint(com, 10, 64)
		if err != nil {
			response.ErrInvalidFormField("com_id").Send(w)
			return
		}
		cid = c
	} else {
		cid = u.Info().Company
	}

	gp := models.Group{
		Name:    fv["name"],
		Company: cid,
		Nation:  r.PostFormValue("nation"),
		Zip:     r.PostFormValue("zip"),
		Addr:    r.PostFormValue("addr"),
		Tel:     r.PostFormValue("tel"),
		Phone:   r.PostFormValue("phone"),
		Note:    models.NewNullStr(r.PostFormValue("note")),
	}

	tx, err := u.DB().Begin()
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	i, err := models.DbCountGroupByComAndName(u.DB(), gp.Company, gp.Name)
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	if i >= 1 {
		response.ErrGroupAlreadyExists(gp.Name).Send(w)
		tx.Rollback()
		return
	}

	gid, err := models.InsertGroup(tx, gp)
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	gp.ID = gid

	newbio := models.Bio{Email: fv["email"], Name: gp.Name + "-Admin", Company: cid, Admin: true, Type: "group"}
	pwd, err := models.HashPWD(fv["pwd"])
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	n, err := u.CountAccountByEmail(newbio.Email)
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	if n >= 1 {
		response.ErrEmailAlreadyExists(newbio.Email).Send(w)
		tx.Rollback()
		return
	}

	acc := models.Account{Bio: newbio, Pwd: string(pwd)}
	id, err := models.InsertAccount(tx, acc)
	if err != nil {
		tx.Rollback()
		response.ErrInternal().Send(w)
		return
	}

	if err := models.InsertGroupAccount(tx, models.GroupAccount{ID: id, Group: gid}); err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	ota.NewGroup(gid)
	response.NewSuccessResponse(gp).Send(w)
}

func PutGroup(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	i, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidFormField("id").Send(w)
		return
	}

	gp, err := models.DbGroupByID(u.DB(), i)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if gp.Company != u.Info().Company && u.Info().Company != 1 {
		response.ErrInsufficientPrivilegeOnGroup(i).Send(w)
		return
	}

	c, err := strconv.ParseUint(r.PostFormValue("com_id"), 10, 64)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidFormField("com_id").Send(w)
		return
	}

	if c != u.Info().Company && u.Info().Company != 1 {
		response.ErrInsufficientPrivilegeOnCompany(c).Send(w)
		return
	}

	putGroup := models.Group{
		ID:      i,
		Name:    r.PostFormValue("name"),
		Company: c,
		Nation:  r.PostFormValue("nation"),
		Zip:     r.PostFormValue("zip"),
		Addr:    r.PostFormValue("addr"),
		Tel:     r.PostFormValue("tel"),
		Phone:   r.PostFormValue("phone"),
		Note:    models.NewNullStr(r.PostFormValue("note")),
	}

	if gp.DefaultGroup && putGroup.Name != gp.Name {
		response.ErrChangeDefaultGroupName().Send(w)
		return
	}

	if err := u.UpdateGroup(i, putGroup); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}

func (h GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}
	i, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidFormField("id").Send(w)
		return
	}

	gp, err := models.DbGroupByID(u.DB(), i)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if gp.Company != u.Info().Company && u.Info().Company != 1 {
		response.ErrInsufficientPrivilegeOnGroup(i).Send(w)
		return
	}

	if gp.DefaultGroup {
		response.ErrDeleteDefaultGroup().Send(w)
		return
	}

	bios, err := models.DbBiosByGroup(u.DB(), i)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	for _, bio := range bios {

		dir := fmt.Sprintf(`photo/company/%d/account/%d`, bio.Company, bio.ID)
		if err := h.gcs.DeleteDir(dir); err != nil {
			logger.Error(err)
			response.ErrInternal().Send(w)
			return
		}

		if err := u.DeleteAccount(bio.Email); err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	if err = u.DeleteGroup(i); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}
