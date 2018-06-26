// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ToString convert the input to a string.
func ToString(value interface{}) string {
	res := fmt.Sprintf("%v", value)
	return string(res)
}

// ToJSON convert the input to a valid JSON string
func ToJSON(value interface{}) string {
	res, err := json.Marshal(value)
	if err != nil {
		res = []byte("")
	}
	return string(res)
}

// ToFloat convert the input string to a float, or 0.0 if the input is not a float.
func ToFloat(value interface{}) float64 {
	floatType := reflect.TypeOf(float64(0))

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		res, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			res = 0.0
		}

		return res
	}

	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0
	}

	return v.Convert(floatType).Float()
}

// ToInt convert the input string to an integer, or 0 if the input is not an integer.
func ToInt(value interface{}) int {
	res, err := strconv.Atoi(Trim(ToString(value), ""))
	if err != nil {
		res = 0
	}
	return res
}

// ToBoolean convert the input string to a boolean.
func ToBoolean(value interface{}) bool {
	res, err := strconv.ParseBool(ToString(value))
	if err != nil {
		res = false
	}
	return res
}

// ToLower convert the value string into lowercase format.
func ToLower(value interface{}) string {
	return strings.ToLower(ToString(value))
}

// ToUpper convert the value string into uppercase format.
func ToUpper(value interface{}) string {
	return strings.ToUpper(ToString(value))
}

// ToCamelCase converts from underscore separated form to camel case form.
// Ex.: my_func => MyFunc
func ToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// ToLowerCamelCase converts from underscore separated form to lower camel case form.
// Ex.: my_func => myFunc
func ToLowerCamelCase(s string) string {
	a := []rune(ToCamelCase(s))
	if len(a) > 0 {
		a[0] = unicode.ToLower(a[0])
	}
	return string(a)
}

// ToUnderscore converts from camel case form to underscore separated form.
// Ex.: MyFunc => my_func
func ToUnderscore(str string) string {
	var output []rune
	var segment []rune
	for _, r := range str {
		if !unicode.IsLower(r) {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

// LeftTrim trim characters from the left-side of the input.
// If second argument is empty, it's will be remove leading spaces.
func LeftTrim(str, chars string) string {
	pattern := "^\\s+"
	if chars != "" {
		pattern = "^[" + chars + "]+"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// RightTrim trim characters from the right-side of the input.
// If second argument is empty, it's will be remove spaces.
func RightTrim(str, chars string) string {
	pattern := "\\s+$"
	if chars != "" {
		pattern = "[" + chars + "]+$"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// Trim trim characters from both sides of the input.
// If second argument is empty, it's will be remove spaces.
func Trim(str, chars string) string {
	return LeftTrim(RightTrim(str, chars), chars)
}

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}
