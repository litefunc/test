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

type GroupStatisticsAuthUserHandler struct {
}

func NewGroupStatisticsAuthUserHandler() GroupStatisticsAuthUserHandler {
	return GroupStatisticsAuthUserHandler{}
}

const (
	GroupStatisticsAuthUserMonthlyAPI = "/group/statistics/auth-user/monthly"
	GroupStatisticsAuthUserWeeklyAPI  = "/group/statistics/auth-user/weekly"
	GroupStatisticsAuthUserDailyAPI   = "/group/statistics/auth-user/daily"
)

func (h GroupStatisticsAuthUserHandler) SetRouter(r *Router) {
	r.GET(GroupStatisticsAuthUserMonthlyAPI, GetGroupStatisticsAuthUserMonthly, mds...)
	r.GET(GroupStatisticsAuthUserWeeklyAPI, GetGroupStatisticsAuthUserWeekly, mds...)
	r.GET(GroupStatisticsAuthUserDailyAPI, GetGroupStatisticsAuthUserDaily, mds...)
}

type AuthUserStatQueryParameter struct {
	Year  int
	Month int
	Gp    uint64
}

func GetGroupStatisticsAuthUserMonthly(w http.ResponseWriter, r *http.Request) {

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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.GroupStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbGroupStatAuthUsersMonthlysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = authUser.DbGroupStatAuthUsersMonthlys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetGroupStatisticsAuthUserWeekly(w http.ResponseWriter, r *http.Request) {
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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.GroupStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbGroupStatAuthUsersWeeklysByGp(u.DB(), q.Year, q.Gp)

	default:
		list, err = authUser.DbGroupStatAuthUsersWeeklys(u.DB(), gids, q.Year)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}

func GetGroupStatisticsAuthUserDaily(w http.ResponseWriter, r *http.Request) {
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
		gid := uint64(i)
		if !util.CheckGroup(gps, gid) {
			response.ErrInsufficientPrivilegeOnGroup(gid).Send(w)
			return
		}
		query.Gp = gid
	}

	q := query
	var list []authUser.GroupStatAuthUser

	switch {
	case q.Gp != 0:
		list, err = authUser.DbGroupStatAuthUsersDailysByGp(u.DB(), q.Year, q.Month, q.Gp)

	default:
		list, err = authUser.DbGroupStatAuthUsersDailys(u.DB(), gids, q.Year, q.Month)
	}
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	response.NewSuccessResponse(list).Send(w)
}
