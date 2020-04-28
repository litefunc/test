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
	File   *multipart.FileHeader `form:"file"`
	// or for multiple files
	// Avatars []*multipart.FileHeader `form:"avatar" binding:"required"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/formfile", func(c *gin.Context) {

		// file, header, err := c.Request.FormFile("file")
		// logger.Debug(err)
		// logger.Debug(file)
		// logger.Debug(header)
		// by, _ := ioutil.ReadAll(file)
		// logger.Debug(string(by))
		// logger.Debug(header)

		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form ProfileForm
		// in this case proper binding will be automatically selected
		if err := c.ShouldBind(&form); err != nil {
			logger.Error(err)
			logger.Debug(form)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		logger.Debug(form.File)

		c.String(http.StatusOK, "ok")
	})

	return router
}
