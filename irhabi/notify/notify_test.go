// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package notify_test

import (
	"fmt"
	"os"
	"testing"

	"git.qasico.com/cuxs/notify"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	// run tests
	res := m.Run()

	os.Exit(res)
}

func TestNotify(t *testing.T) {
	// initialing notify session
	er := notify.New()
	if er != nil {
		panic(er.Error())
	}

	defer notify.Session.Close()

	// test Create
	var cases = []struct {
		Data     *notify.ObjectData
		Expected bool
	}{
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "65536", Action: "google"}, ID: bson.NewObjectId(), UserID: 1, Title: "Testing Notification", DeviceID: "xxxx-xxxx-xxxx-xxxx-xxxx", ActionURL: "http://google.com/a/qasico.com", Message: "Notification Messages"}, true},
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "65536", Action: "google"}, ID: bson.NewObjectId(), UserID: 1, ActionURL: "http://google.com/a/qasico.com", Message: "Notification 2"}, true},
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "131072", Action: "google"}, ID: bson.NewObjectId(), UserID: 0, DeviceID: "xxxx", ActionURL: "http://google.com/a/qasico.com", Message: "Notification 3"}, false},
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "65536", Action: "google"}, ID: bson.NewObjectId(), UserID: 1, Message: "Notification 4"}, false},
		{&notify.ObjectData{ID: bson.NewObjectId(), UserID: 1}, false},
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "131072", Action: "google"}, ID: bson.NewObjectId(), UserID: 2, Title: "Testing Notification", DeviceID: "xxxx-xxxx-xxxx-xxxx-xxx1", ActionURL: "http://google.com/a/qasico.com", Message: "Notification Messages"}, true},
		{&notify.ObjectData{ObjectAction: &notify.ObjectAction{ID: "196608", Action: "google"}, ID: bson.NewObjectId(), UserID: 3, Title: "Testing Notification", DeviceID: "xxxx-xxxx-xxxx-xxxx-xxx2", ActionURL: "http://google.com/a/qasico.com", Message: "Notification Messages"}, true},
	}

	for _, c := range cases {
		err := notify.Create(c.Data)
		if !assert.Equal(t, c.Expected, err == nil) {
			fmt.Println(c.Data, c.Expected)
		}
	}

	var d []*notify.ObjectData
	var e error

	// test GetByUser
	d, _ = notify.GetByUser(4)
	assert.Nil(t, d, "should be nill when no data returned.")

	d, e = notify.GetByUser(1)
	assert.NoError(t, e, "should be no error.")
	assert.Equal(t, 2, len(d), "len data should be equal")
	assert.Equal(t, "65536", d[0].ObjectAction.ID)

	// test GetByDevice
	d, _ = notify.GetByDevice("yyy")
	assert.Nil(t, d, "should be nill when no data returned.")

	d, e = notify.GetByDevice("xxxx-xxxx-xxxx-xxxx-xxxx")
	assert.NoError(t, e, "should be no error.")
	assert.Equal(t, 1, len(d), "len data should be equal")

	// test Read
	e = notify.Read(1)
	assert.NoError(t, e)

	// test ReadByID
	// ambil data terlebih dahulu dengan user id = 2
	d, _ = notify.GetByUser(2)
	e = notify.ReadByID(d[0].ID.Hex())
	assert.NoError(t, e)
	// cek data tersebut harus readed
	dUser, _ := notify.GetByUser(2)
	assert.Equal(t, bool(true), dUser[0].Readed)

	//test UpdateByUserAndObjectAction
	d, _ = notify.GetByUser(3)
	// ambil data terlebih dahulu dengan user id = 3
	e = notify.UpdateByUserAndObjectAction(d[0].UserID, d[0].ObjectAction.ID, d[0].ObjectAction.Action)
	assert.NoError(t, e)
	// cek data tersebut harus readed
	dUser3, _ := notify.GetByUser(3)
	assert.Equal(t, bool(true), dUser3[0].Readed)

	// cleanup
	notify.Clean(1)
	notify.Clean(2)
	notify.Clean(3)
}
