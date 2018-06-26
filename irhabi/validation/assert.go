// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"git.qasico.com/cuxs/common"
)

// IsNotEmpty returns true if value is not nill
func IsNotEmpty(value interface{}) bool {
	if value == nil {
		return false
	}

	if str, ok := value.(string); ok {
		return len(str) > 0
	}
	if _, ok := value.(bool); ok {
		return true
	}
	if i, ok := value.(int); ok {
		return i != 0
	}
	if i, ok := value.(uint); ok {
		return i != 0
	}
	if i, ok := value.(int8); ok {
		return i != 0
	}
	if i, ok := value.(uint8); ok {
		return i != 0
	}
	if i, ok := value.(int16); ok {
		return i != 0
	}
	if i, ok := value.(uint16); ok {
		return i != 0
	}
	if i, ok := value.(uint32); ok {
		return i != 0
	}
	if i, ok := value.(int32); ok {
		return i != 0
	}
	if i, ok := value.(int64); ok {
		return i != 0
	}
	if i, ok := value.(uint64); ok {
		return i != 0
	}
	if t, ok := value.(time.Time); ok {
		tt := time.Time{}
		return (!t.IsZero() && t != tt)
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

// IsNumeric check if the value contains only numbers.
func IsNumeric(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	_, err := strconv.ParseFloat(str, 64)

	return err == nil
}

// IsAlpha check if the value contains only letters (a-zA-Z). Empty string is valid.
func IsAlpha(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	return AlphaPattern.MatchString(str)
}

// IsAlphanumeric check if the value contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	return AlphanumericPattern.MatchString(str)
}

// IsAlphanumericSpace check if the value contains only letters, numbers and space. Empty string is valid.
func IsAlphanumericSpace(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	return AlphanumericSpacePattern.MatchString(str)
}

// IsAlphaSpace check if the value contains only letters and space. Empty string is valid.
func IsAlphaSpace(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	return AlphaSpacePattern.MatchString(str)
}

// IsEmail check if the value is an email.
func IsEmail(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}

	return EmailPattern.MatchString(common.ToString(value))
}

// IsURL check if the value is an URL.
func IsURL(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return URLPattern.MatchString(str)
}

// IsJSON check if the value is valid JSON (note: uses json.Unmarshal).
func IsJSON(value interface{}) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(common.ToString(value)), &js) == nil
}

// IsCreditCard check if the value is a credit card.
func IsCreditCard(value interface{}) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	r, _ := regexp.Compile("[^0-9]+")
	sanitized := r.ReplaceAll([]byte(str), []byte(""))
	if !CreditCardPattern.MatchString(string(sanitized)) {
		return false
	}
	var sum int
	var digit string
	var tmpNum int
	var shouldDouble bool
	for i := len(sanitized) - 1; i >= 0; i-- {
		digit = string(sanitized[i:(i + 1)])
		tmpNum = common.ToInt(digit)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += ((tmpNum % 10) + 1)
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}

	if sum%10 == 0 {
		return true
	}
	return false
}

// IsLowerThanEqual return true if value is greather than equal given number
// this will evaluate value of int, lenght of string and number of slices.
func IsLowerThanEqual(value interface{}, max interface{}) (res bool) {
	if value == nil {
		return true
	}
	v := ChangeType(value)
	m := ChangeType(max)
	res = v <= m
	return res
}

// IsGreaterThanEqual return true if value is greather than equal given number
// this will evaluate value of int, lenght of string and number of slices.
func IsGreaterThanEqual(value interface{}, min interface{}) (res bool) {
	if value == nil {
		return true
	}
	v := ChangeType(value)
	m := ChangeType(min)
	res = v >= m
	return res
}

// IsLowerThan return true if value is lower than given number
// this will evaluate value of int, lenght of string and number of slices.
func IsLowerThan(value interface{}, max interface{}) (res bool) {
	if value == nil {
		return true
	}
	v := ChangeType(value)
	m := ChangeType(max)
	res = v < m
	return res
}

// IsGreaterThan return true if value is greather than given number
// this will evaluate value of int, lenght of string and number of slices.
func IsGreaterThan(value interface{}, min interface{}) (res bool) {
	if value == nil {
		return true
	}
	v := ChangeType(value)
	m := ChangeType(min)
	res = v > m
	return res
}

// ChangeType convert type of data to float64
func ChangeType(val interface{}) (x float64) {
	v := reflect.ValueOf(val)
	// check if type of data is pointer
	if v.Kind() == reflect.Ptr {
		//v is value of *val
		v = v.Elem()
	}
	//switch base on type
	switch v.Kind() {
	case reflect.String:
		//count string length and change it to float64
		str := utf8.RuneCountInString(common.ToString(v))
		x = float64(str)
		break
	case reflect.Slice:
		//length of slice to float64 (by Len() from lib value
		slc := v.Len()
		x = float64(slc)
		break
	case reflect.Float32:
		fl32 := val.(float32)
		x = float64(fl32)
		break
	case reflect.Float64:
		x = val.(float64)
		break
	default:
		num := common.ToInt(v)
		x = float64(num)
	}
	return x
}

// IsOnRange return true if value is greather than equal given min and lowerthan than equal given max
// this will evaluate value of int, lenght of string and number of slices.
func IsOnRange(value interface{}, min interface{}, max interface{}) bool {
	return IsGreaterThanEqual(value, min) && IsLowerThanEqual(value, max)
}

// IsContains check if the value contains the substring.
func IsContains(value interface{}, substring string) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	return strings.Contains(common.ToString(value), substring)
}

// IsMatches check if value matches the pattern (pattern is regular expression)
// In case of error return false
func IsMatches(value interface{}, pattern string) bool {
	str := common.ToString(value)
	if !IsNotEmpty(str) {
		return true
	}
	match, _ := regexp.MatchString(pattern, common.ToString(value))
	return match
}

// IsSame check if the value is identicaly same with given param
func IsSame(value interface{}, param interface{}) bool {
	value = common.ToString(value)
	if !IsNotEmpty(value) {
		return true
	}
	return value == common.ToString(param)
}

// IsIn check if the value is exists in given param
func IsIn(value interface{}, param ...string) bool {
	value = common.ToString(value)
	if !IsNotEmpty(value) {
		return true
	}
	if len(param) > 0 {
		for _, v := range param {
			if v == value {
				return true
			}
		}
	}
	return false
}

// IsNotIn check if the value is not exists in given param
func IsNotIn(value interface{}, param ...string) bool {
	value = common.ToString(value)
	if !IsNotEmpty(value) {
		return true
	}
	if len(param) > 0 {
		for _, v := range param {
			if v == value {
				return false
			}
		}
	}
	return true
}

// IsValidEan check if the value is valid EAN13 code
func IsValidEan(value interface{}) bool {
	ean := common.ToString(value)
	size := 13

	if len(ean) != size {
		return false
	}

	code := ean[:size-1]
	multiplyWhenEven := size%2 == 0
	sum := 0

	for i, v := range code {
		value, err := strconv.Atoi(string(v))
		if err != nil {
			return false
		}

		if (i%2 == 0) == multiplyWhenEven {
			sum += 3 * value
		} else {
			sum += value
		}
	}

	// checksum should
	return strconv.Itoa((10-sum%10)%10) == ean[size-1:size]
}
