package handlers

import (
	"encoding/json"
	"net/http"

	"paraglidingigc/pkg/storage"
)

//TracksCountHandler returns the total amount of tracks in the database
func TracksCountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := storage.TrackDB.CountTracks()
		json.NewEncoder(w).Encode(resp)
	})
}

//TracksDeleteAllHandler deletes all the tracks in the database
func TracksDeleteAllHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := storage.TrackDB.CountTracks()
		storage.TrackDB.DeleteTracks()
		json.NewEncoder(w).Encode(resp)
	})
}
