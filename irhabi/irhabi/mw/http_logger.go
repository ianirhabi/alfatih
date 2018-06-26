package mw

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alfatih/irhabi/common/log"
	"github.com/labstack/echo"
)

// HTTPLogger returns a middleware that logs HTTP requests.
func HTTPLogger() echo.MiddlewareFunc {
	return func(n echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return logRequest(n, c)
		}
	}
}

// logRequest print all http request on consoles.
func logRequest(hand echo.HandlerFunc, c echo.Context) (err error) {
	start := time.Now()
	req := c.Request()
	res := c.Response()
	if err = hand(c); err != nil {
		c.Error(err)
	}
	end := time.Now()
	latency := end.Sub(start) / 1e5

	log := log.Log.WithFields(logrus.Fields{
		"prefix":    fmt.Sprintf("%s/%d", req.Method, res.Status),
		"latecy":    fmt.Sprintf("%1.1fms", float64(int(latency))/10.0),
		"requester": req.RemoteAddr,
		"path":      req.URL.Path,
	})

	if err == nil {
		log.Info(http.StatusText(res.Status))
	} else {
		log.Error(err)
	}
	return
}
