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

type GroupStatisticsLiveHandler struct {
}

func NewGroupStatisticsLiveHandler() GroupStatisticsLiveHandler {
	return GroupStatisticsLiveHandler{}
}

const (
	GroupStatisticsLiveMonthlyAPI = "/group/statistics/live/monthly"
	GroupStatisticsLiveWeeklyAPI  = "/group/statistics/live/weekly"
	GroupStatisticsLiveDailyAPI   = "/group/statistics/live/daily"
)

func (h GroupStatisticsLiveHandler) SetRouter(r *Router) {
	r.GET(GroupStatisticsLiveMonthlyAPI, GetGroupStatisticsLiveMonthly, mds...)
	r.GET(GroupStatisticsLiveWeeklyAPI, GetGroupStatisticsLiveWeekly, mds...)
	r.GET(GroupStatisticsLiveDailyAPI, GetGroupStatisticsLiveDaily, mds...)
}

type LiveStatQueryParameter struct {
	Year  int
	Month int
	Name  string
	Gp    uint64
}

func GetGroupStatisticsLiveMonthly(w http.ResponseWriter, r *http.Request) {

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

	if q := r.URL.Query().Get("name"); q != "" {
		query.Name = q
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
	var list []live.GroupStatLive

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbGroupStatLivesMonthlysByGpAndLive(u.DB(), q.Year, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbGroupStatLivesMonthlysByLive(u.DB(), gids, q.Year, q.Name)

	case q.Gp != 0:
		list, err = live.DbGroupStatLivesMonthlysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = live.DbGroupStatLivesMonthlys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetGroupStatisticsLiveWeekly(w http.ResponseWriter, r *http.Request) {
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

	if q := r.URL.Query().Get("name"); q != "" {
		query.Name = q
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
	var list []live.GroupStatLive

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbGroupStatLivesWeeklysByGpAndLive(u.DB(), q.Year, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbGroupStatLivesWeeklysByLive(u.DB(), gids, q.Year, q.Name)

	case q.Gp != 0:
		list, err = live.DbGroupStatLivesWeeklysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = live.DbGroupStatLivesWeeklys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)

}

func GetGroupStatisticsLiveDaily(w http.ResponseWriter, r *http.Request) {
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

	if q := r.URL.Query().Get("name"); q != "" {
		query.Name = q
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
	var list []live.GroupStatLive

	switch {
	case q.Name != "" && q.Gp != 0:
		list, err = live.DbGroupStatLivesDailysByGpAndLive(u.DB(), q.Year, q.Month, q.Gp, q.Name)

	case q.Name != "":
		list, err = live.DbGroupStatLivesDailysByLive(u.DB(), gids, q.Year, q.Month, q.Name)

	case q.Gp != 0:
		list, err = live.DbGroupStatLivesDailysByGp(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = live.DbGroupStatLivesDailys(u.DB(), gids, q.Year, q.Month)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
