package cf

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func getEnvString(name string) (val string, err error) {
	val = os.Getenv(name)

	if len(val) == 0 {
		err = errors.New("not found environment variable named: " + name)
	}

	return
}

// GetEnvInt get environment variable and convert it to int type.
func getEnvInt(name string) (val int, err error) {
	valenv := os.Getenv(name)

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val, err = strconv.Atoi(valenv)
	}

	return
}

// getEnvBool get environment variable and convert it to bool type.
func getEnvBool(name string) (val bool, err error) {
	valenv := strings.ToLower(os.Getenv(name))

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val = valenv == "true" || valenv == "1"
	}

	return
}

// getEnvStringSlice get environment variable and convert it to []string type.
// Each value is separated by "|"
func getEnvStringSlice(name string) (val []string, err error) {
	valenv := strings.ToLower(os.Getenv(name))

	if len(valenv) == 0 {
		err = errors.New("not found environment variable named: " + name)
	} else {
		val = strings.Split(valenv, "|")
	}

	return
}

// getEnvIntSlice get environment variable and convert it to []int type.
// Each value is separated by "|"
func getEnvIntSlice(name string) (val []int, err error) {
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
