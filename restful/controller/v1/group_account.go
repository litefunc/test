package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/gcs"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/check"
	"cloud/server/portal/web/form"
	"cloud/server/portal/web/response"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type GroupAccountHandler struct {
	gcs *gcs.BucketClient
}

func NewGroupAccountHandler(gcs *gcs.BucketClient) GroupAccountHandler {
	return GroupAccountHandler{gcs: gcs}
}

const (
	GroupAccountAPI = "/group/account"
)

func (h GroupAccountHandler) SetRouter(r *Router) {
	r.POST(GroupAccountAPI, h.PostGroupAccount, mds...)
}

func (h GroupAccountHandler) PostGroupAccount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	v, missing := form.CheckForm(r, "email", "name", "pwd", "group_id")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	var err error

	g, err := strconv.ParseUint(v["group_id"], 10, 64)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidFormField("group_id").Send(w)
		return
	}
	gid := uint64(g)

	if u.Info().Company != 1 {
		n, err := models.CountGroupByIDandCompany(u.DB(), gid, u.Info().Company)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if n == 0 {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
	}

	var admin bool
	ad := r.PostFormValue("admin")
	if ad != "" {
		admin, err = strconv.ParseBool(strings.ToLower(ad))
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
	}
	if !check.ValidEmail(v["email"]) {
		response.ErrInvalidParameter("email", v["email"]).Send(w)
		return
	}
	phone := r.PostFormValue("phone")
	note := r.PostFormValue("note")

	cid, err := models.CompanyIDByGroup(u.DB(), gid)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	newbio := models.Bio{Email: v["email"], Name: v["name"], Company: cid, Admin: admin, Phone: phone, Note: models.NewNullStr(note)}
	newbio.Type = "group"

	pwd, err := models.HashPWD(v["pwd"])
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	acc := models.Account{Bio: newbio, Pwd: string(pwd)}

	n, err := u.CountAccountByEmail(acc.Email)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if n > 0 {
		response.ErrEmailAlreadyExists(acc.Email).Send(w)
		return
	}
	id, err := u.InsertAccount(acc)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	newbio.ID = id

	var photo string

	f, header, err := r.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

	} else {
		obj := fmt.Sprintf("photo/company/%d/account/%d/%s", cid, newbio.ID, header.Filename)
		logger.Debug(obj)

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		if err := h.gcs.SaveObject(obj, buf); err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		url, err := storage.SignedURL(h.gcs.Bucket, obj, &storage.SignedURLOptions{
			GoogleAccessID: h.gcs.Client.GetGoogleAccessID(),
			PrivateKey:     h.gcs.Client.GetPrivateKey(),
			Method:         "GET",
			Expires:        time.Now().Add(365 * 24 * time.Hour),
		})

		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		photo = url
	}
	newbio.Photo = models.NewNullStr(photo)

	if err := models.DbUpdateAccountBio(u.DB(), newbio); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if err := u.InsertGroupAccount(models.GroupAccount{ID: newbio.ID, Group: gid}); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(newbio).Send(w)
}

// func PostGroupAccount(w http.ResponseWriter, r *http.Request, u models.User, gcs *gcs.BucketClient) {

// 	if !u.Info().Admin {
// 		response.ErrNotAdmin(u.Info().Name).Send(w)
// 		return
// 	}

// 	v, missing := form.CheckForm(r, "email", "name", "pwd", "group_id")
// 	if len(missing) != 0 {
// 		response.ErrParameters(missing...).Send(w)
// 		return
// 	}

// 	var err error

// 	g, err := strconv.ParseUint(v["group_id"], 10, 64)
// 	if err != nil {
// 		logger.Error(err)
// 		response.ErrInvalidFormField("group_id").Send(w)
// 		return
// 	}
// 	gid := uint64(g)

// 	if u.Info().Company != 1 {
// 		n, err := models.CountGroupByIDandCompany(u.DB(), gid, u.Info().Company)
// 		if err != nil {
// 			response.ErrRaw(err).Send(w)
// 			return
// 		}
// 		if n == 0 {
// 			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
// 			return
// 		}
// 	}

