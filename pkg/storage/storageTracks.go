package storage

import (
	"time"

	mgo "github.com/globalsign/mgo"
	bson "github.com/globalsign/mgo/bson"
)

//TrackDB ...
var TrackDB TrackStorage

//TrackStorage interface
type TrackStorage interface {
	Init()
	AddTrack(inFile TrackMetaInf) string
	CountTracks() int
	DeleteTracks()
	ReadTrackIDS() []string
	ReadTimeStamps() []int64
	ReadTrack(id string) TrackMetaInf
	ReadTicker() TickerResponse
	ReadSpecificTicker(timestamp int64) TickerResponse
}

//Init initializes the mongo storage
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
}

//AddTrack adds a new set of meta information about a new track
func (db *MongoDB) AddTrack(inFile TrackMetaInf) string {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	inFile.Timestamp = time.Now().Unix()
	inFile.ID = bson.NewObjectId()
	id := inFile.ID.Hex()

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Insert(inFile)
	if err != nil {
		//TODO handle this
		return ""
	}

	return id
}

//CountTracks simply returns an integer, how many tracks there are registered on the database
func (db *MongoDB) CountTracks() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	out, err := session.DB(db.DatabaseName).C(db.TracksCollectionName).Count()
	if err != nil {
		//TODO handle this
		return 0
	}
	return out
}

//DeleteTracks deletes all the tracks from the database
func (db *MongoDB) DeleteTracks() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).DropCollection()
	if err != nil {
		//TODO handle this
		return
	}
	return
}

//ReadTrackIDS returns the array of ids of all tracks in the database
func (db *MongoDB) ReadTrackIDS() []string {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var tempTracks []TrackMetaInf

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(bson.M{}).All(&tempTracks)
	if err != nil {
		//TODO handle this
		return make([]string, 0)
	}

	var idSlice []string
	for i := 0; i < len(tempTracks); i++ {
		idSlice = append(idSlice, tempTracks[i].ID.Hex())
	}

	return idSlice
}

//ReadTimeStamps returns a sorted array of every timestamp in the database
func (db *MongoDB) ReadTimeStamps() []int64 {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var response []int64

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(nil).Select(bson.M{"timestamp": response}).Sort("timestamp").All(&response)
	if err != nil {
		//TODO handle this
		return response
	}

	return response
}

//ReadTrack returns the meta information stored about a track
func (db *MongoDB) ReadTrack(id string) TrackMetaInf {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	requestID := bson.ObjectIdHex(id)
	var response TrackMetaInf

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(bson.M{"_id": requestID}).One(&response)
	if err != nil {
		//TODO handle this
		return response
	}

	return response
}

//ReadTicker reads the ticker but with no input
func (db *MongoDB) ReadTicker() TickerResponse {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var idSlice []bson.ObjectId
	var tsSlice []int64

	type idTS struct {
		ID        bson.ObjectId
		Timestamp int64
	}

	var reciever []idTS

	var out TickerResponse

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(nil).Select(bson.M{"_id": idSlice, "timestamp": tsSlice}).Sort("timestamp").All(&reciever)
	if err != nil {
		//TODO handle this
		return out
	}

	max := len(reciever)
	top := max
	if top > 5 {
		top = 5
	}

	for i := 0; i < top; i++ {
		idSlice = append(idSlice, reciever[i].ID)
	}
	out = TickerResponse{
		TLatest:    reciever[max-1].Timestamp,
		TStart:     reciever[0].Timestamp,
		TStop:      reciever[top-1].Timestamp,
		TrackIDs:   idSlice,
		Processing: 0,
	}
	return out
}

//ReadSpecificTicker reads the ticker but with no input
func (db *MongoDB) ReadSpecificTicker(timeStamp int64) TickerResponse {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var idSlice []bson.ObjectId
	var tsSlice []int64

	type idTS struct {
		ID        bson.ObjectId
		Timestamp int64
	}

	var reciever []idTS

	var out TickerResponse

	err = session.DB(db.DatabaseName).C(db.TracksCollectionName).Find(nil).Select(bson.M{"_id": idSlice, "timestamp": tsSlice}).Sort("timestamp").All(&reciever)
	if err != nil {
		//TODO handle this
		return out
	}

	i := 0
	for done := false; done == false; {
		if reciever[i].Timestamp == timeStamp {
			done = true
		}
		i++
	}

	cap := i + 5

	for j := i; j < cap; j++ {
		idSlice = append(idSlice, reciever[j].ID)
	}
	out = TickerResponse{
		TLatest:    reciever[len(reciever)-1].Timestamp,
		TStart:     reciever[i].Timestamp,
		TStop:      reciever[cap].Timestamp,
		TrackIDs:   idSlice,
		Processing: 0,
	}
	return out
}
