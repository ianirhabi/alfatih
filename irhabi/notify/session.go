// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package notify

import (
	"git.qasico.com/cuxs/env"
	"gopkg.in/mgo.v2"
)

// Config represents all configurable docv library
var Config *configNotify

// configDocv type to store docv configuration
type configNotify struct {
	Host       string
	Database   string
	Collection string
}

// Session instance of Mgo Session
var Session *mgo.Session

// New init set all configurable data from env variable on application started.
// and creating an session mongo db instances
func New() (e error) {
	Config = &configNotify{
		Host:       env.GetString("MGO_HOST", "172.17.0.3"),
		Database:   env.GetString("MGO_DB", "notify"),
		Collection: env.GetString("MGO_NOTIFY_COLLECTION", "notify_document"),
	}

	return session()
}

// session try to dial monggodb server
func session() (err error) {
	Session, err = mgo.Dial(Config.Host)
	return
}

// collection instance from mgo, reading from config
func collection() *mgo.Collection {
	return Session.DB(Config.Database).C(Config.Collection)
}
