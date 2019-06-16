package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/auth"
	"cloud/server/portal/gcs"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/check"
	"cloud/server/portal/web/form"
	"io/ioutil"

	"cloud/server/portal/web/response"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

type ComHandler struct {
	gcs *gcs.BucketClient
}

func NewComHandler(gcs *gcs.BucketClient) ComHandler {
	return ComHandler{gcs: gcs}
}

const (
	CompanyAPI             = "/com"
)

func (h ComHandler) SetRouter(r *Router) {
	r.GET(CompanyAPI, GetCompany, mds...)
	r.POST(CompanyAPI, h.PostCompany, mds...)
	r.PUT(CompanyAPI, h.PutCompany, mds...)
	r.DELETE(CompanyAPI, h.DeleteCompany, mds...)
	r.OPTIONS(CompanyAPI, statusOK, mds...)
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
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
		company, err := u.CompanyByID(cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(company).Send(w)
		return
	}

	switch u.Info().Type {
	case "Noovo":
		cvs, err := view.DBCompanies(u.DB())
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(cvs).Send(w)

	case "SI":
		cvs, err := view.DBCompaniesById(u.DB(), u.Info().Company)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(cvs).Send(w)
	default:
		response.NewSuccessResponse(nil).Send(w)
		return
	}
}

func (h ComHandler) PostCompany(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Company != 1 {
		response.ErrNotNoovoCreateCompany().Send(w)
		return
	}

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	ffv, missing := form.CheckFormFile(r, "logo", "name", "nation", "zip", "addr", "tel", "phone", "email", "pwd")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	phone := ffv.Values["phone"]
	if phone != "" {
		if !check.ValidPhone(phone) {
			response.ErrInvalidParameter("phone number", phone).Send(w)
			return
		}
	}

	c := models.Company{
		Name:   ffv.Values["name"],
		Nation: ffv.Values["nation"],
		Zip:    ffv.Values["zip"],
		Addr:   ffv.Values["addr"],
		Tel:    ffv.Values["tel"],
		Phone:  phone,
		Note:   models.NewNullStr(r.PostFormValue("note")),
	}

	tx, err := u.DB().Begin()
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	id, err := models.NewCompany(tx, c)
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}
	c.ID = id

	pwd, err := models.HashPWD(ffv.Values["pwd"])
	if err != nil {
		response.ErrRaw(err).Send(w)
		tx.Rollback()
		return
	}

	newbio := models.Bio{Email: ffv.Values["email"], Name: c.Name + "-Admin", Company: c.ID, Admin: true, Type: "SI", Phone: c.Phone}
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
	_, err = models.InsertAccount(tx, acc)
	if err != nil {
		response.ErrInternal().Send(w)
		tx.Rollback()
		return
	}

	obj := fmt.Sprintf("logo/company/%d/%s", c.ID, ffv.Header.Filename)
	logger.Debug(obj)
	if err := h.gcs.SaveObject(obj, ffv.Buf); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		tx.Rollback()
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
		tx.Rollback()
		return
	}
	c.Logo = models.NewNullStr(url)

	if err := models.UpdateCompany(tx, c); err != nil {
		response.ErrInternal().Send(w)
		tx.Rollback()
		return
	}

	if _, err := models.CreateDefaultGroup(tx, c.ID); err != nil {
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

	response.NewSuccessResponse(c).Send(w)
}

func (h ComHandler) PutCompany(w http.ResponseWriter, r *http.Request) {
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

	com, err := u.CompanyByID(i)
	if err != nil {
		logger.Error(err)
		response.ErrInternal().Send(w)
		return
	}
	com.Name = r.PostFormValue("name")
	com.Nation = r.PostFormValue("nation")
	com.Zip = r.PostFormValue("zip")
	com.Addr = r.PostFormValue("addr")
	com.Tel = r.PostFormValue("tel")
	com.Phone = r.PostFormValue("phone")
	com.Note = models.NewNullStr(r.PostFormValue("note"))

	f, header, err := r.FormFile("logo")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error(err)
			response.ErrInternal().Send(w)
			return
		}

	} else {
		if header.Filename != com.Logo.String() {
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				logger.Error(err)
				response.ErrInternal().Send(w)
				return
			}
			defer f.Close()
			obj := fmt.Sprintf("logo/company/%d/%s", com.ID, header.Filename)

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
			com.Logo = models.NewNullStr(url)
		}
	}

	if err := u.UpdateCompany(com); err != nil {
		logger.Error(err)
		response.ErrInternal().Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}

func (h ComHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Company != 1 {
		response.ErrNotNoovoDeleteCompany().Send(w)
		return
	}

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

	bios, err := models.DbBiosByCom(u.DB(), i)
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

	if err = u.DeleteCompany(i); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	defer func() {
		auth.DeleteUsersByCom(i)
	}()

	dir := fmt.Sprintf(`logo/company/%d`, i)
	if err := h.gcs.DeleteDir(dir); err != nil {
		logger.Error(err)
		response.ErrInternal().Send(w)
		return
	}

	response.NewSuccessResponse(nil).Send(w)
}

func saveLogoLocally(r *http.Request, com uint64) (string, error) {
	var url string
	file, header, err := r.FormFile("logo")
	if err != nil {
		logger.Error(err)
		return url, err
	}

	defer file.Close()
	prefix := "assets/"
	path := fmt.Sprintf(`photo/company/%d/`, com)
	dir := prefix + path
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			logger.Error(err)
			return url, err
		}
	}

	if err != nil {
		logger.Error(err)
		return url, err
	}

	f, err := os.OpenFile(dir+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		return url, err
	}
	defer f.Close()
	io.Copy(f, file)

	host := fmt.Sprintf(`%s%s/`, os.Getenv("FILE_SERVER_HOST"), os.Getenv("FILE_SERVER_PORT"))
	url = host + path + header.Filename
	return url, nil

}
