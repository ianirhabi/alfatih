// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeDateRange(t *testing.T) {

	start, _ := time.Parse(time.RFC3339, "2018-01-01T00:00:00+07:00")
	end, _ := time.Parse(time.RFC3339, "2018-03-15T00:00:00+07:00")

	x := NewDateRange(start, end)

	diff := end.Sub(start)
	numberOfDays := int(diff.Hours()/24) + 1

	assert.Equal(t, numberOfDays, len(x.ByDate()))
	assert.Equal(t, 3, len(x.ByMonth()))
}
