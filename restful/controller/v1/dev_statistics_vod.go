package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/models/stat/vod"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"net/http"

	"strconv"
	"time"
)

type DeviceStatisticsVodHandler struct {
}

func NewDeviceStatisticsVodHandler() DeviceStatisticsVodHandler {
	return DeviceStatisticsVodHandler{}
}

const (
	DeviceStatisticsVodMonthlyAPI = "/dev/statistics/vod/monthly"
	DeviceStatisticsVodWeeklyAPI  = "/dev/statistics/vod/weekly"
	DeviceStatisticsVodDailyAPI   = "/dev/statistics/vod/daily"
)

func (h DeviceStatisticsVodHandler) SetRouter(r *Router) {
	r.GET(DeviceStatisticsVodMonthlyAPI, GetDeviceStatisticsVodMonthly, mds...)
	r.GET(DeviceStatisticsVodWeeklyAPI, GetDeviceStatisticsVodWeekly, mds...)
	r.GET(DeviceStatisticsVodDailyAPI, GetDeviceStatisticsVodDaily, mds...)
}

func GetDeviceStatisticsVodMonthly(w http.ResponseWriter, r *http.Request) {
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

	var query VodStatQueryParameter

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

	if q := r.URL.Query().Get("vod"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("vod", q).Send(w)
			return
		}
		query.Vod = i
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
	var list []vod.DeviceStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbDeviceStatVodsMonthlyByGroupAndVod(u.DB(), q.Year, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbDeviceStatVodsMonthlyBySnsAndVod(u.DB(), q.Year, sns, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbDeviceStatVodsMonthlyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = vod.DbDeviceStatVodsMonthlyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetDeviceStatisticsVodWeekly(w http.ResponseWriter, r *http.Request) {
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

	var query VodStatQueryParameter

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

	if q := r.URL.Query().Get("vod"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("vod", q).Send(w)
			return
		}
		query.Vod = i
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
	var list []vod.DeviceStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbDeviceStatVodsWeeklyByGroupAndVod(u.DB(), q.Year, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbDeviceStatVodsWeeklyBySnsAndVod(u.DB(), q.Year, sns, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbDeviceStatVodsWeeklyByGroup(u.DB(), q.Year, q.Gp)

	default:
		list, err = vod.DbDeviceStatVodsWeeklyBySns(u.DB(), q.Year, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetDeviceStatisticsVodDaily(w http.ResponseWriter, r *http.Request) {
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

	var query VodStatQueryParameter

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

	if q := r.URL.Query().Get("vod"); q != "" {
		i, err := strconv.Atoi(q)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("vod", q).Send(w)
			return
		}
		query.Vod = i
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
	var list []vod.DeviceStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbDeviceStatVodsDailyByGroupAndVod(u.DB(), q.Year, q.Month, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbDeviceStatVodsDailyBySnsAndVod(u.DB(), q.Year, q.Month, sns, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbDeviceStatVodsDailyByGroup(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = vod.DbDeviceStatVodsDailyBySns(u.DB(), q.Year, q.Month, sns)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
