package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/stat/live"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
	"time"
)

type DeviceStatisticsLiveHandler struct {
}

func NewDeviceStatisticsLiveHandler() DeviceStatisticsLiveHandler {
	return DeviceStatisticsLiveHandler{}
}

const (
	DeviceStatisticsLiveMonthlyAPI = "/dev/statistics/live/monthly"
	DeviceStatisticsLiveWeeklyAPI  = "/dev/statistics/live/weekly"
	DeviceStatisticsLiveDailyAPI   = "/dev/statistics/live/daily"
)

func (h DeviceStatisticsLiveHandler) SetRouter(r *Router) {
	r.GET(DeviceStatisticsLiveMonthlyAPI, GetDeviceStatisticsLiveMonthly, mds...)
	r.GET(DeviceStatisticsLiveWeeklyAPI, GetDeviceStatisticsLiveWeekly, mds...)
	r.GET(DeviceStatisticsLiveDailyAPI, GetDeviceStatisticsLiveDaily, mds...)
}

func GetDeviceStatisticsLiveMonthly(w http.ResponseWriter, r *http.Request) {
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

	var query LiveStatQueryParameter

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

	query.Name = r.URL.Query().Get("name")

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
	var list []live.DeviceStatLiveData

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbDeviceStatLivesMonthlyByGroupAndLive(u.DB(), q.Year, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbDeviceStatLivesMonthlyBySnsAndLive(u.DB(), q.Year, sns, q.Name)

	case q.Gp != 0:
		list, err = live.DbDeviceStatLivesMonthlyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = live.DbDeviceStatLivesMonthlyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsLiveWeekly(w http.ResponseWriter, r *http.Request) {
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

	var query LiveStatQueryParameter

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

	query.Name = r.URL.Query().Get("name")

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
	var list []live.DeviceStatLiveData

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbDeviceStatLivesWeeklyByGroupAndLive(u.DB(), q.Year, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbDeviceStatLivesWeeklyBySnsAndLive(u.DB(), q.Year, sns, q.Name)

	case q.Gp != 0:
		list, err = live.DbDeviceStatLivesWeeklyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = live.DbDeviceStatLivesWeeklyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsLiveDaily(w http.ResponseWriter, r *http.Request) {
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

	var query LiveStatQueryParameter

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

	query.Name = r.URL.Query().Get("name")

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
	var list []live.DeviceStatLiveData

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbDeviceStatLivesDailyByGroupAndLive(u.DB(), q.Year, q.Month, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbDeviceStatLivesDailyBySnsAndLive(u.DB(), q.Year, q.Month, sns, q.Name)

	case q.Gp != 0:
		list, err = live.DbDeviceStatLivesDailyByGroup(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = live.DbDeviceStatLivesDailyBySns(u.DB(), q.Year, q.Month, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
