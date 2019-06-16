package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"math"
	"net/http"
	"sort"

	"strconv"
	"time"
)

type GroupStatisticsTrafficHandler struct {
}

func NewGroupStatisticsTrafficHandler() GroupStatisticsTrafficHandler {
	return GroupStatisticsTrafficHandler{}
}

const (
	GroupStatisticsAPI = "/group/statistics"
)

func (h GroupStatisticsTrafficHandler) SetRouter(r *Router) {
	r.GET(GroupStatisticsAPI, GetGroupStatistics, mds...)
}

type groupQuery struct {
	GId     uint64
	Date    string
	Regular string
}

type StatisticsResponse struct {
	Traffic []TrafficData `json:"traffic"`
	// To do
	// Ads click count
}

func GetGroupStatistics(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(web.UserContextKey).(models.User)

	q := r.URL.Query()

	// validate gp query string
	gp, err := strconv.ParseUint(q.Get("gp"), 10, 64)
	if err != nil {
		response.ErrInvalidFormField("gp").Send(w)
		return
	}
	// validate regular query string
	regular := q.Get("regular")
	if regular != "weekly" && regular != "daily" && regular != "monthly" {
		response.ErrInvalidFormField("regular").Send(w)
		return
	}
	// validate date query string
	// ISO 8601
	date := q.Get("date")
	_, err = time.Parse("2006-1-2", date)
	if err != nil {
		response.ErrInvalidFormField("date").Send(w)
		return
	}

	groups, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	allow := false
	for _, group := range groups {
		if group.ID == gp {
			allow = true
		}
	}

	if allow == false {
		response.ErrInsufficientPrivilegeOnGroup(gp).Send(w)
		return
	}

	gpQuery := &groupQuery{
		GId:     gp,
		Date:    date,
		Regular: regular,
	}

	logger.Info("/group/statistics ", r.Method, gpQuery)

	getGroupStatistics(w, r, u, gpQuery)

}

func getGroupStatistics(w http.ResponseWriter, r *http.Request, u models.User, g *groupQuery) {

	var resp StatisticsResponse
	db := u.DB()
	switch g.Regular {
	case "daily":
		// return last 10 days
		s, err := models.DbGetGroupTrafficByDay(db, g.Date, g.GId, -9)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}

		date, err := time.Parse("2006-01-02", g.Date)
		if err != nil {
			logger.Error(err)
			response.ErrRaw(err).Send(w)
			return
		}

		resp.Traffic = TrafficResponseDaily(s, date)
		response.NewSuccessResponse(resp).Send(w)
	case "weekly":
		// return last 10 weeks
		s, err := models.DbGetGroupTrafficByWeek(db, g.Date, g.GId, -9)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		resp.Traffic = TrafficDataFromModel(s)
		response.NewSuccessResponse(resp).Send(w)
	case "monthly":
		// return last 10 months
		s, err := models.DbGetGroupTrafficByMonth(db, g.Date, g.GId, -9)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		resp.Traffic = TrafficDataFromModel(s)
		response.NewSuccessResponse(resp).Send(w)
	}
}

func TrafficResponseDaily(trs []models.GroupTraffic, now time.Time) []TrafficData {
	var new []TrafficData

	m := make(map[string]float64)
	var dates []string

	for i := 0; i < 10; i++ {
		h := i * (-24)
		t := now.Add(time.Duration(h) * time.Hour)
		m[t.Format("2006-01-02")] = 0
		dates = append(dates, t.Format("2006-01-02"))
	}

	for _, tr := range trs {
		if _, ok := m[tr.Date]; ok {
			mb := float64(tr.Total) / 1000
			m[tr.Date] = math.Round(mb)
		}
	}

	sort.Strings(dates)
	for _, d := range dates {
		tr := TrafficData{Date: d, Total: m[d]}
		new = append(new, tr)
	}

	return new
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func TrafficDataFromModel(trs []models.GroupTraffic) []TrafficData {
	var new []TrafficData
	for _, tr := range trs {

		mb := float64(tr.Total) / 1000
		new = append(new, TrafficData{Date: tr.Date, Total: math.Round(mb)})
	}
	return new
}
