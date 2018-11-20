package storage

import (
	mgo "github.com/globalsign/mgo"
	bson "github.com/globalsign/mgo/bson"
)

//WebhookDB ...
var WebhookDB WebhookStorage

//WebhookStorage ...
type WebhookStorage interface {
	AddWebhook(inFile WebhookInfo) string
	CountWebhooks() int
	DeleteWebhook(id string) WebhookInfo
	ReadHookIDS() []string
	ReadWebhook(id string) WebhookInfo
	CheckWebhooks() []WebhookInfo
}

//AddWebhook ...
func (db *MongoDBWebHook) AddWebhook(inFile WebhookInfo) string {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	inFile.ID = bson.NewObjectId()
	inFile.LastTrackCount = TrackDB.CountTracks()
	id := inFile.ID.Hex()

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Insert(inFile)
	if err != nil {
		//!Handle me
		return ""
	}

	return id
}

//CountWebhooks ...
func (db *MongoDBWebHook) CountWebhooks() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	out, err := session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Count()
	if err != nil {
		//!Handle me
		panic(err)
	}
	return out
}

//DeleteWebhook ...
func (db *MongoDBWebHook) DeleteWebhook(id string) WebhookInfo {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	requestID := bson.ObjectIdHex(id)
	var response WebhookInfo

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{"_id": requestID}).One(&response)
	if err != nil {
		//!Handle me
		panic(err)
	}

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Remove(bson.M{"_id": requestID})
	if err != nil {
		//!Handle me
		panic(err)
	}

	return response
}

//ReadHookIDS ...
func (db *MongoDBWebHook) ReadHookIDS() []string {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var tempHooks []WebhookInfo

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{}).All(&tempHooks)
	if err != nil {
		//!Handle me
		return make([]string, 0)
	}

	var idSlice []string
	for i := 0; i < len(tempHooks); i++ {
		idSlice = append(idSlice, tempHooks[i].ID.Hex())
	}

	return idSlice
}

//ReadWebhook ...
func (db *MongoDBWebHook) ReadWebhook(id string) WebhookInfo {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	requestID := bson.ObjectIdHex(id)
	var response WebhookInfo

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{"_id": requestID}).One(&response)
	if err != nil {
		//!Handle me
		panic(err)
	}

	return response
}

//CheckWebhooks ...
func (db *MongoDBWebHook) CheckWebhooks() []WebhookInfo {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allHooks []WebhookInfo

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{}).All(&allHooks)
	if err != nil {
		//!Handle me
		panic(err)
	}

	var response []WebhookInfo

	for key := range allHooks {
		if allHooks[key].TriggerValue == 1 {
			response = append(response, allHooks[key])
		} else {
			if TrackDB.CountTracks()-allHooks[key].LastTrackCount == allHooks[key].TriggerValue {
				response = append(response, allHooks[key])
				updateID := allHooks[key].ID
				err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Update(bson.M{"_id": updateID}, bson.M{"lastTrackCount": TrackDB.CountTracks()})
				if err != nil {
					//!Handle me
					panic(err)
				}
			}
		}
	}

	return response
}
