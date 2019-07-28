package gotils

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// GetEnvInt get environment variable and convert it to int type.
func GetEnvInt(name string) (val int, err error) {
	valenv := os.Getenv(name)

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val, err = strconv.Atoi(valenv)
	}

	return
}

// GetEnvBool get environment variable and convert it to bool type.
func GetEnvBool(name string) (val bool, err error) {
	valenv := strings.ToLower(os.Getenv(name))

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val = valenv == "true" || valenv == "1"
	}

	return
}

// GetEnvStringSlice get environment variable and convert it to []string type.
// Each value is separated by "|"
func GetEnvStringSlice(name string) (val []string, err error) {
	valenv := strings.ToLower(os.Getenv(name))

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val = strings.Split(valenv, "|")
	}

	return
}

// GetEnvIntSlice get environment variable and convert it to []int type.
// Each value is separated by "|"
func GetEnvIntSlice(name string) (val []int, err error) {
	valenv := strings.ToLower(os.Getenv(name))

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		strVals := strings.Split(valenv, "|")
		val = make([]int, len(strVals))

		for _, s := range strVals {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

			val = append(val, i)
		}
	}

	return
}
