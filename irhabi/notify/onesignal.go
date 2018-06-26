// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package notify

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/env"
)

const (
	// ONESIGNALURL endpoint apis from onesignal
	ONESIGNALURL = "https://onesignal.com/api/v1/notifications"
)

type (
	// OnesignalAdaptor struct data will be passed to onesignal api
	OnesignalAdaptor struct {
		AppID     string            `json:"app_id"`
		Players   []string          `json:"include_player_ids"`
		SmallIcon string            `json:"small_icon"`
		LargeIcon string            `json:"large_icon"`
		Contents  *OnesignalContent `json:"contents"`
		Header    *OnesignalHeader  `json:"headings"`
		Data      *OnesignalData    `json:"data"`
		URL       string            `json:"url"`
	}

	// OnesignalData additional data, containts callback url
	OnesignalData struct {
		URL string `json:"url"`
	}

	// OnesignalContent containts the content notification
	OnesignalContent struct {
		EN string `json:"en"`
		ID string `json:"id"`
	}

	// OnesignalHeader containts the title notification
	OnesignalHeader struct {
		EN string `json:"en"`
		ID string `json:"id"`
	}
)

func notify(data *ObjectData) error {
	ones := new(OnesignalAdaptor)
	ones.SmallIcon = "ic_stat_onesignal_default"
	ones.LargeIcon = "ic_stat_onesignal_default"

	ones.AppID = env.GetString("ONES_KEY", "eaf26238-dd12-44ec-9b42-4e98972962b0")
	ones.Players = append(ones.Players, data.DeviceID)
	ones.Contents = &OnesignalContent{EN: data.Message}
	ones.Header = &OnesignalHeader{EN: data.Title}
	ones.Data = &OnesignalData{URL: data.ActionURL}

	body, e := json.Marshal(ones)
	r, e := http.Post(ONESIGNALURL, "application/json", bytes.NewBuffer(body))
	if e != nil {
		log.Debug("Failed to dispatch request to %s \n %s", ONESIGNALURL, e.Error())
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	log.Debug("Push notification was dispatched: %v", string(bodyBytes))

	return e
}
