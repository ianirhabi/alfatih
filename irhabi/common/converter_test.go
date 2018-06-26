// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	tests := []string{"1000", "-123", "abcdef", "100000000000000000000000000000000000000000000", " 1"}
	expected := []int{1000, -123, 0, 0, 1}
	for i := 0; i < len(tests); i++ {
		assert.Equal(t, expected[i], ToInt(tests[i]))
	}
}

func TestToBoolean(t *testing.T) {
	tests := []string{"true", "1", "True", "false", "0", "abcdef"}
	expected := []bool{true, true, true, false, false, false}
	for i := 0; i < len(tests); i++ {
		assert.Equal(t, expected[i], ToBoolean(tests[i]))
	}
}

func toString(t *testing.T, test interface{}, expected string) {
	assert.Equal(t, expected, ToString(test))
}

func TestToString(t *testing.T) {
	toString(t, "str123", "str123")
	toString(t, 123, "123")
	toString(t, 12.3, "12.3")
	toString(t, true, "true")
	toString(t, 1.5+10i, "(1.5+10i)")
}

func TestToFloat(t *testing.T) {
	tests := []interface{}{"", "123", "-.01", "10.", "string", "1.23e3", ".23e10", []string{"asd"}, 0.1}
	expected := []float64{0, 123, -0.01, 10.0, 0, 1230, 0.23e10, 0, 0.1}
	for i := 0; i < len(tests); i++ {
		assert.Equal(t, expected[i], ToFloat(tests[i]))
	}
}

func TestToJSON(t *testing.T) {
	tests := []interface{}{"test", map[string]string{"a": "b", "b": "c"}, func() error {
		return fmt.Errorf("Error")
	}}
	expected := [][]string{
		[]string{"\"test\"", "<nil>"},
		[]string{"{\"a\":\"b\",\"b\":\"c\"}", "<nil>"},
		[]string{"", "json: unsupported type: func() error"},
	}
	for i, test := range tests {
		assert.Equal(t, expected[i][0], ToJSON(test))
	}
}

func TestToCamelCase(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "ABC"},
		{"my_func", "MyFunc"},
		{"1ab_cd", "1abCd"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ToCamelCase(test.param))
	}

	tests = []struct {
		param    string
		expected string
	}{
		{"abc", "abc"},
		{"a_b_c", "aBC"},
		{"my_func", "myFunc"},
		{"1ab_cd", "1abCd"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ToLowerCamelCase(test.param))
	}
}

func TestToUnderscore(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"MyFunc", "my_func"},
		{"ABC", "a_b_c"},
		{"1B", "1_b"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ToUnderscore(test.param))
	}
}

func TestToUpper(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "A_B_C"},
		{"my_func", "MY_FUNC"},
		{"1ab_cd", "1AB_CD"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ToUpper(test.param))
	}
}

func TestToLower(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "A_B_C"},
		{"my_func", "MY_FUNC"},
		{"1ab_cd", "1AB_CD"},
	}
	for _, test := range tests {
		assert.Equal(t, test.param, ToLower(test.expected))
	}
}

func TestLeftTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "foo  \r\n\t   "},
		{"010100201000", "01", "201000"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, LeftTrim(test.param1, test.param2))
	}
}

func TestRightTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "  \r\n\tfoo"},
		{"010100201000", "01", "0101002"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, RightTrim(test.param1, test.param2))
	}
}

func TestTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "foo"},
		{"010100201000", "01", "2"},
		{"1234567890987654321", "1-8", "909"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, Trim(test.param1, test.param2))
	}
}
