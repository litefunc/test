package v1

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"cloud/server/portal/models"
	"cloud/server/portal/models/view"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/form"
	"cloud/server/portal/web/response"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"
	"time"
)

type DeviceHandler struct {
}

func NewDeviceHandler() DeviceHandler {
	return DeviceHandler{}
}

const (
	DeviceAPI = "/dev"
)

func (h DeviceHandler) SetRouter(r *Router) {
	r.GET(DeviceAPI, GetDevice, mds...)
	r.POST(DeviceAPI, PostDevice, mds...)
	r.PUT(DeviceAPI, PutDevice, mds...)
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	gp := r.URL.Query().Get("gp")
	if gp != "" {
		c, err := strconv.ParseUint(gp, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", gp).Send(w)
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
		ds, err := view.DBDeviceGroupsByGroup(u.DB(), gid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(ds).Send(w)
		return
	}

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
		ds, err := view.DBDeviceGroupsByCom(u.DB(), cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(ds).Send(w)
		return
	}

	if u.Info().Type == "group" {
		gps, err := u.Groups()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		gids := models.Groups(gps).IDs()
		ds, err := view.DBDeviceGroupsByGroups(u.DB(), gids)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(view.DeviceGroups(ds).Devices()).Send(w)
		return
	}

	d, err := getDevices(u)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(d).Send(w)
}

func getDevices(u models.User) (ds []view.DeviceGroup, err error) {
	var rows *sql.Rows

	switch {

	// noovo can view all accounts
	case u.Info().Type == "Noovo" && u.Info().Company == 1:
		rows, err = u.DB().Query(fmt.Sprintf(view.SelectDeviceGroupSQL, ""))

	// SI can view its company's accounts
	case u.Info().Type == "SI":
		rows, err = u.DB().Query(fmt.Sprintf(view.SelectDeviceGroupSQL, "WHERE com=$1"), u.Info().Company)

	// group monitor can vew its group's accounts
	case u.Info().Type == "group":
		gp, err := models.DbGroupByAccountID(u.DB(), u.Info().ID)
		if err != nil {
			return nil, err
		}
		rows, err = u.DB().Query(fmt.Sprintf(view.SelectDeviceGroupSQL, "WHERE gp=$1"), gp)

	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return view.ScanDeviceGroups(rows)
}

func PostDevice(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	if u.Info().Company != 1 {
		response.ErrNotNoovoCreateDevice().Send(w)
		return
	}

	if !u.Info().Admin {
		response.ErrNotAdmin(u.Info().Name).Send(w)
		return
	}

	v, missing := form.CheckForm(r, "sn", "name", "mac", "com_id")
	if len(missing) != 0 {
		response.ErrParameters(missing...).Send(w)
		return
	}
	var dev models.Device
	dev.SN = v["sn"]
	dev.Name = v["name"]
	dev.MAC = models.MAC(v["mac"])
	dev.Model = r.PostFormValue("model")
	dev.Build = r.PostFormValue("build")
	dev.Docker = r.PostFormValue("docker")
	dev.Router = r.PostFormValue("router")
	dev.Tuner = r.PostFormValue("tuner")

	com, err := strconv.ParseUint(v["com_id"], 10, 64)
	if err != nil {
		response.ErrInvalidFormField("com_id").Send(w)
		return
	}
	dev.Company = com

	if d := r.PostFormValue("disk"); d != "" {
		disk, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			response.ErrInvalidFormField("disk").Send(w)
			return
		}
		dev.Disk = null.NewInt64(disk)
	}

	dev.InitDate = time.Now().UTC()

	if u.Info().Email == "localserver@noovo.co" && u.Info().Name == "localserver" {
		n, err := models.DbCountDevicesBySn(u.DB(), dev.SN)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		// if device already exists, update it
		if n > 0 {
			if err := models.DbUpdateDevice(u.DB(), dev); err != nil {
				response.ErrRaw(err).Send(w)
				return
			}

		} else {
			if err := u.NewDevice(dev); err != nil {
				response.ErrRaw(err).Send(w)
				return
			}
		}

	} else {
		if err := u.NewDevice(dev); err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
	}

	if u.Info().Email == "localserver@noovo.co" && u.Info().Name == "localserver" {
		i, err := models.DbCountDeviceGroupBySn(u.DB(), dev.SN)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		// if device has no group, default group is 1
		if i == 0 {
			var gid uint64
			if dev.Tuner == "sony" {
				gp, err := models.DbGroupByComAndName(u.DB(), 1, "Noovo sony")
				if err != nil {
					response.ErrRaw(err).Send(w)
					return
				}
				gid = gp.ID
			} else {
				gid = 1
			}
			if err := models.DbInsertDeviceGroup(u.DB(), models.DeviceGroup{SN: dev.SN, Group: gid}); err != nil {
				response.ErrRaw(err).Send(w)
				return
			}
		}
		response.NewSuccessResponse(dev).Send(w)
		return
	}

	response.NewSuccessResponse(dev).Send(w)
}

type UpdateDevice struct {
	SN   string   `json:"sn"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func PutDevice(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	var ud UpdateDevice
	if err := json.NewDecoder(r.Body).Decode(&ud); err != nil {
		logger.Error(err)
		response.ErrRaw(err).Send(w)
		return
	}
	if ud.SN == "" {
		response.ErrInvalidParameter("sn", ud.SN).Send(w)
		return
	}

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if !util.CheckDevice(devs, ud.SN) {
		response.ErrInsufficientPrivilegeOnDevice(ud.SN).Send(w)
		return
	}

	if err := models.DbUpdateDeviceName(u.DB(), ud.SN, ud.Name); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	if err := models.DbUpdateDeviceStatusTags(u.DB(), ud.SN, ud.Tags); err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(ud).Send(w)
}
