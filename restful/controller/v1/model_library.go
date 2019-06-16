package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"strconv"
	"strings"
)

var librarySchema string = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "Model library schema",
  "definitions": {
    "demod-service": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string"
        },
        "provider": {
          "type": "string"
        },
        "pid": {
          "type": "number"
        }
      },
      "required": ["type", "provider", "pid"]
    },
    "demod-dvbt": {
      "type": "object",
      "properties": {
        "no": {
          "type": "number"
        },
        "mode": {
          "type": "number",
          "enum": [0, 1]
        },
        "upnp_enable": {
          "type": "boolean"
        },
        "upnp_service": {
          "type": "string"
        },
        "channel": {
          "type": "number"
        },
        "bandwidth": {
          "type": "number"
        },
        "services": {
          "type": "array",
          "items": { "$ref": "#/definitions/demod-service" }
        }
      },
      "required": [
        "no",
        "mode",
        "upnp_enable",
        "upnp_service",
        "channel",
        "bandwidth",
        "services"
      ]
    },
    "demod-dvbs": {
      "type": "object",
      "properties": {
        "no": {
          "type": "number"
        },
        "mode": {
          "type": "number",
          "enum": [2]
        },
        "upnp_enable": {
          "type": "boolean"
        },
        "upnp_service": {
          "type": "string"
        },
        "channel": {
          "type": "number"
        },
        "symobolrate": {
          "type": "number"
        },
        "lnb_enable": {
          "type": "boolean"
        },
        "lnb_tone": {
          "type": "string"
        },
        "lnb_level": {
          "type": "string"
        },
        "services": {
          "type": "array",
          "items": { "$ref": "#/definitions/demod-service" }
        }
      },
      "required": [
        "no",
        "mode",
        "upnp_enable",
        "upnp_service",
        "channel",
        "symobolrate",
        "lnb_enable",
        "lnb_tone",
        "lnb_level",
        "services"
      ]
    }
  },
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "model": {
      "type": "string"
    },
    "demods": {
      "type": "array",
      "items": {
        "oneOf": [
          { "$ref": "#/definitions/demod-dvbt" },
          { "$ref": "#/definitions/demod-dvbs" }
        ]
      }
    }
  },
  "required": ["name", "model", "demods"],
  "additionalProperties": false
}
`
var noovoShareSchema string = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "Share library schema",
  "properties": {
    "name": {
      "type": "string"
    },
    "com_id": {
      "type": "number"
    }
  },
  "required": ["name", "com_id"],
  "additionalProperties": false
}
`

