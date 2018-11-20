package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Macpon7/paraglidingigc/pkg/storage"
	"github.com/gorilla/mux"
)

//TickerOutOldestHandler ...
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

//TickerOutLatestHandler ...
func TickerOutLatestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp := storage.TrackDB.ReadTimeStamps()
		max := len(temp) - 1
		out := temp[max]
		json.NewEncoder(w).Encode(out)
	})
}

//TickerOutSpecificHandler ...
func TickerOutSpecificHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now().Unix()
		vars := mux.Vars(r)
		temp := vars["timestamp"]
		tsIn, err := strconv.ParseInt(temp, 10, 64)
		if err != nil {
			//!Handle me
			return
		}
		var out storage.TickerResponse
		out = storage.TrackDB.ReadSpecificTicker(tsIn)
		t2 := time.Now().Unix()
		out.Processing = t2 - t1
		json.NewEncoder(w).Encode(out)
	})
}
