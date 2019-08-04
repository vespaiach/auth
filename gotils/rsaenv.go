package gotils

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"os"
	"path"

	jwtgo "github.com/dgrijalva/jwt-go"
)

func hasPath(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func getAbsPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	c, _ := os.Getwd()
	return path.Join(c, p)
}

// LoadRsaPublicKey loads RSA public key file
func LoadRsaPublicKey(filePath string) (key *rsa.PublicKey, err error) {
	abspath := getAbsPath(filePath)

	if !hasPath(abspath) {
		err = errors.New(abspath + " doesn't exist")
		return
	}

	var dt []byte
	dt, err = ioutil.ReadFile(abspath)

	if err == nil {
		key, err = jwtgo.ParseRSAPublicKeyFromPEM(dt)
	}

	return
}

// LoadRsaPrivateKey loads RSA private key file
func LoadRsaPrivateKey(filePath string) (key *rsa.PrivateKey, err error) {
	abspath := getAbsPath(filePath)

	if !hasPath(abspath) {
		err = errors.New(abspath + " doesn't exist")
		return
	}

	var dt []byte
	dt, err = ioutil.ReadFile(abspath)

	if err == nil {
		key, err = jwtgo.ParseRSAPrivateKeyFromPEM(dt)
	}

	return
}
