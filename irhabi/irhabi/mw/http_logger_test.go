package mw

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alfatih/irhabi/common/log"
	"github.com/labstack/echo"
)

func TestLogger(t *testing.T) {
	var w bytes.Buffer
	setup(&w)

	// Note: Just for the test coverage, not a real test.
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := HTTPLogger()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Status 2xx
	h(c)

	// Status 3xx
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = HTTPLogger()(func(c echo.Context) error {
		return c.String(http.StatusTemporaryRedirect, "test")
	})
	h(c)

	// Status 4xx
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = HTTPLogger()(func(c echo.Context) error {
		return c.String(http.StatusNotFound, "test")
	})
	h(c)

	// Status 5xx with empty path
	req, _ = http.NewRequest(echo.GET, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = HTTPLogger()(func(c echo.Context) error {
		return errors.New("error")
	})
	h(c)
}

func setup(w io.Writer) {
	log.Log.Out = w
}
