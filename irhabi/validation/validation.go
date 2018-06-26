package validation

import (
	"fmt"
	"reflect"

	"time"

	"git.qasico.com/cuxs/common"
)

type (
	// Validation holding the tag name and taglists available
	Validation struct {
		TagName       string                   // Default tag name for struct validation.
		ValidatorTags map[string]ValidatorFunc // TagRules is a map of RuleFunc indexed by tag name.
	}

	// Request interface validation requests
	Request interface {
		Validate() *Output
		Messages() map[string]string
	}
)

// New creates a new Validation instances.
func New() *Validation {
	return &Validation{
		TagName:       tagName,
		ValidatorTags: ValidatorTags,
	}
}

// Field validates a value based on the provided
// tags and returns result struct.
func (v *Validation) Field(value interface{}, tag string) *Output {
	tags, err := toTag(tag, v.ValidatorTags)
	if err != nil {
		return nil
	}

	r := new(Output)
	r.failureMessages = make(map[string]string)
	for _, t := range tags {
		var e string
		if r.Valid, e = t.Fn(value, t.Param); !r.Valid {
			r.Failure(t.Name, e)
			r.tag = t.Name
			return r.error()
		}
	}

	return r
}

// Struct validates the object of a struct based
// on 'valid' tags and returns errors found indexed
// by the field name.
func (v *Validation) Struct(object interface{}) *Output {
	sv := reflect.ValueOf(object)
	st := reflect.TypeOf(object)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return v.Struct(sv.Elem().Interface())
	}
	if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Interface {
		return &Output{}
	}

	res := &Output{
		Valid:           true,
		failureMessages: make(map[string]string),
	}

	nf := sv.NumField()
	for i := 0; i < nf; i++ {
		f := sv.Field(i)

		ftag := st.Field(i).Tag.Get(v.TagName)

		// take the name from json first,
		// if theres no tag json we get from the field name
		// and it will convert into snake_case string instead of CamelCase.
		fname := st.Field(i).Tag.Get("json")
		if fname == "" {
			fname = common.ToUnderscore(st.Field(i).Name)
		}

		if f.Kind() == reflect.Ptr && !f.IsNil() {
			if oq, oqk := f.Interface().(Request); oqk {
				if e := v.Request(oq); !e.Valid {
					for j, k := range e.failureMessages {
						res.Failure(fname+"."+j, k)
					}

					continue
				}
			}
		}

		for f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}

		if ftag != "" {
			r := v.Field(f.Interface(), ftag)
			if r == nil {
				continue
			}

			if !r.Valid {
				res.Failure(fname+"."+r.tag, fmt.Sprintf(r.Message(), fname))
			}
		}

		if (f.Kind() == reflect.Struct || f.Kind() == reflect.Interface) && f.Type() != reflect.TypeOf(time.Time{}) {
			var e *Output
			if o, ok := f.Interface().(Request); ok {
				e = o.Validate()
			} else {
				e = v.Struct(f.Interface())
			}

			if !e.Valid {
				for j, k := range e.failureMessages {
					res.Failure(fname+"."+j, k)
				}
			}
		}

		if f.Kind() == reflect.Slice && f.Len() > 0 {
			if f.Index(0).Kind() == reflect.Struct || f.Index(0).Kind() == reflect.Ptr {
				for i := 0; i < f.Len(); i++ {
					var e *Output

					if o, ok := f.Index(i).Interface().(Request); ok {
						e = v.Request(o)
					} else {
						e = v.Struct(f.Index(i).Interface())
					}

					if !e.Valid {
						for j, k := range e.failureMessages {
							kk := fmt.Sprintf("%s.%d.%s", fname, i, j)
							res.Failure(kk, k)
						}
					}
				}
			}
		}
	}

	if len(res.failureMessages) > 0 {
		res.Valid = false
		return res.error()
	}

	return res
}

// Request same as Validation.Struct but, this
// should be implement an ValidationRequest interfaces
// so we can do some custom validation and custome error messages.
func (v *Validation) Request(object Request) *Output {
	o := &Output{
		Valid:           true,
		customMessages:  object.Messages(),
		failureMessages: make(map[string]string),
	}

	if os := v.Struct(object); !os.Valid {
		o.Valid = false
		for x, y := range os.failureMessages {
			o.Failure(x, y)
		}
	}

	if or := object.Validate(); or != nil && !or.Valid {
		or.error()

		o.Valid = false
		for x, y := range or.failureMessages {
			o.Failure(x, y)
		}
	}

	if !o.Valid {
		o.error()
	}

	return o
}