var comShareSchema string = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "Company share library schema",
  "properties": {
    "name": {
      "type": "string"
    },
    "group_id": {
      "type": "number"
    }
  },
  "required": ["name", "group_id"],
  "additionalProperties": false
}
`

type ModelLibraryHandler struct {
}

func NewModelLibraryHandler() ModelLibraryHandler {
	return ModelLibraryHandler{}
}

const (
	NoovoModelLibraryAPI      = "/noovo/model/library"
	NoovoModelLibraryInfoAPI  = "/noovo/model/library/info"
	NoovoModelLibraryShareAPI = "/noovo/model/library/share"

	CompanyModelLibraryAPI      = "/com/model/library"
	CompanyModelLibraryInfoAPI  = "/com/model/library/info"
	CompanyModelLibraryShareAPI = "/com/model/library/share"
)

func (h ModelLibraryHandler) SetRouter(r *Router) {

	// Noovo model library
	r.GET(NoovoModelLibraryAPI, GetNoovoModelLibrary, mds...)
	r.POST(NoovoModelLibraryAPI, PostNoovoModelLibrary, mds...)
	r.PUT(NoovoModelLibraryAPI, PutNoovoModelLibrary, mds...)
	r.DELETE(NoovoModelLibraryAPI, DeleteNoovoModelLibrary, mds...)

	r.GET(NoovoModelLibraryInfoAPI, GetNoovoModelLibraryInfo, mds...)
	r.PUT(NoovoModelLibraryShareAPI, PutNoovoModelLibraryShare, mds...)

	// Company model library
	r.GET(CompanyModelLibraryAPI, GetCompanyModelLibrary, mds...)
	r.PUT(CompanyModelLibraryAPI, PutCompanyModelLibrary, mds...)
	r.DELETE(CompanyModelLibraryAPI, DeleteCompanyModelLibrary, mds...)

	r.GET(CompanyModelLibraryInfoAPI, GetCompanyModelLibraryInfo, mds...)
	r.PUT(CompanyModelLibraryShareAPI, PutCompanyModelLibraryShare, mds...)
}

func GetNoovoModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}
	GetLibraryList(w, r, u)
}

func PostNoovoModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}
	CreateLibrary(w, r, u)
}

func PutNoovoModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}
	ModifyLibrary(w, r, u)
}

func DeleteNoovoModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}
	DeleteLibrary(w, r, u)
}

func GetNoovoModelLibraryInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}

	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}

	GetLibraryInfo(w, r, u)
}

type noovoShareParams struct {
	Name    string `json:"name"`
	Company uint64 `json:"com_id"`
}

type comShareParams struct {
	Name  string `json:"name"`
	Group uint64 `json:"group_id"`
}

func PutNoovoModelLibraryShare(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {

		return
	}

	if u.Info().Type != "Noovo" {
		response.ErrNotNoovo("access").Send(w)
		return
	}

	ShareLibraryToCom(w, r, u)
}

func ShareLibraryToCom(w http.ResponseWriter, r *http.Request, u models.User) {
	var params noovoShareParams

	// check permission
	if u.Info().Type != "Noovo" || !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = util.ValidateJsonSchema(body, noovoShareSchema)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	db := u.DB()
	libId, _ := models.DbGetLibraryIdByName(db, params.Name, u.Info().Company)
	if libId == 0 {
		response.ErrRaw(errors.New("Library not found")).Send(w)
		return
	}
	_, err = models.DbCompanyById(db, params.Company)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	lib, err := models.DbGetLibraryByName(db, params.Name, u.Info().Company)
	shareLib := models.ModelLibrary{
		Name:    lib.Name,
		Model:   lib.Model,
		Company: params.Company,
	}

	shareId, err := models.DbInsertLibrary(db, shareLib)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	demods, err := models.DbGetLibraryDemodList(db, lib.Id)

	for _, demod := range demods {
		// insert library demod
		demodId, err := models.DbInsertLibraryDemod(db, shareId, demod)
		if err != nil {
			models.DbDeleteLibraryById(db, shareId)
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		switch demod.Mode {
		// DVBT/T2
		case 0, 1:
			dvbt := models.LibraryDVBT{
				Channel:   demod.Channel,
				Bandwidth: demod.Bandwidth,
			}
			_, err := models.DbInsertLibraryDVBT(db, demodId, dvbt)
			if err != nil {
				models.DbDeleteLibraryById(db, shareId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		// DVBS/S2
		case 2:
			dvbs := models.LibraryDVBS{
				Channel:     demod.Channel,
				Symobolrate: demod.Symobolrate,
				LnbEnable:   demod.LnbEnable,
				LnbTone:     demod.LnbTone,
				LnbLevel:    demod.LnbLevel,
			}
			_, err := models.DbInsertLibraryDVBS(db, demodId, dvbs)
			if err != nil {
				models.DbDeleteLibraryById(db, shareId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
		for _, service := range demod.Services {
			_, err = models.DbInsertLibraryService(db, demodId, service)
			if err != nil {
				models.DbDeleteLibraryById(db, shareId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
	}
	response.NewSuccessResponse(nil).Send(w)
}

func GetCompanyModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	GetLibraryList(w, r, u)
}

func PutCompanyModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	ModifyLibrary(w, r, u)
}

func DeleteCompanyModelLibrary(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	DeleteLibrary(w, r, u)
}

func GetCompanyModelLibraryInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}
	GetLibraryInfo(w, r, u)
}

func PutCompanyModelLibraryShare(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)
	if reflect.ValueOf(u).IsNil() {
		return
	}

	ShareLibraryToGroup(w, r, u)

}

func ShareLibraryToGroup(w http.ResponseWriter, r *http.Request, u models.User) {
	var params comShareParams

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = util.ValidateJsonSchema(body, comShareSchema)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	db := u.DB()
	gids, err := models.DbGidsByCom(db, u.Info().Company)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	libId, _ := models.DbGetLibraryIdByName(db, params.Name, u.Info().Company)
	if libId == 0 {
		response.ErrRaw(errors.New("Library not found")).Send(w)
		return
	}

	gpMatch := false
	for _, gid := range gids {
		if gid == params.Group {
			gpMatch = true
		}
	}
	if !gpMatch {
		response.ErrRaw(errors.New("Group not found")).Send(w)
		return
	}

	lib, err := models.DbGetLibraryByName(db, params.Name, u.Info().Company)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	demods, err := models.DbGetLibraryDemodList(db, lib.Id)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	device, err := models.DbDevicesByGroup(db, params.Group)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	for _, dev := range device {
		count, err := models.DbCountDeviceDemodBySn(db, dev.SN)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if count != 0 {
			err = models.DbDeleteDeviceDemodBySn(db, dev.SN)
			if err != nil {
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
		for _, demod := range demods {
			demodId, err := models.DbInsertDeviceDemod(db, dev.SN, models.DeviceDemod(demod))
			if err != nil {
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
			switch demod.Mode {
			// DVBT/T2
			case 0, 1:
				dvbt := models.DemodDVBT{
					Channel:   demod.Channel,
					Bandwidth: demod.Bandwidth,
				}
				_, err := models.DbInsertDemodDVBT(db, demodId, dvbt)
				if err != nil {
					logger.Error(err)
					response.ErrRaw(err).Send(w)
					return
				}
				// DVBS/S2
			case 2:
				dvbs := models.DemodDVBS{
					Channel:     demod.Channel,
					Symobolrate: demod.Symobolrate,
					LnbEnable:   demod.LnbEnable,
					LnbTone:     demod.LnbTone,
					LnbLevel:    demod.LnbLevel,
				}
				_, err := models.DbInsertDemodDVBS(db, demodId, dvbs)
				if err != nil {
					logger.Error(err)
					response.ErrRaw(err).Send(w)
					return
				}
			}
			for _, service := range demod.Services {
				_, err = models.DbInsertDemodService(db, demodId, service)
				if err != nil {
					logger.Error(err)
					response.ErrRaw(err).Send(w)
					return
				}
			}
		}
	}
	response.NewSuccessResponse(nil).Send(w)
}

func GetLibraryInfo(w http.ResponseWriter, r *http.Request, u models.User) {
	// get library name from form
	name := r.FormValue("name")
	name = strings.TrimSuffix(name, "\n")
	if name == "" {
		response.ErrParameters("name").Send(w)
		return
	}

	db := u.DB()
	id, _ := models.DbGetLibraryIdByName(db, name, u.Info().Company)
	if id == 0 {
		response.ErrRaw(errors.New("Library not found")).Send(w)
		return
	}
	lib, err := models.DbGetLibraryByName(db, name, u.Info().Company)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	demods, err := models.DbGetLibraryDemodList(db, lib.Id)
	lib.Demods = demods
	response.NewSuccessResponse(lib).Send(w)
}

func GetLibraryList(w http.ResponseWriter, r *http.Request, u models.User) {
	var com uint64

	if r.FormValue("com") == "" {
		com = u.Info().Company
	} else {
		c, err := strconv.Atoi(r.FormValue("com"))
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		com = uint64(c)
	}

	if u.Info().Company > 1 && com != u.Info().Company {
		response.ErrInsufficientPrivilegeOnCompany(com).Send(w)
		return
	}

	db := u.DB()
	libs, err := models.DbGetLibrary(db, com)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	for idx, lib := range libs {
		demods, err := models.DbGetLibraryDemodList(db, lib.Id)
		libs[idx].Demods = demods
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
	}
	response.NewSuccessResponse(libs).Send(w)
}

func CreateLibrary(w http.ResponseWriter, r *http.Request, u models.User) {
	var lib models.ModelLibrary

	// check permission
	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = util.ValidateJsonSchema(body, librarySchema)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = json.Unmarshal(body, &lib)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	db := u.DB()
	n, err := models.DbCountModelDemods(db, lib.Model)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if len(lib.Demods) != int(n) {
		response.ErrRaw(errors.New("Numbers of Demods is not match")).Send(w)
		return
	}

	lib.Company = u.Info().Company
	libId, err := models.DbInsertLibrary(db, lib)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	for _, demod := range lib.Demods {
		// insert library demod
		demodId, err := models.DbInsertLibraryDemod(db, libId, demod)
		if err != nil {
			models.DbDeleteLibraryById(db, libId)
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		switch demod.Mode {
		// DVBT/T2
		case 0, 1:
			dvbt := models.LibraryDVBT{
				Channel:   demod.Channel,
				Bandwidth: demod.Bandwidth,
			}
			_, err := models.DbInsertLibraryDVBT(db, demodId, dvbt)
			if err != nil {
				models.DbDeleteLibraryById(db, libId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		// DVBS/S2
		case 2:
			dvbs := models.LibraryDVBS{
				Channel:     demod.Channel,
				Symobolrate: demod.Symobolrate,
				LnbEnable:   demod.LnbEnable,
				LnbTone:     demod.LnbTone,
				LnbLevel:    demod.LnbLevel,
			}
			_, err := models.DbInsertLibraryDVBS(db, demodId, dvbs)
			if err != nil {
				models.DbDeleteLibraryById(db, libId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
		for _, service := range demod.Services {
			_, err = models.DbInsertLibraryService(db, demodId, service)
			if err != nil {
				models.DbDeleteLibraryById(db, libId)
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
	}
	response.NewSuccessResponse(nil).Send(w)
}

func ModifyLibrary(w http.ResponseWriter, r *http.Request, u models.User) {
	var lib models.ModelLibrary

	// check permission
	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = util.ValidateJsonSchema(body, librarySchema)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = json.Unmarshal(body, &lib)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	db := u.DB()
	libId, _ := models.DbGetLibraryIdByName(db, lib.Name, u.Info().Company)
	if libId == 0 {
		response.ErrRaw(errors.New("Library not found")).Send(w)
		return
	}

	n, err := models.DbCountModelDemods(db, lib.Model)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if len(lib.Demods) != int(n) {
		response.ErrRaw(errors.New("Numbers of Demods is not match")).Send(w)
		return
	}

	err = models.DbUpdateLibraryById(db, libId, lib.Name, lib.Model)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = models.DbDeleteLibraryDemodByLib(db, libId)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	for _, demod := range lib.Demods {
		// insert library demod
		demodId, err := models.DbInsertLibraryDemod(db, libId, demod)
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}
		switch demod.Mode {
		// DVBT/T2
		case 0, 1:
			dvbt := models.LibraryDVBT{
				Channel:   demod.Channel,
				Bandwidth: demod.Bandwidth,
			}
			_, err := models.DbInsertLibraryDVBT(db, demodId, dvbt)
			if err != nil {
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		// DVBS/S2
		case 2:
			dvbs := models.LibraryDVBS{
				Channel:     demod.Channel,
				Symobolrate: demod.Symobolrate,
				LnbEnable:   demod.LnbEnable,
				LnbTone:     demod.LnbTone,
				LnbLevel:    demod.LnbLevel,
			}
			_, err := models.DbInsertLibraryDVBS(db, demodId, dvbs)
			if err != nil {
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
		for _, service := range demod.Services {
			_, err = models.DbInsertLibraryService(db, demodId, service)
			if err != nil {
				logger.Error(err)
				response.ErrRaw(err).Send(w)
				return
			}
		}
	}
	response.NewSuccessResponse(nil).Send(w)
}

func DeleteLibrary(w http.ResponseWriter, r *http.Request, u models.User) {

	// check permission
	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Email).Send(w)
		return
	}
	// get library name from form
	name := r.FormValue("name")
	name = strings.TrimSuffix(name, "\n")
	if name == "" {
		response.ErrParameters("name").Send(w)
		return
	}

	id, _ := models.DbGetLibraryIdByName(u.DB(), name, u.Info().Company)
	if id == 0 {
		response.ErrRaw(errors.New("Library not found")).Send(w)
		return
	}

	if err := models.DbDeleteLibraryById(u.DB(), id); err != nil {
		response.ErrInternal().Send(w)
		return
	}
	response.NewSuccessResponse(nil).Send(w)
}
