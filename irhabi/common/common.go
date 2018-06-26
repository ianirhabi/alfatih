// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric       = "0123456789"
)

// RandomStr return random string with defined length.
func RandomStr(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomNumeric return random numeric with defined length.
func RandomNumeric(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numeric) {
			b[i] = numeric[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// FloatPrecision make a float with precision decimal
func FloatPrecision(num float64, precision int) float64 {
	m := "1" + strings.Repeat("0", precision)
	mx, _ := strconv.Atoi(m)
	return float64(int(num*float64(mx))) / float64(mx)
}

// Rounder to round float number with some intermed and precision decimal
// e.g: 2.558--->2.6
// Rounder(2.558, 0.5, 1)--->2.6 (intermed is a number that determine range of round ceil and round floor
// can accept negatif number
// from:https://play.golang.org/p/KNhgeuU5sT
func Rounder(num float64, intermed float64, precision int) float64 {

	pow := math.Pow(10, float64(precision))
	digit := pow * num
	_, div := math.Modf(digit)

	var round float64
	if num > 0 {
		if div >= intermed {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if div >= intermed {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}
	return round / pow
}

// Encrypt perform simple encryption and decription values.
func Encrypt(n interface{}) string {
	num := ToInt(n)
	return ToString(((0x0000FFFF & num) << 16) + ((0xFFFF0000 & num) >> 16))
}

// Decrypt return real values of encripted values.
func Decrypt(v interface{}) (id int64, err error) {
	if num := ToInt(Encrypt(v)); num != 0 {
		id = int64(num)
		return id, nil
	}

	err = &DecryptionError{
		Values: v,
	}
	return
}

// DecryptionError error type caused by decription failure.
type DecryptionError struct {
	Message string
	Values  interface{}
}

// Error implement error interfaces.
func (e *DecryptionError) Error() string {
	return "Invalid encryption values: " + ToString(e.Values)
}

// PasswordHash compares hashed password with its possible
// plaintext equivalent using bcrypt algorithm.
func PasswordHash(hashed string, plain string) error {
	h := []byte(hashed)
	p := []byte(plain)

	return bcrypt.CompareHashAndPassword(h, p)
}

// PasswordHasher returns the bcrypt hash of the password
// using DefaultCost
func PasswordHasher(p string) (h string, err error) {
	if hx, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost); err == nil {
		h = string(hx)
	}
	return
}

// Contains cek is slice contains a strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
