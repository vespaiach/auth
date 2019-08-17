package gotils

import (
	"crypto/rsa"
	"io/ioutil"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// LoadRsaPublicKey loads RSA public key file
func LoadRsaPublicKey(filePath string) (key *rsa.PublicKey, err error) {
	var dt []byte
	dt, err = ioutil.ReadFile(filePath)

	if err == nil {
		key, err = jwtgo.ParseRSAPublicKeyFromPEM(dt)
	}

	return
}

// LoadRsaPrivateKey loads RSA private key file
func LoadRsaPrivateKey(filePath string) (key *rsa.PrivateKey, err error) {
	var dt []byte
	dt, err = ioutil.ReadFile(filePath)

	if err == nil {
		key, err = jwtgo.ParseRSAPrivateKeyFromPEM(dt)
	}

	return
}
