package v1

import (
	"cloud/server/portal/gcs"
	"database/sql"
	"net/http"
)

type ServeMux struct {
	Router *Router
	db     *sql.DB
	gcs    *gcs.BucketClient
}

func NewServeMux(router *Router, db *sql.DB, gcs *gcs.BucketClient) *ServeMux {

	rc := ServeMux{
		Router: router,
		db:     db,
		gcs:    gcs,
	}
	rc.setV1Router()

	return &rc
}

func (rc *ServeMux) setV1Router() {

	rc.Router.Router = rc.Router.PathPrefix("/v1").Subrouter()

	NewAccountProfileHandler().SetRouter(rc.Router)

	NewAccountPwdHandler().SetRouter(rc.Router)

	NewAccountHandler(rc.gcs).SetRouter(rc.Router)

	NewComDefaultGroupHandler().SetRouter(rc.Router)

	NewComDeviceAmountHandler().SetRouter(rc.Router)

	NewComDeviceHandler().SetRouter(rc.Router)

	NewComInfoHandler().SetRouter(rc.Router)

	NewComHandler(rc.gcs).SetRouter(rc.Router)

	NewConnectivityHandler().SetRouter(rc.Router)

	NewDeviceDVBHandler().SetRouter(rc.Router)

	NewDeviceGroupHandler().SetRouter(rc.Router)

	NewDeviceModelHandler().SetRouter(rc.Router)

	NewDeviceMsgHandler().SetRouter(rc.Router)

	NewDeviceScanHandler(gcs).SetRouter(rc.Router)

	NewDeviceStatisticsAuthUserHandler().SetRouter(rc.Router)

	NewDeviceStatisticsLiveHandler().SetRouter(rc.Router)

	NewDeviceStatisticsTrafficHandler().SetRouter(rc.Router)

	NewDeviceStatisticsVodHandler().SetRouter(rc.Router)

	NewDeviceStatisticsHandler().SetRouter(rc.Router)

	NewDeviceStatusInfoHandler().SetRouter(rc.Router)

	NewDeviceStatusHandler().SetRouter(rc.Router)

	NewDeviceSyncHandler().SetRouter(rc.Router)

	NewDeviceHandler().SetRouter(rc.Router)

	NewGroupAccountHandler(gcs).SetRouter(rc.Router)

	NewGroupDeviceAmountHandler().SetRouter(rc.Router)

	NewGroupInfoHandler().SetRouter(rc.Router)

	NewGroupStatisticsAuthUserHandler().SetRouter(rc.Router)

	NewGroupStatisticsLiveHandler().SetRouter(rc.Router)

	NewGroupStatisticsVodHandler().SetRouter(rc.Router)

	NewGroupStatisticsHandler().SetRouter(rc.Router)

	NewGroupHandler(gcs).SetRouter(rc.Router)

	NewModelLibraryHandler().SetRouter(rc.Router)

	NewTagHandler().SetRouter(rc.Router)
}

func statusOK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}
