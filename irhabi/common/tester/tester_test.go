// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tester

import (
	"net/http"
	"runtime"
	"testing"

	"git.qasico.com/cuxs/common/log"
	"github.com/buger/jsonparser"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var goVersion = runtime.Version()

type echoContent struct {
	Hello string `json:"hello"`
	Foo   string `json:"foo"`
	A     string `json:"a"`
	B     string `json:"b"`
	C     string `json:"c"`
	D     string `json:"d"`
}

// Binding from JSON
type echoJSONContent struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
}

func echoHelloHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echoContent{
			Hello: "world",
		})
	}
}

func echoTextHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func echoQueryHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		text := c.QueryParam("text")
		foo := c.QueryParam("foo")

		return c.JSON(http.StatusOK, &echoContent{
			Hello: text,
			Foo:   foo,
		})
	}
}

func echoPostFormHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		a := c.FormValue("a")
		b := c.FormValue("b")

		return c.JSON(http.StatusOK, &echoContent{
			A: a,
			B: b,
		})
	}
}

func echoJSONHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		json := new(echoJSONContent)
		err := c.Bind(json)

		if err != nil {
			log.Error(err)
		}

		return c.JSON(http.StatusOK, json)
	}
}

func echoPutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		foo := c.FormValue("c")
		bar := c.FormValue("d")

		return c.JSON(http.StatusOK, &echoContent{
			C: foo,
			D: bar,
		})
	}
}

func echoDeleteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echoContent{
			Hello: "world",
		})
	}
}

// EchoEngine is echo router.
func engine() *echo.Echo {
	e := echo.New()

	e.GET("/hello", echoHelloHandler())
	e.GET("/text", echoTextHandler())
	e.GET("/query", echoQueryHandler())

	e.POST("/form", echoPostFormHandler())
	e.POST("/json", echoJSONHandler())
	e.PUT("/update", echoPutHandler())
	e.DELETE("/delete", echoDeleteHandler())

	e.PATCH("/patch", echoHelloHandler())
	e.OPTIONS("/options", echoHelloHandler())
	e.HEAD("/head", echoHelloHandler())

	return e
}

func TestEchoHelloWorld(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetDebug(true).
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHeader(t *testing.T) {
	r := New()

	r.GET("/text").
		SetHeader(H{
			"Content-Type": "text/plain",
			"Go-Version":   goVersion,
		}).
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {

			assert.Equal(t, goVersion, rq.Header.Get("Go-Version"))
			assert.Equal(t, r.Body.String(), "Hello World")
			assert.Equal(t, r.Code, http.StatusOK)
		})
}

func TestEchoQuery(t *testing.T) {
	r := New()

	r.GET("/query?text=world&foo=bar").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")
			foo, _ := jsonparser.GetString(data, "foo")

			assert.Equal(t, "world", hello)
			assert.Equal(t, "bar", foo)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostFormData(t *testing.T) {
	r := New()

	r.POST("/form").
		SetBody("a=1&b=2").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetString(data, "a")
			b, _ := jsonparser.GetString(data, "b")

			assert.Equal(t, "1", a)
			assert.Equal(t, "2", b)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPostJSONData(t *testing.T) {
	r := New()

	r.POST("/json").
		SetJSON(D{
			"a": 1,
			"b": 2,
		}).
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			a, _ := jsonparser.GetInt(data, "a")
			b, _ := jsonparser.GetInt(data, "b")

			assert.Equal(t, 1, int(a))
			assert.Equal(t, 2, int(b))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPut(t *testing.T) {
	r := New()

	r.PUT("/update").
		SetBody("c=1&d=2").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			c, _ := jsonparser.GetString(data, "c")
			d, _ := jsonparser.GetString(data, "d")

			assert.Equal(t, "1", c)
			assert.Equal(t, "2", d)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoDelete(t *testing.T) {
	r := New()

	r.DELETE("/delete").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			hello, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", hello)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoPatch(t *testing.T) {
	r := New()

	r.PATCH("/patch").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoHead(t *testing.T) {
	r := New()

	r.HEAD("/head").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestEchoOptions(t *testing.T) {
	r := New()

	r.OPTIONS("/options").
		Run(engine(), func(r HTTPResponse, rq HTTPRequest) {
			data := []byte(r.Body.String())

			value, _ := jsonparser.GetString(data, "hello")

			assert.Equal(t, "world", value)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestSetQueryString(t *testing.T) {
	r := New()

	r.GET("/hello").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?a=1&b=2", r.Path)

	r.GET("/hello?foo=bar").
		SetQuery(H{
			"a": "1",
			"b": "2",
		})

	assert.Equal(t, "/hello?foo=bar&a=1&b=2", r.Path)
}
