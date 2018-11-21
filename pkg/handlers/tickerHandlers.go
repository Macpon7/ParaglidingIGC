package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"paragliding/pkg/storage"

	"github.com/gorilla/mux"
)

/*TickerOutOldestHandler returns the most recent timestamp,
the first timestamp, and the last timestamp in the returning
array. It then returns an array containing the IDs of the first
five tracks in the catabase that were added, and the amount
of time the request took to process in ms.*/
func TickerOutOldestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now().Unix()
		var out storage.TickerResponse
		out = storage.TrackDB.ReadTicker()
		t2 := time.Now().Unix()
		out.Processing = t2 - t1
		json.NewEncoder(w).Encode(out)
	})
}

/*TickerOutLatestHandler returns the id of the last track to
have been added.*/
func TickerOutLatestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp := storage.TrackDB.ReadTimeStamps()
		max := len(temp) - 1
		out := temp[max]
		json.NewEncoder(w).Encode(out)
	})
}

/*TickerOutSpecificHandler returns the most recent timestamp,
the timestamp after the input, and the last timestamp in the
returning array. It then returns an array containing the IDs
of the 5 next tracks chronologically after the input, and the
amount of time the request took to process in ms. */
func TickerOutSpecificHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now().Unix()
		temp := mux.Vars(r)["timestamp"]
		tsIn, err := strconv.ParseInt(temp, 10, 64)
		if err != nil {
			log.Print("Coult not parse timestamp. Bad Request.", err)
			return
        }
		var out storage.TickerResponse
		out = storage.TrackDB.ReadSpecificTicker(tsIn)
		t2 := time.Now().Unix()
		out.Processing = t2 - t1
		json.NewEncoder(w).Encode(out)
	})
}
