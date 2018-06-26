// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package notify

import (
	"time"

	"git.qasico.com/cuxs/validation"
	"gopkg.in/mgo.v2/bson"
)

// ObjectAction Document structure for action notification
type ObjectAction struct {
	ID     string `bson:"id"`
	Action string `bson:"action"`
}

// ObjectData Document structure of document version
type ObjectData struct {
	ID           bson.ObjectId `bson:"_id"`
	UserID       int64         `bson:"user_id" valid:"required|gt:0"`
	DeviceID     string        `bson:"device_id"`
	Title        string        `bson:"title"`
	Message      string        `bson:"message" valid:"required"`
	ActionURL    string        `bson:"action_url" valid:"required"`
	Readed       bool          `bson:"readed"`
	CreatedAt    time.Time     `bson:"created_at"`
	ReadedAt     time.Time     `bson:"readed_at"`
	ObjectAction *ObjectAction `bson:"object_action"`
}

// Create saving new notification
func Create(d *ObjectData) (e error) {
	v := validation.New()

	if o := v.Struct(d); !o.Valid {
		return o
	}

	d.CreatedAt = time.Now()
	d.Readed = false

	e = collection().Insert(d)
	if e == nil {
		e = notify(d)
	}

	return
}

// GetByUser get all notification by user_id
func GetByUser(UserID int64) (d []*ObjectData, e error) {
	bs := bson.M{"user_id": UserID}

	e = collection().Find(bs).Sort("-created_at").All(&d)
	return
}

// GetByDevice get all notification by device_id
func GetByDevice(DeviceID string) (d []*ObjectData, e error) {
	bs := bson.M{"device_id": DeviceID}

	e = collection().Find(bs).Sort("-created_at").All(&d)

	return
}

// Read update status notification into readed
func Read(UserID int64) (e error) {
	bs := bson.M{"user_id": UserID, "readed": false}
	data := bson.M{"readed": true, "readed_at": time.Now()}

	_, e = collection().UpdateAll(bs, bson.M{"$set": data})

	return
}

// ReadByID update status notification into readed by id notify
func ReadByID(ID string) (e error) {
	objectID := bson.ObjectIdHex(ID)

	bs := bson.M{"_id": objectID, "readed": false}
	data := bson.M{"readed": true, "readed_at": time.Now()}

	e = collection().Update(bs, bson.M{"$set": data})

	return
}

// UpdateByUserAndObjectAction get notification by user id, object action id dan action
// update status notification into readed by id notify
func UpdateByUserAndObjectAction(UserID int64, ObjectActionID string, Action string) (e error) {
	var data ObjectData
	bs := bson.M{"user_id": UserID, "object_action": bson.M{"id": ObjectActionID, "action": Action}}

	if e = collection().Find(bs).One(&data); e == nil {
		e = ReadByID(data.ID.Hex())
		return e
	}
	return e
}

// Clean delete all version of matching document.
func Clean(UserID int64) error {
	bs := bson.M{"user_id": UserID}
	_, e := collection().RemoveAll(bs)
	return e
}
