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

type GroupStatisticsVodHandler struct {
}

func NewGroupStatisticsVodHandler() GroupStatisticsVodHandler {
	return GroupStatisticsVodHandler{}
}

const (
	GroupStatisticsVodMonthlyAPI = "/group/statistics/vod/monthly"
	GroupStatisticsVodWeeklyAPI  = "/group/statistics/vod/weekly"
	GroupStatisticsVodDailyAPI   = "/group/statistics/vod/daily"
)

func (h GroupStatisticsVodHandler) SetRouter(r *Router) {
	r.GET(GroupStatisticsVodMonthlyAPI, GetGroupStatisticsVodMonthly, mds...)
	r.GET(GroupStatisticsVodWeeklyAPI, GetGroupStatisticsVodWeekly, mds...)
	r.GET(GroupStatisticsVodDailyAPI, GetGroupStatisticsVodDaily, mds...)
}

type VodStatQueryParameter struct {
	Year  int
	Month int
	Vod   int
	Gp    uint64
}

func GetGroupStatisticsVodMonthly(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	gps, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var gids []uint64
	for _, gp := range gps {
		gids = append(gids, gp.ID)
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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []vod.GroupStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbGroupStatVodsMonthlysByGpAndVod(u.DB(), q.Year, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbGroupStatVodsMonthlysByVod(u.DB(), gids, q.Year, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbGroupStatVodsMonthlysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = vod.DbGroupStatVodsMonthlys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetGroupStatisticsVodWeekly(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	gps, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var gids []uint64
	for _, gp := range gps {
		gids = append(gids, gp.ID)
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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []vod.GroupStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbGroupStatVodsWeeklysByGpAndVod(u.DB(), q.Year, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbGroupStatVodsWeeklysByVod(u.DB(), gids, q.Year, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbGroupStatVodsWeeklysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = vod.DbGroupStatVodsWeeklys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetGroupStatisticsVodDaily(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	gps, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	var gids []uint64
	for _, gp := range gps {
		gids = append(gids, gp.ID)
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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []vod.GroupStatVod

	switch {
	case q.Vod != 0 && q.Gp != 0:
		list, err = vod.DbGroupStatVodsDailysByGpAndVod(u.DB(), q.Year, q.Month, q.Gp, q.Vod)

	case q.Vod != 0:
		list, err = vod.DbGroupStatVodsDailysByVod(u.DB(), gids, q.Year, q.Month, q.Vod)

	case q.Gp != 0:
		list, err = vod.DbGroupStatVodsDailysByGp(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = vod.DbGroupStatVodsDailys(u.DB(), gids, q.Year, q.Month)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
