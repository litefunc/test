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

	//"strconv"
	"strings"
)

var modifyDVBSchema string = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "Modify device DVB setting",
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
    "sn": {
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
  "required": ["sn", "demods"],
  "additionalProperties": false
}
`

type DeviceDVBHandler struct {
}

func NewDeviceDVBHandler() DeviceDVBHandler {
	return DeviceDVBHandler{}
}

const (
	DeviceDVBAPI     = "/dev/dvb"
	DeviceDVBInfoAPI = "/dev/dvb/info"
)

func (h DeviceDVBHandler) SetRouter(r *Router) {
	r.GET(DeviceDVBInfoAPI, GetDeviceDVBInfo, mds...)
	r.PUT(DeviceDVBAPI, PutDeviceDVB, mds...)
}

func GetDeviceDVBInfo(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	// get sn from form
	sn := r.FormValue("sn")
	sn = strings.TrimSuffix(sn, "\n")
	if sn == "" {
		response.ErrParameters("sn").Send(w)
		return
	}

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if !util.CheckDevice(devs, sn) {
		response.ErrInsufficientPrivilegeOnDevice(sn).Send(w)
		return
	}

	demods, err := models.DbGetDeviceDemodList(u.DB(), sn)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(demods).Send(w)
}

func PutDeviceDVB(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	var dev UpdateDeviceDVB

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = util.ValidateJsonSchema(body, modifyDVBSchema)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	err = json.Unmarshal(body, &dev)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if !util.CheckDevice(devs, dev.SN) {
		response.ErrInsufficientPrivilegeOnDevice(dev.SN).Send(w)
		return
	}

	db := u.DB()
	d, err := models.DbDeviceBySn(db, dev.SN)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	n, err := models.DbCountModelDemods(db, d.Model)
	if err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if len(dev.Demods) != int(n) {
		response.ErrRaw(errors.New("Numbers of Demods is not match")).Send(w)
		return
	}

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

	for _, demod := range dev.Demods {
		// insert device demod
		demodId, err := models.DbInsertDeviceDemod(db, dev.SN, demod)
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
	response.NewSuccessResponse(nil).Send(w)
}

type UpdateDeviceDVB struct {
	SN     string               `json:"sn"`
	Demods []models.DeviceDemod `json:"demods"`
}
