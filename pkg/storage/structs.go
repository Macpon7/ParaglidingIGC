package storage

import (
	"time"

	bson "github.com/globalsign/mgo/bson"
)

//TickerResponse represents the ticker
type TickerResponse struct {
	TLatest    int64           `json:"t_latest"`
	TStart     int64           `json:"t_start"`
	TStop      int64           `json:"t_stop"`
	TrackIDs   []bson.ObjectId `json:"tracks"`
	Processing int64           `json:"processing"`
}

//MongoDB stores the details of the DB connection
type MongoDB struct {
	DatabaseURL          string
	DatabaseName         string
	TracksCollectionName string
}

//MongoDBWebHook ...
type MongoDBWebHook struct {
	DatabaseURL            string
	DatabaseName           string
	WebhooksCollectionName string
}

//TrackMetaInf represents the main persistent data structure
type TrackMetaInf struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Timestamp   int64         `json:"timestamp"`
	Hdate       time.Time     `json:"H_date"`
	Pilot       string        `json:"pilot"`
	GliderType  string        `json:"glider"`
	GliderID    string        `json:"glider_id"`
	TrackLength float64       `json:"track_length"`
	TrackURL    string        `json:"track_src_url"`
}

//WebhookInfo contains the webhook url and the frequency of updates
type WebhookInfo struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	URL            string        `json:"webhookURL"`
	TriggerValue   int           `json:"minTriggerValue"`
	LastTrackCount int           `bson:"lastTrackCount"`
}
