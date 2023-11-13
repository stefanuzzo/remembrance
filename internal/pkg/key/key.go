// TODO
package key

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/stefanuzzo/utilities"
)

const defaultKeysPath string = "./keys"

func InitializeKeys(keysPath string) error {
	var actualKeysPath string
	if keysPath == "" {
		actualKeysPath = defaultKeysPath

	} else {
		actualKeysPath = keysPath
	}

	err := utilities.EnsureDirectoryExists(actualKeysPath)
	if err != nil {
		return err
	}

	// TODO verify slash count
	privateKeyFilepath := actualKeysPath + "/privateKey"
	_, err = os.Stat(privateKeyFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			key, keyErr := generatePrivateKey("rsa", 2048)
			if keyErr != nil {
				return keyErr
			}

			err = persistPrivateKey(key, privateKeyFilepath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err = readPrivateKey(privateKeyFilepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func generatePrivateKey(keyAlgorithm string, bits int) (*crypto.PrivateKey, error) {
	actualKeyAlgorithm := strings.ToUpper(keyAlgorithm)
	if actualKeyAlgorithm == "RSA" {
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, err
		}

		var cryptoKey crypto.PrivateKey = *key
		return &cryptoKey, nil

	} else if actualKeyAlgorithm == "ED25519" {
		_, key, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}

		var cryptoKey crypto.PrivateKey = key
		return &cryptoKey, nil

	} else if actualKeyAlgorithm == "ECDSA" {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		var cryptoKey crypto.PrivateKey = *key
		return &cryptoKey, nil

	} else {
		return nil, fmt.Errorf("unsupported key algorithm: %s", keyAlgorithm)
	}
}

func readPrivateKey(filepath string) (*crypto.PrivateKey, error) {

	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(fileContents)
	if err != nil {
		return nil, err
	}

	var cryptoKey crypto.PrivateKey = *key
	return &cryptoKey, nil
}

func persistPrivateKey(key *crypto.PrivateKey, filepath string) error {
	var keyBytes []byte
	switch k := (*key).(type) {
	case rsa.PrivateKey:
		keyBytes = x509.MarshalPKCS1PrivateKey(&k)

	default:
		keyBytes = nil
	}
	return os.WriteFile(filepath, keyBytes, 0)
}
