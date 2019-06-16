package v1

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"cloud/server/portal/models"
	"cloud/server/portal/util"
	"cloud/server/portal/web"
	"cloud/server/portal/web/response"
	"database/sql"
	"fmt"
	"net/http"

	"strconv"
	"strings"
	"time"
)

type GroupInfoHandler struct {
}

func NewGroupInfoHandler() GroupInfoHandler {
	return GroupInfoHandler{}
}

const (
	GroupInfoAPI = "/group/info"
)

func (h GroupInfoHandler) SetRouter(r *Router) {
	r.GET(GroupInfoAPI, GetGroupInfo, mds...)
}

type GroupInfo struct {
	ID         uint64         `json:"id"`
	Name       string         `json:"name"`
	Note       models.NullStr `json:"note,omitempty"`
	DevCount   int            `json:"dev_count"`
	LastUpdate LastUpdate     `json:"last_update"`
	Traffic    Traffic        `json:"traffic"`
	Status     Status         `json:"status"`
}

type LastUpdate struct {
	MsFwTime        null.Time `json:"msfw_time"`
	MsFwVersion     string    `json:"msfw_version"`
	RouterFwTime    null.Time `json:"routerfw_time"`
	RouterFwVersion string    `json:"routerfw_version"`
}

type Traffic struct {
	Daily   int `json:"daily"`
	Weekly  int `json:"weekly"`
	Monthly int `json:"monthly"`
}

type Status struct {
	State       int    `json:"state"`
	Description string `json:"description"`
}

func GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(web.UserContextKey).(models.User)

	if com := r.URL.Query().Get("com"); com != "" {
		cid, err := strconv.ParseUint(com, 10, 64)
		if err != nil {
			logger.Error(err)
			response.ErrInvalidQueryParameter("com", com).Send(w)
			return
		}
		coms, err := u.Companies()
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		if !util.CheckCompany(coms, cid) {
			response.ErrInsufficientPrivilegeOnCompany(cid).Send(w)
			return
		}
		gps, err := models.DbGidsByCom(u.DB(), cid)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		gis, err := DBGroupInfo(u.DB(), gps)
		if err != nil {
			response.ErrRaw(err).Send(w)
			return
		}
		response.NewSuccessResponse(gis).Send(w)
		return
	}

	g, err := u.Groups()
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	gis, err := DBGroupInfo(u.DB(), util.GroupIDs(g))
	if err != nil {
		response.ErrRaw(err).Send(w)
		return
	}

	response.NewSuccessResponse(gis).Send(w)
}

func DBGroupInfo(db *sql.DB, gps []uint64) ([]GroupInfo, error) {

	if len(gps) == 0 {
		return nil, nil
	}
	query := `SELECT a.id, a.name, a.note, a.dev_count, 
	b.msfw_time, COALESCE(b.msfw_version,'') AS msfw_version, b.routerfw_time, COALESCE(b.routerfw_version,'') AS routerfw_version 
	FROM cloud.group_view AS a LEFT JOIN 
	cloud.group_last_update AS b on a.id=b.gp WHERE a.id in (%s);`

	rows, err := db.Query(fmt.Sprintf(query, strings.Join(util.Uint64sToStrs(gps), ", ")))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	var gis []GroupInfo
	for rows.Next() {
		var gi GroupInfo
		var lu LastUpdate
		var tr Traffic

		if err := rows.Scan(&gi.ID, &gi.Name, &gi.Note, &gi.DevCount, &lu.MsFwTime, &lu.MsFwVersion, &lu.RouterFwTime, &lu.RouterFwVersion); err != nil {
			logger.Error(err)
			return nil, err
		}
		gi.LastUpdate = lu

		now := time.Now().UTC()
		// now := time.Date(1999, 3, 1, 0, 0, 0, 0, time.UTC)

		if err := db.QueryRow(`SELECT total FROM cloud.get_group_monthly_traffic($1, $2, $3) LIMIT 1;`, gi.ID, now, 0).Scan(&tr.Monthly); err != nil && err != sql.ErrNoRows {
			logger.Error(err)
			return nil, err
		}
		if err := db.QueryRow(`SELECT total FROM cloud.get_group_weekly_traffic($1, $2, $3) LIMIT 1;`, gi.ID, now, 0).Scan(&tr.Weekly); err != nil && err != sql.ErrNoRows {
			logger.Error(err)
			return nil, err
		}
		if err := db.QueryRow(`SELECT total FROM cloud.get_group_daily_traffic($1, $2, $3) LIMIT 1;`, gi.ID, now, 0).Scan(&tr.Daily); err != nil && err != sql.ErrNoRows {
			logger.Error(err)
			return nil, err
		}
		gi.Traffic = tr

		gis = append(gis, gi)
	}
	return gis, nil
}
