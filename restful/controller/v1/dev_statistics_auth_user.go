package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/stat/authUser"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
	"time"
)

type DeviceStatisticsAuthUserHandler struct {
}

func NewDeviceStatisticsAuthUserHandler() DeviceStatisticsAuthUserHandler {
	return DeviceStatisticsAuthUserHandler{}
}

const (
	DeviceStatisticsAuthUserMonthlyAPI = "/dev/statistics/auth-user/monthly"
	DeviceStatisticsAuthUserWeeklyAPI  = "/dev/statistics/auth-user/weekly"
	DeviceStatisticsAuthUserDailyAPI   = "/dev/statistics/auth-user/daily"
)

func (h DeviceStatisticsAuthUserHandler) SetRouter(r *Router) {
	r.GET(DeviceStatisticsAuthUserMonthlyAPI, GetDeviceStatisticsAuthUserMonthly, mds...)
	r.GET(DeviceStatisticsAuthUserWeeklyAPI, GetDeviceStatisticsAuthUserWeekly, mds...)
	r.GET(DeviceStatisticsAuthUserDailyAPI, GetDeviceStatisticsAuthUserDaily, mds...)
}

func GetDeviceStatisticsAuthUserMonthly(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var sns []string
	for _, dev := range devs {
		sns = append(sns, dev.SN)
	}

	var query AuthUserStatQueryParameter

	now := time.Now()
	query.Year = now.Year()
	query.Month = int(now.Month())
	if y := r.URL.Query().Get("year"); y != "" {
		i, err := strconv.Atoi(y)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("year", y).Send(w)
			return
		}
		query.Year = i
	}

	if q := r.URL.Query().Get("gp"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", q).Send(w)
			return
		}

		gps, err := u.Groups()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.DeviceStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbDeviceStatAuthUsersMonthlyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = authUser.DbDeviceStatAuthUsersMonthlyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsAuthUserWeekly(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var sns []string
	for _, dev := range devs {
		sns = append(sns, dev.SN)
	}

	var query AuthUserStatQueryParameter

	now := time.Now()
	query.Year = now.Year()
	if y := r.URL.Query().Get("year"); y != "" {
		i, err := strconv.Atoi(y)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("year", y).Send(w)
			return
		}
		query.Year = i
	}
	if q := r.URL.Query().Get("gp"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", q).Send(w)
			return
		}

		gps, err := u.Groups()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.DeviceStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbDeviceStatAuthUsersWeeklyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = authUser.DbDeviceStatAuthUsersWeeklyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsAuthUserDaily(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var sns []string
	for _, dev := range devs {
		sns = append(sns, dev.SN)
	}

	var query AuthUserStatQueryParameter

	now := time.Now()
	query.Year = now.Year()
	query.Month = int(now.Month())
	if y := r.URL.Query().Get("year"); y != "" {
		i, err := strconv.Atoi(y)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("year", y).Send(w)
			return
		}
		query.Year = i
	}
	m := r.URL.Query().Get("month")
	if m == "" {
		response.ErrEmptyQueryParameter("month").Send(w)
		return
	}
	i, err := strconv.Atoi(m)
	if err != nil {
		logger.Error(err)
		response.ErrInvalidQueryParameter("month", m).Send(w)
		return
	}
	query.Month = i

	if q := r.URL.Query().Get("gp"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("gp", q).Send(w)
			return
		}

		gps, err := u.Groups()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.DeviceStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbDeviceStatAuthUsersDailyByGroup(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = authUser.DbDeviceStatAuthUsersDailyBySns(u.DB(), q.Year, q.Month, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}
