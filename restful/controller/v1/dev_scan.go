package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/gcs"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/form"
	"cloud/server/portal/web/response"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

type DeviceScanHandler struct {
	gcs *gcs.BucketClient
}

func NewDeviceScanHandler(gcs *gcs.BucketClient) DeviceScanHandler {
	return DeviceScanHandler{gcs: gcs}
}

const (
	DeviceScanAPI = "/dev/scan"
)

func (h DeviceScanHandler) SetRouter(r *Router) {
	r.POST(DeviceScanAPI, h.PostDeviceScan, mds...)
}

func (h DeviceScanHandler) PostDeviceScan(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Company != 1 {
		response.ErrInsufficientPrivilege().Send(w)
		return
	}

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	v, missing := form.CheckForm(r, "sn")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}

	f, header, err := r.FormFile("file")
	if err != nil {
		logger.Error(err)
		if err != http.ErrMissingFile {
			response.ErrInternal().Send(w)
			return
		}
		response.ErrRaw(err).Send(w)
		return

	}

	sn := v["sn"]
	obj := fmt.Sprintf("scan/%d/%s", sn, header.Filename)
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

	ds := models.NewDeviceScan(sn, header.Filename, h.gcs.Bucket, obj, url)

	if err := models.DbInsertDeviceScan(u.DB(), ds); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(ds).Send(w)
}
