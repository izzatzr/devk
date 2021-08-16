package genrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/izzatzr/devk/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	Log = logger.NewLogger()
)

// Create generate temporary public & private RSA Key
func Create() (os.File, os.File, error) {

	pvKeyFile, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	pbKeyFile, err := ioutil.TempFile(os.TempDir(), "*.pub")
	if err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	key, pvKeyData, err := privateKey()
	if err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	pbKeyData, err := publicKey(&key.PublicKey)
	if err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	if err = ioutil.WriteFile(pvKeyFile.Name(), pvKeyData, os.ModePerm); err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	if err = ioutil.WriteFile(pbKeyFile.Name(), pbKeyData, os.ModePerm); err != nil {
		return os.File{}, os.File{}, errors.Wrap(err, err.Error())
	}

	return *pvKeyFile, *pbKeyFile, nil
}

func privateKey() (*rsa.PrivateKey, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	err = key.Validate()
	if err != nil {
		return nil, nil, err
	}

	return key, encodePrivateKey(key), nil
}

func publicKey(pvKey *rsa.PublicKey) ([]byte, error) {
	key, err := ssh.NewPublicKey(pvKey)
	if err != nil {
		return nil, err
	}

	keyBytes := ssh.MarshalAuthorizedKey(key)
	return keyBytes, nil
}

func encodePrivateKey(pvKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Bytes: x509.MarshalPKCS1PrivateKey(pvKey)})
}
