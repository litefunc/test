package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"
)

type TagHandler struct {
}

func NewTagHandler() TagHandler {
	return TagHandler{}
}

const (
	TagAPI = "/tag"
)

func (h TagHandler) SetRouter(r *Router) {
	r.GET(TagAPI, GetTag, mds...)
	r.PUT(TagAPI, PutTag, mds...)
}

func GetTag(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	t, err := u.Tags()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(t).Send(w)
}

func PutTag(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	r.ParseForm()
	if err := u.Tag(models.TagList(r.Form)); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(nil).Send(w)
}
