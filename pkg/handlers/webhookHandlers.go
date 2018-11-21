package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"paragliding/pkg/storage"
	"strconv"
	"time"
)

//WebhookOut returns
type WebhookOut struct {
	Content string `json:"content"`
}

//WebhookRegisterHandler registers new webhooks
func WebhookRegisterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var input storage.WebhookInfo

		err := dec.Decode(&input)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		idStr := storage.WebhookDB.AddWebhook(input)

		json.NewEncoder(w).Encode(idStr)
	})
}

//NotifyWebhookSubscribers notifies all subscribers
func NotifyWebhookSubscribers() {
	startTime := time.Now()
	hooks := storage.WebhookDB.CheckWebhooks()
	trackIDS := storage.TrackDB.ReadTrackIDS()
	maxID := len(trackIDS) - 1
	for key := range hooks {
		out := WebhookOut{}
		out.Content = "Latest timestamp: "
		out.Content += strconv.FormatInt(storage.TrackDB.ReadTrack(trackIDS[maxID]).Timestamp, 10)
		out.Content += ", "
		out.Content += strconv.Itoa(hooks[key].TriggerValue)
		out.Content += " new tracks are: "
		if hooks[key].TriggerValue > 1 {
			for i := maxID - hooks[key].TriggerValue + 1; i > maxID; i++ {
				out.Content += trackIDS[i]
				out.Content += ", "
			}
		}
		out.Content += trackIDS[maxID]
		out.Content += ". (processing: "
		processTime := ((time.Now().UnixNano() - startTime.UnixNano()) / 1000000) % 1000
		out.Content += strconv.FormatInt((processTime/1000)%60, 10)
		out.Content += "s "
		out.Content += strconv.FormatInt(processTime, 10)
		out.Content += "ms)\n"
		raw, err := json.Marshal(out)
		if err != nil {
			log.Print("Could not marshal", err)
			return
		}
		resp, err := http.Post(hooks[key].URL, "application/json", bytes.NewBuffer(raw))
		if err != nil {
			log.Print("Could not notify webhook")
			return
		}
		defer resp.Body.Close()
	}
}
