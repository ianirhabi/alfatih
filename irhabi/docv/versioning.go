// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docv

import (
	"time"

	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

// Document structure of document version
type Document struct {
	Type      string      `bson:"type" json:"-"`
	ID        int64       `bson:"id" json:"-"`
	Version   int         `bson:"v" json:"version"`
	Data      interface{} `bson:"data" json:"-"`
	UpdatedBy interface{} `bson:"updated_by" json:"-"`
	StoredAt  time.Time   `bson:"stored_at" json:"stored_at"`

	DataPresented interface{} `bson:"-" json:"data,omitempty"`
	UserPresented interface{} `bson:"-" json:"updated_by,omitempty"`
}

// Save store a document and set to the new version
func (d *Document) Save() error {
	var lastV int

	// dapatkan last version berdasarkan type dan id document.
	if prev, e := Show(d.Type, d.ID); e == nil {
		lastV = prev.Version
	}
	// increast version
	d.Version = lastV + 1

	// save ke mongo
	return collection().Insert(d)
}

// Create initial function to store data into document
func Create(doc string, id int64, data interface{}, user interface{}) error {
	// should convert to json first so the version
	// will be the same format as real document
	jdata, _ := json.Marshal(data)
	juser, _ := json.Marshal(user)

	d := &Document{
		Type:      doc,
		ID:        id,
		Data:      string(jdata),
		UpdatedBy: string(juser),
		StoredAt:  time.Now(),
	}

	return d.Save()
}

// Show get detail of matching document,
// if version is not set will get the latest version
func Show(doc string, id int64, v ...int) (*Document, error) {
	bs := bson.M{"type": doc, "id": id}

	if len(v) > 0 && v[0] != 0 {
		bs["v"] = v[0]
	}

	var d Document
	var e error
	if e = collection().Find(bs).Sort("-$natural").Limit(1).One(&d); e == nil {
		json.Unmarshal([]byte(d.Data.(string)), &d.DataPresented)
		json.Unmarshal([]byte(d.UpdatedBy.(string)), &d.UserPresented)
	}

	return &d, e
}

// Get return all version of matching document.
func Get(doc string, id int64) ([]*Document, error) {
	bs := bson.M{"type": doc, "id": id}

	var d []*Document
	e := collection().Find(bs).Sort("-v").All(&d)

	return d, e
}

// Clean delete all version of matching document.
func Clean(doc string, id int64) error {
	bs := bson.M{"type": doc, "id": id}
	_, e := collection().RemoveAll(bs)
	return e
}
