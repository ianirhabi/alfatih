// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package faker

import (
	"math/rand"
	"reflect"
	"time"

	"regexp"
	"strconv"
	"strings"

	"github.com/alfatih/irhabi/common"
	"github.com/jmcvetta/randutil"
)

// Fill passed interface with random data based on the struct field type,
// take a look at fuzzValueFor for details on supported data types.
func Fill(e interface{}, except ...string) {
	v := reflect.TypeOf(e)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		value := reflect.ValueOf(e).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := value.Field(i)
			name := value.Type().Field(i).Name

			tag := v.Field(i).Tag.Get("orm")
			// cek field can be set and not in exception param
			if field.CanSet() && !common.Contains(except, name) {
				fillElement(field, tag)
			}
		}
	}
}

type OrmTag struct {
	Name    string
	Size    int
	Options []string
	Digit   int
	Decimal int
}

// ParseOrmTag parse struct tag string
func ParseOrmTag(data string) *OrmTag {
	tag := new(OrmTag)
	tag.Size = 25
	tag.Digit = 2
	tag.Decimal = 2
	for _, v := range strings.Split(data, ";") {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		t := strings.ToLower(v)
		if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i]
			opt := regexp.MustCompile(`\(([^)]+)\)`).FindStringSubmatch(v)
			val := opt[1]

			switch name {
			case "column":
				tag.Name = val
				break
			case "size":
				tag.Size, _ = strconv.Atoi(val)
				break
			case "options":
				tag.Options = strings.Split(val, ",")
				break
			case "digits":
				tag.Digit, _ = strconv.Atoi(val)
				break
			case "decimals":
				tag.Decimal, _ = strconv.Atoi(val)
				break
			default:
				continue
			}
		}
	}
	return tag
}

// fuzzValueFor Generates random values for the following types:
// string, bool, int, int32, int64, float32, float64
func fillElement(e reflect.Value, ormTag string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	vals := reflect.ValueOf(e)
	tag := ParseOrmTag(ormTag)

	rd := rand.Intn(tag.Digit)
	switch e.Kind() {
	case reflect.String:
		// get size
		if len(tag.Options) > 0 {
			v, _ := randutil.ChoiceString(tag.Options)
			vals = reflect.ValueOf(v)
		} else {
			vals = reflect.ValueOf(common.RandomStr(tag.Size))
		}
		break
	case reflect.Uint:
		vals = reflect.ValueOf(uint(rd))
		break
	case reflect.Uint8:
		vals = reflect.ValueOf(uint8(rd))
		break
	case reflect.Uint16:
		vals = reflect.ValueOf(uint16(rd))
		break
	case reflect.Uint32:
		vals = reflect.ValueOf(uint32(rd))
		break
	case reflect.Uint64:
		vals = reflect.ValueOf(uint64(rd))
		break
	case reflect.Int:
		rand.Seed(time.Now().Unix())
		vals = reflect.ValueOf(rand.Intn(999999999-1) + 1)
		break
	case reflect.Int8:
		vals = reflect.ValueOf(int8(rd))
		break
	case reflect.Int16:
		vals = reflect.ValueOf(int16(rd))
		break
	case reflect.Int32:
		vals = reflect.ValueOf(r.Int31())
		break
	case reflect.Int64:
		vals = reflect.ValueOf(r.Int63())
		break
	case reflect.Float32:
		f := common.FloatPrecision(r.Float64(), tag.Decimal)
		vals = reflect.ValueOf(float32(f))
		break
	case reflect.Float64:
		vals = reflect.ValueOf(common.FloatPrecision(r.Float64(), tag.Decimal))
		break
	case reflect.Bool:
		val := r.Intn(2) > 0
		vals = reflect.ValueOf(val)
		break
	default:
		if e.Type() == reflect.TypeOf(time.Time{}) {
			ft := time.Now().Add(time.Second * 3600)
			vals = reflect.ValueOf(ft)
		} else {
			vals = reflect.ValueOf(nil)
		}
	}

	if vals.IsValid() {
		e.Set(vals)
	}
}
