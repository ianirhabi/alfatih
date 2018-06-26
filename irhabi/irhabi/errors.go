package irhabi

type (
	// DataNotExistsError error for inexisting data
	DataNotExistsError struct {
		Message string
		Errors  map[string]string
	}

	// DataDuplicateError error that data is duplicate.
	DataDuplicateError struct {
		Message string
		Errors  map[string]string
	}
)

// Error makes it compatible with `error` interface.
func (e *DataNotExistsError) Error() string {
	return e.Message
}

// Error makes it compatible with `error` interface.
func (e *DataDuplicateError) Error() string {
	return e.Message
}

// ErrDataNotExists error interface for inexists data
// this will cause error 422 with errors data
// fail should has 2 string, first is the key follow with values
// e := ErrDataNotExists("email", "This email is not exists.")
func ErrDataNotExists(fail ...string) error {
	e := &DataNotExistsError{}
	e.Errors = map[string]string{fail[0]: fail[1]}
	e.Message = fail[1]

	return e
}

// ErrDataExists error interface for already exists data
// this will cause error 422 with errors data
// fail should has 2 string, first is the key follow with values
// e := ErrDataNotExists("email", "This email is not exists.")
func ErrDataExists(fail ...string) error {
	e := &DataDuplicateError{}
	e.Errors = map[string]string{fail[0]: fail[1]}
	e.Message = fail[1]

	return e
}
