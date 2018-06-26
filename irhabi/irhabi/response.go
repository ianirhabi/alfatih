package irhabi

import (
	"net/http"

	"github.com/alfatih/irhabi/orm"
	"github.com/alfatih/irhabi/validation"
	"github.com/labstack/echo"
)

const (
	// HTTPResponseSuccess default status for success responses
	HTTPResponseSuccess = "success"

	// HTTPResponseFail default status when responses has an errors.
	HTTPResponseFail = "fail"
)

// ResponseFormat is standart response formater of the applicatin.
type ResponseFormat struct {
	Code    int               `json:"-"`
	Status  string            `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// NewResponse return new instances of response formater.
func NewResponse() *ResponseFormat {
	return &ResponseFormat{
		Status: HTTPResponseSuccess,
	}
}

// SetData fill data and total into response formater.
func (r *ResponseFormat) SetData(d interface{}, t ...int64) *ResponseFormat {
	r.Status = HTTPResponseSuccess
	r.Code = http.StatusOK
	r.Data = d

	if len(t) > 0 {
		r.Total = t[0]
	}

	return r
}

// SetError set an error into response formater.
func (r *ResponseFormat) SetError(err error) *ResponseFormat {
	r.Code = http.StatusBadRequest
	r.Status = HTTPResponseFail
	r.Message = err.Error()

	// Check error based on type
	if he, ok := err.(*echo.HTTPError); ok {
		// Error cause of http failure should return status as is the errors
		// using standart http code.
		r.Code = he.Code
		r.Message = http.StatusText(r.Code)

		// if debug mode error messages will
		// send as response message.
		if IsDebug() {
			r.Message = he.Message
		}
	} else if o, ok := err.(*validation.Output); ok {
		// Error cause of validation failure should return
		// status 422 and returning all failure messages as errors.
		r.Code = http.StatusUnprocessableEntity
		r.Errors = o.Messages()
		r.Message = http.StatusText(r.Code)
	} else if oe, ok := err.(*orm.OrmError); ok {
		// Error cause of error from databases.
		// this should be treaten as bad requests.
		// but for debuging we need to know what the
		// error messages.
		r.Code = http.StatusBadRequest
		r.Message = http.StatusText(r.Code)

		// if debug mode error messages will
		// send as response message.
		if IsDebug() {
			r.Message = oe.Message
		}
	} else if ne, ok := err.(*DataNotExistsError); ok {
		// Error cause of data not exists or invalid
		// threated as Unprocessable Entity with code 422.
		r.Code = http.StatusUnprocessableEntity
		r.Errors = ne.Errors
		r.Message = http.StatusText(r.Code)
	} else if ne, ok := err.(*DataDuplicateError); ok {
		// Error cause of data not exists or invalid
		// threated as Unprocessable Entity with code 422.
		r.Code = http.StatusUnprocessableEntity
		r.Errors = ne.Errors
		r.Message = http.StatusText(r.Code)
	}

	return r
}

// reset all data in response formater
func (r *ResponseFormat) reset() {
	r.Data = nil
	r.Errors = nil
	r.Message = nil
	r.Total = 0
}
