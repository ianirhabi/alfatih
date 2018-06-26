// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    interface{}
		expected string
	}{
		{1, "65536"},
		{4040, "264765440"},
		{"264765440", "4040"},
		{"65536", "1"},
		{"randomstring", "0"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, Encrypt(test.param))
	}
}

func TestDecrypt(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected int64
		valid    bool
	}{
		{"264765440", 4040, true},
		{"65536", 1, true},
		{"randomstring", 0, false},
	}

	for _, test := range tests {
		v, e := Decrypt(test.param)
		if e != nil {
			_, ok := e.(*DecryptionError)
			if !ok {
				assert.Fail(t, "Error type is invalid")
			}
		}

		assert.Equal(t, test.expected, v)
		assert.Equal(t, e == nil, test.valid)
	}
}

func TestRandomStr(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		len int
	}{
		{10},
		{4},
		{1000},
		{5},
	}

	for _, test := range tests {
		str := RandomStr(test.len)
		assert.Equal(t, test.len, len(str), "Not equal length given.")
	}
}

func TestRandomNumeric(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		len int
	}{
		{10},
		{4},
		{1000},
		{5},
	}

	for _, test := range tests {
		str := RandomNumeric(test.len)
		assert.Equal(t, test.len, len(str), "Not equal length given.")
	}
}

func TestPasswordHash(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param string
	}{
		{"randompassword"},
		{"123456"},
	}
	for _, test := range tests {
		hash, e := PasswordHasher(test.param)
		assert.NoError(t, e)

		match := PasswordHash(hash, test.param)
		assert.NoError(t, match)
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		slicer   []string
		value    string
		expected bool
	}{
		{[]string{"one", "two"}, "one", true},
		{[]string{"one", "two"}, "three", false},
	}

	for _, test := range tests {
		res := Contains(test.slicer, test.value)
		assert.Equal(t, test.expected, res)
	}
}

func TestFloatPrecision(t *testing.T) {
	d := 10.123123123
	r := FloatPrecision(d, 2)
	assert.Equal(t, 10.12, r)
}

func TestRounder(t *testing.T) {
	d := Rounder(2.588, 0.5, 1)
	assert.Equal(t, float64(2.6), d)

	d2 := Rounder(2.588, 0.5, 0)
	assert.Equal(t, float64(3), d2)

	d3 := Rounder(2.588, 0.9, 1)
	assert.Equal(t, float64(2.5), d3)

	d4 := Rounder(-38.288888, 0.5, 2)
	assert.Equal(t, float64(-38.28), d4)

	d5 := Rounder(3.333333, 0.3, 1)
	assert.Equal(t, float64(3.4), d5)
}