// 	var admin bool
// 	ad := r.PostFormValue("admin")
// 	if ad != "" {
// 		admin, err = strconv.ParseBool(strings.ToLower(ad))
// 		if err != nil {
// 			logger.Error(err)
// 			response.ErrInternal().Send(w)
// 			return
// 		}
// 	}
// 	if !check.ValidEmail(v["email"]) {
// 		response.ErrInvalidParameter("email", v["email"]).Send(w)
// 		return
// 	}
// 	phone := r.PostFormValue("phone")
// 	// if phone != "" {
// 	// 	if !check.ValidPhone(phone) {
// 	// 		response.ErrInvalidParameter("phone number", phone).Send(w)
// 	// 		return
// 	// 	}
// 	// }
// 	note := r.PostFormValue("note")

// 	cid, err := models.CompanyIDByGroup(u.DB(), gid)
// 	if err != nil {
// 		response.ErrRaw(err).Send(w)
// 		return
// 	}

// 	newbio := models.Bio{Email: v["email"], Name: v["name"], Company: cid, Admin: admin, Phone: phone, Note: models.NewNullStr(note)}
// 	newbio.Type = "group"

// 	pwd, err := models.HashPWD(v["pwd"])
// 	if err != nil {
// 		response.ErrRaw(err).Send(w)
// 		return
// 	}

// 	acc := models.Account{Bio: newbio, Pwd: string(pwd)}

// 	n, err := u.CountAccountByEmail(acc.Email)
// 	if err != nil {
// 		response.ErrInternal().Send(w)
// 		return
// 	}
// 	if n > 0 {

// 		ac, err := models.DbAccountByEmail(u.DB(), acc.Email)
// 		if err != nil {
// 			response.ErrInternal().Send(w)
// 			return
// 		}

// 		m, err := u.CountGroupAccountById(ac.ID)
// 		if err != nil {
// 			response.ErrInternal().Send(w)
// 			return
// 		}

// 		// si user
// 		if m == 0 {
// 			response.ErrAssignSItoGroup().Send(w)
// 			return
// 		}

// 		i, err := u.CountGroupAccount(models.GroupAccount{ID: ac.ID, Group: gid})
// 		if err != nil {
// 			response.ErrInternal().Send(w)
// 			return
// 		}
// 		if i > 0 {
// 			response.ErrEmailAlreadyExists(acc.Email).Send(w)
// 			return
// 		}
// 		newbio.ID = ac.ID

// 	} else {
// 		id, err := u.InsertAccount(acc)
// 		if err != nil {
// 			response.ErrInternal().Send(w)
// 			return
// 		}
// 		newbio.ID = id
// 	}

// 	var photo string

// 	f, header, err := r.FormFile("photo")
// 	if err != nil {
// 		if err != http.ErrMissingFile {
// 			logger.Error(err)
// 			response.ErrInternal().Send(w)
// 			return
// 		}

// 	} else {
// 		obj := fmt.Sprintf("photo/company/%d/account/%d/%s", cid, newbio.ID, header.Filename)
// 		logger.Debug(obj)

// 		buf, err := ioutil.ReadAll(f)
// 		if err != nil {
// 			logger.Error(err)
// 			response.ErrRaw(err).Send(w)
// 			return
// 		}

// 		if err := gcs.SaveObject(obj, buf); err != nil {
// 			logger.Error(err)
// 			response.ErrRaw(err).Send(w)
// 			return
// 		}

// 		url, err := storage.SignedURL(gcs.Bucket, obj, &storage.SignedURLOptions{
// 			GoogleAccessID: gcs.Client.GetGoogleAccessID(),
// 			PrivateKey:     gcs.Client.GetPrivateKey(),
// 			Method:         "GET",
// 			Expires:        time.Now().Add(365 * 24 * time.Hour),
// 		})

// 		if err != nil {
// 			logger.Error(err)
// 			response.ErrInternal().Send(w)
// 			return
// 		}
// 		photo = url
// 	}
// 	newbio.Photo = models.NewNullStr(photo)

// 	if err := models.DbUpdateAccountBio(u.DB(), newbio); err != nil {
// 		response.ErrRaw(err).Send(w)
// 		return
// 	}

// 	if err := u.InsertGroupAccount(models.GroupAccount{ID: newbio.ID, Group: gid}); err != nil {
// 		response.ErrRaw(err).Send(w)
// 		return
// 	}

// 	response.NewSuccessResponse(newbio).Send(w)
// }
