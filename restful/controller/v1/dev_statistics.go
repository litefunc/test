package v1

import (
	"cloud/lib/logger"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"math"
	"net/http"
	"sort"
	"time"

	"strconv"
)

type DeviceStatisticsHandler struct {
}

func NewDeviceStatisticsHandler() DeviceStatisticsHandler {
	return DeviceStatisticsHandler{}
}

const (
	DeviceStatisticsAPI      = "/dev/statistics"
	DeviceStatisticsTotalAPI = "/dev/statistics/total"
)

func (h DeviceStatisticsHandler) SetRouter(r *Router) {
	r.GET(DeviceStatisticsAPI, GetDeviceStatistics, mds...)
	r.GET(DeviceStatisticsTotalAPI, GetDeviceStatisticsTotal, mds...)
}

type DeviceStatistics struct {
	Traffic []TrafficData `json:"traffic"`
	Ads     []models.Ad   `json:"ads"`
}

func GetDeviceStatistics(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	sn := r.URL.Query().Get("sn")
	if sn == "" {
		response.ErrInvalidQueryParameter("sn", sn).Send(w)
		return
	}
	var err error
	devs, err := u.Devices()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}
	if !util.CheckDevice(devs, sn) {
		response.ErrInsufficientPrivilegeOnDevice(sn).Send(w)
		return
	}

	var offset, limit uint64

	if q := r.URL.Query().Get("offset"); q != "" {
		offset, err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("offset", q).Send(w)
			return
		}
	}
	if q := r.URL.Query().Get("limit"); q != "" {
		limit, err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("limit", q).Send(w)
			return
		}
	}
	list, err := models.DbDeviceTraffic(u.DB(), sn, offset, limit)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	trs := TrafficResponse(list, int(limit))
	response.NewSuccessResponse(DeviceStatistics{Traffic: trs}).Send(w)
}

func GetDeviceStatisticsTotal(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	var err error
	var offset, limit uint64

	if q := r.URL.Query().Get("offset"); q != "" {
		offset, err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("offset", q).Send(w)
			return
		}
	}
	if q := r.URL.Query().Get("limit"); q != "" {
		limit, err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("limit", q).Send(w)
			return
		}
	}
	list, err := models.DbDeviceTrafficTotal(u.DB(), offset, limit)
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	trs := TrafficResponse(list, int(limit))
	response.NewSuccessResponse(DeviceStatistics{Traffic: trs}).Send(w)

}

type TrafficData struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

func TrafficResponse(trs []models.DeviceTraffic, limit int) []TrafficData {

	var new []TrafficData

	now := time.Now()

	m := make(map[string]float64)
	var dates []string

	for i := 0; i < limit; i++ {
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
