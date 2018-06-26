package validation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ValidatorFunc is a function that receives the value of a
// field and a parameter used for the respective validation tag.
type ValidatorFunc func(v interface{}, param string) (valid bool, failMessage string)

// ValidatorTags register tag and validator function
// that available on validation.
var ValidatorTags = map[string]ValidatorFunc{
	"required":        validRequired,
	"numeric":         validNumeric,
	"alpha":           validAlpha,
	"alpha_num":       validAlphaNum,
	"alpha_num_space": validAlphaNumSpace,
	"alpha_space":     validAlphaSpace,
	"email":           validEmail,
	"url":             validURL,
	"json":            validJSON,
	"cc":              validCc,
	"lte":             validLte,
	"gte":             validGte,
	"lt":              validLt,
	"gt":              validGt,
	"range":           validRange,
	"contains":        validContains,
	"match":           validMatch,
	"same":            validSame,
	"in":              validIn,
	"not_in":          validNotIn,
	"ean":             validEan,
}

type tag struct {
	Name  string        // name of the tag
	Fn    ValidatorFunc // validation function to call
	Param string        // parameter to send to the validation function
}

func toTag(value string, vfn map[string]ValidatorFunc) ([]tag, error) {
	if value == "-" {
		return []tag{}, errors.New("Tag skipped")
	}

	tl := strings.Split(value, "|")
	tags := make([]tag, 0, len(tl))

	for _, i := range tl {
		t := tag{}
		p := strings.SplitN(i, ":", 2)
		t.Name = strings.Trim(p[0], " ")
		if t.Name == "" {
			return []tag{}, errors.New("Tag cannot be empty")
		}

		if len(p) > 1 {
			t.Param = strings.Trim(p[1], " ")
		}

		// check is tag has been declared
		var found bool
		if t.Fn, found = vfn[t.Name]; !found {
			return []tag{}, fmt.Errorf("Cannot find any tag with name %s", t.Name)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func validRequired(value interface{}, _ string) (v bool, m string) {
	if v = IsNotEmpty(value); !v {
		m = "The %s field is required."
	}
	return
}

func validNumeric(value interface{}, _ string) (v bool, m string) {
	if v = IsNumeric(value); !v {
		m = "The %s must be a number."
	}
	return
}

func validAlpha(value interface{}, _ string) (v bool, m string) {
	if v = IsAlpha(value); !v {
		m = "The %s may only contain letters."
	}
	return
}

func validAlphaNum(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphanumeric(value); !v {
		m = "The %s may only contain letters and numbers."
	}
	return
}

func validAlphaNumSpace(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphanumericSpace(value); !v {
		m = "The %s may only contain letters, numbers and spaces."
	}
	return
}

func validAlphaSpace(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphaSpace(value); !v {
		m = "The %s may only contain letters and spaces."
	}
	return
}

func validEmail(value interface{}, _ string) (v bool, m string) {
	if v = IsEmail(value); !v {
		m = "The %s must be a valid email address."
	}
	return
}

func validURL(value interface{}, _ string) (v bool, m string) {
	if v = IsURL(value); !v {
		m = "The %s format is invalid."
	}
	return
}

func validJSON(value interface{}, _ string) (v bool, m string) {
	if v = IsJSON(value); !v {
		m = "The %s must be a valid JSON string."
	}
	return
}

func validCc(value interface{}, _ string) (v bool, m string) {
	if v = IsCreditCard(value); !v {
		m = "The :attribute must be a valid credit card number."
	}
	return
}

func validLte(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsLowerThanEqual(value, p); !v {
		m = fmt.Sprintf("The %s may not be greater than %v", "%s", p)
	}
	return
}

func validGte(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsGreaterThanEqual(value, p); !v {
		m = fmt.Sprintf("The %s should be greater than %v", "%s", p)
	}
	return
}

func validLt(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsLowerThan(value, p); !v {
		m = fmt.Sprintf("The %s may not be greater than %v", "%s", p)
	}
	return
}

func validGt(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsGreaterThan(value, p); !v {
		m = fmt.Sprintf("The %s should be greater than %v", "%s", p)
	}
	return
}

func validRange(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")

	if len(p) == 2 {
		min := convert(p[0])
		max := convert(p[1])

		if v = IsOnRange(value, min, max); !v {
			m = fmt.Sprintf("The %s must be between %v and %v.", "%s", min, max)
		}
	}
	return
}

//convert string to interface {}
func convert(param string) (p interface{}) {
	var errInt, errFlt error
	p, errInt = strconv.Atoi(param)
	if errInt != nil {
		p, errFlt = strconv.ParseFloat(param, 64)
		if errFlt != nil {

			p = param
		}
	}
	return p
}
func validContains(value interface{}, param string) (v bool, m string) {
	if v = IsContains(value, param); !v {
		m = "The %s format is invalid."
	}
	return
}

func validMatch(value interface{}, param string) (v bool, m string) {
	if v = IsMatches(value, param); !v {
		m = "The %s format is invalid."
	}
	return
}

func validSame(value interface{}, param string) (v bool, m string) {
	if v = IsSame(value, param); !v {
		m = "The %s format is invalid."
	}
	return
}

func validIn(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")
	if v = IsIn(value, p...); !v {
		m = "The selected %s is invalid."
	}
	return
}

func validNotIn(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")
	if v = IsNotIn(value, p...); !v {
		m = "The selected %s is invalid."
	}
	return
}

func validEan(value interface{}, _ string) (v bool, m string) {
	if v = IsValidEan(value); !v {
		m = "The %s field is not valid ean13 code."
	}
	return
}
