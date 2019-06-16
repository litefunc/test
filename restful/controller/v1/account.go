package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/auth"
	"cloud/server/portal/gcs"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
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

type AccountHandler struct {
	gcs *gcs.BucketClient
}

func NewAccountHandler(gcs *gcs.BucketClient) AccountHandler {
	return AccountHandler{gcs: gcs}
}
   
const (
	AccountAPI    = "/account"
)

func (h AccountHandler) SetRouter(r *Router) {
	r.GET(AccountAPI, GetAccount, mds...)
	r.POST(AccountAPI, h.PostAccount, mds...)
	r.PUT(AccountAPI, h.PutAccount, mds...)
	r.DELETE(AccountAPI, h.DeleteAccount, mds...)
	r.OPTIONS(AccountAPI, statusOK, mds...)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
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
		bios, err := models.DbBiosByCom(u.DB(), cid)
		if err != nil {
			response.ErrInternal().Send(w)
			return
		}
		response.NewSuccessResponse(bios).Send(w)
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
		bios, err := models.DbBiosByGroup(u.DB(), gid)
		if err != nil {
			response.ErrInternal().Send(w)
			return
		}
		response.NewSuccessResponse(bios).Send(w)
		return
	}

	b, err := u.Accounts()
	if err != nil {
		response.ErrInternal().Send(w)
		return
	}
	response.NewSuccessResponse(b).Send(w)
}

func (h AccountHandler) PostAccount(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	v, missing := form.CheckForm(r, "email", "name", "pwd")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	var err error
	var cid uint64
	com := r.PostFormValue("com_id")
	if com == "" {
		cid = u.Info().Company
	} else {
		c, err := strconv.ParseUint(com, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidFormField("com_id").Send(w)
			return
		}
		cid = uint64(c)
	}

	if u.Info().Company != cid {
		if u.Info().Company != 1 {
			response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
			return
		}
	}

	var admin bool
	ad := r.PostFormValue("admin")
	if ad != "" {
		admin, err = strconv.ParseBool(strings.ToLower(ad))
		if err != nil {
			logger.Error(err)
			response.ErrInternal().Send(w)
			return
		}
	}
	if !check.ValidEmail(v["email"]) {
		response.ErrInvalidParameter("email", v["email"]).Send(w)
		return
	}

	phone := r.PostFormValue("phone")
	// if phone != "" {
	// 	if !check.ValidPhone(phone) {
	// 		response.ErrInvalidParameter("phone number", phone).Send(w)
	// 		return
	// 	}
	// }

	note := r.PostFormValue("note")

	newbio := models.Bio{Email: v["email"], Name: v["name"], Company: cid, Admin: admin, Phone: phone, Note: models.NewNullStr(note)}
	if newbio.Company == 1 {
		newbio.Type = "Noovo"
	} else {
		if u.Info().Type == "group" {
			newbio.Type = "group"
		} else {
			newbio.Type = "SI"
		}
	}

	pwd, err := models.HashPWD(v["pwd"])
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	acc := models.Account{Bio: newbio, Pwd: string(pwd)}

	n, err := u.CountAccountByEmail(acc.Email)
	if err != nil {
		response.ErrInternal().Send(w)
		return
	}
	if n > 0 {
		response.ErrEmailAlreadyExists(acc.Email).Send(w)
		return
	}

	id, err := u.InsertAccount(acc)
	if err != nil {
		response.ErrInternal().Send(w)
		return
	}
	newbio.ID = id

	var photo string
	f, header, err := r.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error(err)
			response.ErrInternal().Send(w)
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
			response.ErrInternal().Send(w)
			return
		}
		photo = url
	}
	newbio.Photo = models.NewNullStr(photo)

	if err := u.UpdateAccount(newbio); err != nil {
		response.ErrRaw(err).Send(w)
	}

	if u.Info().Type == "group" {
		gp, err := u.GroupByAccountID(u.Info().ID)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if err := u.InsertGroupAccount(models.GroupAccount{Group: gp, ID: newbio.ID}); err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	response.NewSuccessResponse(newbio).Send(w)
}

type Bio struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	Phone string `json:"phone"`
	Note  string `json:"note"`
}

func (h AccountHandler) PutAccount(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	i, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidFormField("id").Send(w)
		return
	}

	var new models.Bio
	new.ID = i
	new.Email = r.PostFormValue("email")
	new.Name = r.PostFormValue("name")
	new.Phone = r.PostFormValue("phone")
	new.Note = models.NewNullStr(r.PostFormValue("note"))

	admin, err := strconv.ParseBool(strings.ToLower(r.PostFormValue("admin")))
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	new.Admin = admin
	new.Note = models.NewNullStr(r.PostFormValue("note"))
	// check parameters
	var missing []string
	if new.ID == 0 {
		missing = append(missing, "id")
	}
	if new.Email == "" {
		missing = append(missing, "email")
	}
	if new.Name == "" {
		missing = append(missing, "name")
	}
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}
	if !check.ValidEmail(new.Email) {
		response.ErrInvalidParameter("email", new.Email).Send(w)
		return
	}
	// if new.Phone != "" {
	// 	if !check.ValidPhone(new.Phone) {
	// 		response.ErrInvalidParameter("phone number", new.Phone).Send(w)
	// 		return
	// 	}
	// }

	// check permission
	if new.Admin && !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}
	bios, err := u.Accounts()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if !util.CheckAccount(bios, new.ID) {
		response.ErrInsufficientPrivilegeOnAccount(new.Email).Send(w)
		return
	}

	bio, err := models.DbBioByID(u.DB(), new.ID)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	// user's company and type are the same, can not be changed
	bio.Email = new.Email
	bio.Name = new.Name
	bio.Admin = new.Admin
	bio.Phone = new.Phone
	bio.Note = new.Note

	f, header, err := r.FormFile("photo")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error(err)
			response.ErrInternal().Send(w)
			return
		}

	} else {
		if header.Filename != bio.Photo.String() {
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				logger.Error(err)
				response.ErrInternal().Send(w)
				return
			}
			defer f.Close()
			obj := fmt.Sprintf("photo/company/%d/account/%d/%s", bio.Company, bio.ID, header.Filename)

			if err := h.gcs.SaveObject(obj, buf); err != nil {
				logger.Error(err)
				response.ErrInternal().Send(w)
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
				response.ErrInternal().Send(w)
				return
			}
			logger.Debug(url)
			bio.Photo = models.NewNullStr(url)
		}
	}

	if err := models.DbUpdateAccountBio(u.DB(), bio); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}

type IDPwd struct {
	ID  uint64 `json:"id"`
	Pwd string `json:"pwd"`
}

func (h AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		response.ErrParameters("email").Send(w)
		return
	}
	if u.Info().Email == email {
		response.ErrDeleteOwnAccount().Send(w)
		return
	}

	if u.Info().Company != 1 {
		n, err := models.CountAccountByEmailandCompany(u.DB(), email, u.Info().Company)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if n == 0 {
			response.ErrInsufficientPrivilegeOnAccount(email).Send(w)
			return
		}
	}

	dir := fmt.Sprintf(`photo/company/%d/account/%d`, u.Info().Company, u.Info().ID)
	if err := h.gcs.DeleteDir(dir); err != nil {
		logger.Error(err)
		response.ErrInternal().Send(w)
		return
	}

	if err := u.DeleteAccount(email); err != nil {
		response.ErrInternal().Send(w)
		return
	}

	auth.DeleteUser(email)

	response.NewSuccessResponse(nil).Send(w)
}
