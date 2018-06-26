package irhabi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.GET("/ok", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	c, b := request(echo.GET, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)

	c, b = request(echo.GET, "/404", e)
	assert.Equal(t, http.StatusNotFound, c)
	assert.Equal(t, `{"status":"fail","message":"Not Found"}`, b)
}

func TestEchoRoutes(t *testing.T) {
	e := New()
	routes := []echo.Route{
		{Method: echo.GET, Path: "/users/:user/events", Handler: ""},
		{Method: echo.GET, Path: "/users/:user/events/public", Handler: ""},
		{Method: echo.POST, Path: "/repos/:owner/:repo/git/refs", Handler: ""},
		{Method: echo.POST, Path: "/repos/:owner/:repo/git/tags", Handler: ""},
		{Method: echo.PUT, Path: "/repos/:owner/:repo/git/tags", Handler: ""},
	}
	for _, r := range routes {
		if r.Method == echo.GET {
			e.GET(r.Path, func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			})
		} else if r.Method == echo.POST {
			e.POST(r.Path, func(c echo.Context) error {
				return c.String(http.StatusOK, "OK")
			})
		}
	}

	for _, r := range e.Routes() {
		found := false
		for _, rr := range routes {
			if r.Method == rr.Method && r.Path == rr.Path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Route %s : %s not found", r.Method, r.Path)
		}
	}
}

func request(method, path string, e *echo.Echo) (int, string) {
	req, _ := http.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
