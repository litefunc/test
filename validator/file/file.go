package file

import (
	"cloud/lib/logger"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type va struct {
	A string `form:"a" json:"a" binding:"required"`
	B string `form:"b" json:"b" binding:"required"`
	C string `form:"c" json:"c" `
}

type ProfileForm struct {
	Name   string                `form:"name" binding:"required"`
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`

	// or for multiple files
	// Avatars []*multipart.FileHeader `form:"avatar" binding:"required"`
}

func (md va) m(exclude ...string) map[string]interface{} {
	m := map[string]interface{}{
		"a": md.A,
		"b": md.B,
		"c": md.C,
		"d": "d",
	}
	for _, k := range exclude {
		delete(m, k)
	}
	return m
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/formfile", func(c *gin.Context) {

		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form ProfileForm
		// in this case proper binding will be automatically selected
		if err := c.ShouldBind(&form); err != nil {
			logger.Error(err)
			logger.Debug(form.Avatar == nil)
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		c.String(http.StatusOK, "ok")
	})

	return router
}
