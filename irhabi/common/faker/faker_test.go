// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package faker_test

import (
	"testing"

	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/faker"
	"github.com/stretchr/testify/assert"
)

type ExtSimple struct {
	Value string
}

type Simple struct {
	ValueString            string
	ValueUint              uint
	ValueUint8             uint8
	ValueUint16            uint16
	ValueUint32            uint32
	ValueUint64            uint64
	Value                  int
	Value8                 int8
	Value16                int16
	Value32                int32
	Value64                int64
	ValueF32               float32
	ValueF64               float64
	ValueBool              bool
	ValueTime              time.Time
	ValuePtr               *ExtSimple
	ValueSlice             []*ExtSimple
	ValuePtrInt            *int
	ValuePtrTime           *time.Time
	ValueSliceString       []string `json:"name"`
	ValueWithTagOrmSize    string   `orm:"column(employee_addr);size(3)"`
	ValueWithTagOrmOptions string   `orm:"column(religion);options(islam,katolik,kristen)"`
}

func TestFill(t *testing.T) {
	d := Simple{}

	faker.Fill(&d)
	assert.NotEqual(t, 0, d.Value)
	assert.NotEqual(t, 0, d.Value8)
	assert.NotEqual(t, 0, d.Value16)
	assert.NotEqual(t, 0, d.Value32)
	assert.NotEqual(t, 0, d.Value64)
	assert.NotEqual(t, 0.0, d.ValueF32)
	assert.NotEqual(t, 0.0, d.ValueF64)
	assert.Nil(t, d.ValuePtr)
	assert.NotEmpty(t, d.ValueTime)
	assert.Nil(t, d.ValuePtrInt)
	assert.Equal(t, 0, len(d.ValueSlice))
	assert.Equal(t, 0, len(d.ValueSliceString))
	assert.Equal(t, 3, len(d.ValueWithTagOrmSize))
	assert.Condition(t, func() bool {
		return common.Contains([]string{"islam", "katolik", "kristen"}, d.ValueWithTagOrmOptions)
	})

	count := 0
	for i := 0; i < 10; i++ {
		faker.Fill(&d)
		if d.ValueBool {
			count++
		}
	}

	assert.True(t, count > 0)
}

func TestFillExcept(t *testing.T) {
	d := Simple{}

	faker.Fill(&d, "Value")
	assert.Empty(t, d.Value)
}

func TestParseOrmTag(t *testing.T) {
	tag := faker.ParseOrmTag("column(religion);options(islam,katolik,kristen);size(3)")

	assert.Equal(t, "religion", tag.Name)
	assert.Len(t, tag.Options, 3)
	assert.Equal(t, 3, tag.Size)
}
