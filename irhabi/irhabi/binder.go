package irhabi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/alfatih/irhabi/validation"
	"github.com/labstack/echo"
)

type (
	// Custom echo binder
	binder struct{}

	// Validator custom validation object.
	binderValidator struct {
		once      sync.Once
		validator *validation.Validation
	}
)

// Set binding validator using cuxs validation.
var bindValidator = &binderValidator{}

// ValidateStruct evaluate an object,
// will run validation request if the object
// is implementing validatonRequests.
func (v *binderValidator) validate(obj interface{}) error {
	v.lazyinit()

	var o *validation.Output
	if vr, ok := obj.(validation.Request); ok {
		o = v.validator.Request(vr)
	} else {
		o = v.validator.Struct(obj)
	}

	if !o.Valid {
		return o
	}
	return nil
}

// lazyinit initialing validator instances for one of time only.
func (v *binderValidator) lazyinit() {
	v.once.Do(func() {
		v.validator = validation.New()
	})
}

// Bind is decode request body and injecting into interfaces,
// We only accept json data type other type will return error bad requests.
// Also automaticly validate data with interfaces.
func (b binder) Bind(i interface{}, ctx echo.Context) (err error) {
	bindValidator.lazyinit()
	req := ctx.Request()
	ctype := req.Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		if err = json.NewDecoder(req.Body).Decode(i); err != nil && err != io.EOF {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("unmarshal type error: expected=%v, got=%v, offset=%v", ute.Type, ute.Value, ute.Offset))
			} else if se, ok := err.(*json.SyntaxError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("syntax error: offset=%v, error=%v", se.Offset, se.Error()))
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}
		return bindValidator.validate(i)
	}
	return echo.ErrUnsupportedMediaType
}
