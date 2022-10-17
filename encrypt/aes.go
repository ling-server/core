package encrypt

import (
	"os"
	"sync"

	"github.com/ling-server/core/log"
)

var (
	defaultKeyPath = "/etc/core/key"
)

// AESEncryptor uses AES to encrypt or decrypt string
type AESEncryptor struct {
	keyProvider KeyProvider
	keyParams   map[string]interface{}
}

// NewAESEncryptor returns an instance of an AESEncryptor
func NewAESEncryptor(keyProvider KeyProvider) Encryptor {
	return &AESEncryptor{
		keyProvider: keyProvider,
	}
}

var encryptInstance Encryptor
var encryptOnce sync.Once

// AesInstance ... Get instance of encryptor
func AesInstance() Encryptor {
	encryptOnce.Do(func() {
		kp := os.Getenv("KEY_PATH")
		if len(kp) == 0 {
			kp = defaultKeyPath
		}
		log.Infof("the path of key used by key provider: %s", kp)
		encryptInstance = NewAESEncryptor(NewFileKeyProvider(kp))
	})
	return encryptInstance
}

// Encrypt ...
func (a *AESEncryptor) Encrypt(plaintext string) (string, error) {
	key, err := a.keyProvider.Get(a.keyParams)
	if err != nil {
		return "", err
	}
	return ReversibleEncrypt(plaintext, key)
}

// Decrypt ...
func (a *AESEncryptor) Decrypt(ciphertext string) (string, error) {
	key, err := a.keyProvider.Get(a.keyParams)
	if err != nil {
		return "", err
	}
	return ReversibleDecrypt(ciphertext, key)
}
