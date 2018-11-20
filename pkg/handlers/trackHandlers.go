package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/adrianceng/paraglidingigc/pkg/storage"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

type trackIn struct {
	URL string `json:"url"`
}

type trackOut struct {
	Hdate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	GliderType  string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
	TrackURL    string    `json:"track_src_url"`
}

//TracksInHandler parses an IGC file and sends it to the database
func TracksInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var input trackIn

		err := dec.Decode(&input)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		tempTrack, err := igc.ParseLocation(input.URL)

		tempLength := 0.0
		for i := 0; i < len(tempTrack.Points)-1; i++ {
			tempLength += tempTrack.Points[i].Distance(tempTrack.Points[i+1])
		}

		var tempMeta = storage.TrackMetaInf{
			ID:          "",
			Timestamp:   0,
			Pilot:       tempTrack.Header.Pilot,
			GliderType:  tempTrack.Header.GliderType,
			GliderID:    tempTrack.Header.GliderID,
			TrackLength: tempLength,
			Hdate:       tempTrack.Header.Date,
			TrackURL:    input.URL,
		}

		idStr := storage.TrackDB.AddTrack(tempMeta)
		NotifyWebhookSubscribers()
		json.NewEncoder(w).Encode(idStr)
	})
}

//TracksOutHandler returns the array of ids currently in the database
func TracksOutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idSlice := storage.TrackDB.ReadTrackIDS()
		json.NewEncoder(w).Encode(idSlice)
	})
}

//TrackMetaHandler returns all the stored info on one track
func TrackMetaHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		tempTrack := storage.TrackDB.ReadTrack(id)
		var responseTrack = trackOut{
			Hdate:       tempTrack.Hdate,
			Pilot:       tempTrack.Pilot,
			GliderType:  tempTrack.GliderType,
			GliderID:    tempTrack.GliderID,
			TrackLength: tempTrack.TrackLength,
			TrackURL:    tempTrack.TrackURL,
		}

		json.NewEncoder(w).Encode(responseTrack)
	})
}

//TrackSpecificHandler returns a specifiv field of info for a track
func TrackSpecificHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		field := vars["field"]
		tempTrack := storage.TrackDB.ReadTrack(id)
		switch field {
		case "pilot":
			fmt.Fprintf(w, tempTrack.Pilot)
		case "glider":
			fmt.Fprintf(w, tempTrack.GliderType)
		case "glider_id":
			fmt.Fprintf(w, tempTrack.GliderID)
		case "track_length":
			fmt.Fprintln(w, tempTrack.TrackLength)
		case "H_date":
			fmt.Fprintf(w, tempTrack.Hdate.String())
		default:
			http.Error(w, "Bad Request", 400)
			return
		}
	})
}
