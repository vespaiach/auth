package conf

import (
	"crypto/rsa"

	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/pkg/gotils"
)

var (
	defaultPrivateKeyPath = "configs/rsa/id_rsa"
	defaultPublicKeyPath  = "configs/rsa/id_rsa_pub"
)

// RsaKeyConfig holds all rsa's configuration
type RsaKeyConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey
}

func loadRsaConfig() (config *RsaKeyConfig, err error) {
	appbase, e := gotils.GetEnvString("APP_DIR")
	if e != nil {
		appbase = defaultAppDir
	}

	PrivateKeyPath, e := gotils.GetEnvString("PRIVATE_KEY_PATH")
	if e != nil {
		log.Println(e)
		PrivateKeyPath = defaultPrivateKeyPath
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	PublicKeyPath, e := gotils.GetEnvString("PUBLIC_KEY_PATH")
	if e != nil {
		log.Println(e)
		PublicKeyPath = defaultPublicKeyPath
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	if !filepath.IsAbs(PublicKeyPath) {
		PublicKeyPath = path.Join(appbase, PublicKeyPath)
	}

	if !filepath.IsAbs(PrivateKeyPath) {
		PrivateKeyPath = path.Join(appbase, PrivateKeyPath)
	}

	PrivateKey, err := gotils.LoadRsaPrivateKey(PrivateKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	PublicKey, err := gotils.LoadRsaPublicKey(PublicKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	config = &RsaKeyConfig{
		PrivateKeyPath,
		PublicKeyPath,
		PrivateKey,
		PublicKey,
	}

	return
}
