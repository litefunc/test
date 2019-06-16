package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/stat/traffic"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
	"time"
)

type DeviceStatisticsTrafficHandler struct {
}

func NewDeviceStatisticsTrafficHandler() DeviceStatisticsTrafficHandler {
	return DeviceStatisticsTrafficHandler{}
}

const (
	DeviceStatisticsTrafficMonthlyAPI = "/dev/statistics/traffic/monthly"
	DeviceStatisticsTrafficWeeklyAPI  = "/dev/statistics/traffic/weekly"
	DeviceStatisticsTrafficDailyAPI   = "/dev/statistics/traffic/daily"
)

func (h DeviceStatisticsTrafficHandler) SetRouter(r *Router) {
	r.GET(DeviceStatisticsTrafficMonthlyAPI, GetDeviceStatisticsTrafficMonthly, mds...)
	r.GET(DeviceStatisticsTrafficWeeklyAPI, GetDeviceStatisticsTrafficWeekly, mds...)
	r.GET(DeviceStatisticsTrafficDailyAPI, GetDeviceStatisticsTrafficDaily, mds...)
}

type TrafficStatQueryParameter struct {
	Year  int
	Month int
	Gp    uint64
}

func GetDeviceStatisticsTrafficMonthly(w http.ResponseWriter, r *http.Request) {
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

	var query TrafficStatQueryParameter

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
	var list []traffic.DeviceTraffic

	switch {

	case q.Gp != 0:
		list, err = traffic.DbDeviceTrafficsMonthlyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = traffic.DbDeviceTrafficsMonthlyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsTrafficWeekly(w http.ResponseWriter, r *http.Request) {
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

	var query TrafficStatQueryParameter

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
	var list []traffic.DeviceTraffic

	switch {

	case q.Gp != 0:
		list, err = traffic.DbDeviceTrafficsWeeklyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = traffic.DbDeviceTrafficsWeeklyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetDeviceStatisticsTrafficDaily(w http.ResponseWriter, r *http.Request) {
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

	var query TrafficStatQueryParameter

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
	var list []traffic.DeviceTraffic

	switch {

	case q.Gp != 0:
		list, err = traffic.DbDeviceTrafficsDailyByGroup(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = traffic.DbDeviceTrafficsDailyBySns(u.DB(), q.Year, q.Month, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
