package comtype

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Defined error codes
const (
	ErrInvalidCredential  = 1
	ErrInvalidData        = 2
	ErrDuplicatedData     = 3
	ErrDataValidationFail = 4
	ErrMissingCredential  = 5
	ErrNoPermision        = 6

	ErrQueryDataFail    = 100
	ErrHandleDataFail   = 101
	ErrDataNotFound     = 102
	ErrBadRequest       = 103
	ErrRequestCancelled = 104
)

var errorMessageMapping map[int]string

var defaultErrorMessageMappings = map[int]string{
	ErrInvalidCredential:  "invalid credential",
	ErrInvalidData:        "invalid data",
	ErrDuplicatedData:     "duplicated data",
	ErrDataValidationFail: "data validation fail",

	ErrQueryDataFail:  "queried data fail",
	ErrHandleDataFail: "handled data fail",
	ErrDataNotFound:   "data not found",
	ErrBadRequest:     "bad request",
}

// CommonError is system common error
type CommonError struct {
	Err      error
	location string
	Code     int
	Content  map[string]string
}

// NewCommonError create common error instance
func NewCommonError(err error, location string, code int, content map[string]string) *CommonError {
	return &CommonError{
		err,
		location,
		code,
		content,
	}
}

func (err *CommonError) Error() string {
	return getErrorMessage(err.Code)
}

// Is will check if matched error code
func (err *CommonError) Is(code int) bool {
	return err.Code == code
}

// Debug print real error messag
func (err *CommonError) Debug() string {
	if err != nil && err.Err != nil {
		return err.Err.Error()
	}
	return ""
}

// Info return error info for logging
func (err *CommonError) Info() (string, error) {
	return err.location, err.Err
}

func init() {
	e := loadErrorMessages("./error.yml")
	if e != nil {
		log.Println("err: counldn't find error.yml", e)
		errorMessageMapping = make(map[int]string)
	}
}

func getErrorMessage(code int) string {
	str, ok := errorMessageMapping[code]
	if !ok {
		defaultStr, _ := defaultErrorMessageMappings[code]
		return defaultStr
	}
	return str
}

func loadErrorMessages(file string) error {
	appbase, _ := os.Getwd()

	ymlfile, err := ioutil.ReadFile(path.Join(appbase, file))
	if err != nil {
		return err
	}

	m := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(ymlfile, &m); err != nil {
		return err
	}

	errorMessageMapping := make(map[int]string)

	for key, val := range m {
		k, ok := key.(int)

		if !ok {
			continue
		}

		v, ok := val.(string)

		if !ok {
			continue
		}

		errorMessageMapping[k] = v
	}

	return nil
}
