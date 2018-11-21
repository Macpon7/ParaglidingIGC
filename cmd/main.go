package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"paragliding/pkg/handlers"
	"paragliding/pkg/storage"

	"github.com/gorilla/mux"
)

func main() {
	startTime := time.Now()

	storage.TrackDB = &storage.MongoDB{
		DatabaseURL:          "mongodb://adrianceng:Glassheisen101@ds111608.mlab.com:11608/paraglide",
		DatabaseName:         "paraglide",
		TracksCollectionName: "IGCTracks",
	}
	storage.WebhookDB = &storage.MongoDBWebHook{
		DatabaseURL:            "mongodb://adrianceng:Glassheisen101@ds111608.mlab.com:11608/paraglide",
		DatabaseName:           "paraglide",
		WebhooksCollectionName: "IGCHooks",
	}

	storage.TrackDB.Init()

	r := mux.NewRouter()

	//api base path
	r.Handle("/paragliding/api", rootHandler(startTime))

	//paragliding base path
	r.Handle("/paragliding", http.RedirectHandler("/paragliding/api", 301))

	//track path
	r.Handle("/paragliding/api/track", handlers.TracksInHandler()).Methods("POST")
	r.Handle("/paragliding/api/track", handlers.TracksOutHandler()).Methods("GET")
	r.Handle("/paragliding/api/track/{id}", handlers.TrackMetaHandler()).Methods("GET")
	r.Handle("/paragliding/api/track/{id}/{field}", handlers.TrackSpecificHandler()).Methods("GET")

	//ticker path
	r.Handle("/paragliding/api/ticker", handlers.TickerOutOldestHandler()).Methods("GET")
	r.Handle("/paragliding/api/ticker/latest", handlers.TickerOutLatestHandler()).Methods("GET")
	r.Handle("/paragliding/api/ticker/{timestamp}", handlers.TickerOutSpecificHandler()).Methods("GET")

	//webhook path
	r.Handle("/paragliding/api/webhook/new_track", handlers.WebhookRegisterHandler()).Methods("POST")
	//r.Handle("/paragliding/api/webhook/new_track/{webhook_id}", handlers.WebhookAccessHandler()).Methods("GET")
	//r.Handle("/paragliding/api/webhook/new_track/{webhook_id}", handlers.WebhookDeleteHandler()).Methods("DELETE")

	//admin path
	r.Handle("/admin/api/tracks_count", handlers.TracksCountHandler())
	r.Handle("/admin/api/tracks", handlers.TracksDeleteAllHandler())

	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("past listen and serve")

}

func getPort() string {
	var port = os.Getenv("PORT")
	return ":" + port
}

func rootHandler(t time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type metaInfo struct {
			Uptime  string `json:"uptime"`
			Info    string `json:"info"`
			Version string `json:"version"`
		}

		d := time.Since(t)
		dur := d.Seconds()
		upTime := durationFormat(dur)

		metadata := metaInfo{upTime, "Service app for IGC tracks", "v1"}
		json.NewEncoder(w).Encode(metadata)
	})
}

func durationFormat(sec float64) string {
	var days, hours, minutes, seconds float64
	upTime := "P"

	if sec > 86400 {
		seconds = math.Mod(sec, 86400.0)
		days = math.Trunc(sec / 86400.0)
		sec = seconds
		upTime += strconv.FormatFloat(days, 'f', 0, 64) + "DT"
	}
	if sec > 3600 {
		seconds = math.Mod(sec, 3600.0)
		hours = math.Trunc(sec / 3600.0)
		sec = seconds
		upTime += strconv.FormatFloat(hours, 'f', 0, 64) + "H"
	}
	if sec > 60 {
		seconds = math.Mod(sec, 60.0)
		minutes = math.Trunc(sec / 60.0)
		sec = seconds
		upTime += strconv.FormatFloat(minutes, 'f', 0, 64) + "M"
	}
	if sec != 0 {
		upTime += strconv.FormatFloat(sec, 'f', 0, 64) + "S"
	}
	return upTime
}
