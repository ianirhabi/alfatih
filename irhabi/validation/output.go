package validation

import (
	"path/filepath"
	"regexp"
	"strings"

	"git.qasico.com/cuxs/common"
)

// Output format response when running validations
type Output struct {
	Valid           bool              // state of validation
	tag             string            // failing tags
	messages        map[string]string // compiled error messages
	failureMessages map[string]string // failing error messages
	customMessages  map[string]string // custom messages
	failureKeys     []string
}

// Messages is a map which contains all errors from validating a struct.
func (o *Output) Messages() map[string]string {
	return o.messages
}

// Message returns failure message by key provided as parameter,
// if key is not provided will return the first failure message.
func (o *Output) Message(k ...string) string {
	if len(k) > 0 {
		return o.messages[k[0]]
	}

	var vl string
	for _, vl = range o.messages {
		break
	}

	return vl
}

// Failure set an failure message as key and value
func (o *Output) Failure(k string, e string) {
	if o.failureMessages == nil {
		o.failureMessages = make(map[string]string)
	}

	o.Valid = false
	o.failureKeys = append(o.failureKeys, k)
	o.failureMessages[k] = e
}

// error wrapping up if validation failing.
func (o *Output) error() *Output {
	o.applyCustomMessage()

	res := make(map[string]string)
	for _, i := range o.failureKeys {
		k := strings.TrimSuffix(i, filepath.Ext(i))
		if _, ok := res[k]; !ok {
			res[k] = o.failureMessages[i]
		}
	}

	o.messages = res
	return o
}

// applyCustomMessage compile any message with customMessage provided.
func (o *Output) applyCustomMessage() {
	for i := range o.failureMessages {
		if c := o.customMessages[i]; c != "" {
			o.Failure(i, c)
			continue
		}

		if IsMatches(i, "(\\.[0-9]+\\.[a-z]+\\.[a-z]*)$") {
			re := regexp.MustCompile("[^a-z.]")
			ix := re.ReplaceAllString(i, "*")
			if c := o.customMessages[ix]; c != "" {
				o.Failure(i, c)
			}
		}
	}
}

// Error implement error type interfaces
func (o *Output) Error() string {
	return common.ToJSON(o.Messages())
}

// SetError is helper to manualy set error where ever its needed
func SetError(field string, value string) *Output {
	o := new(Output)
	o.Failure(field, value)

	return o.error()
}
