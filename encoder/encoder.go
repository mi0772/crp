package encoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func EncryptFile(f string, password string, remove bool) (newFileName string, err error) {
	body, err := os.ReadFile(f)
	newFileName = fmt.Sprintf("%s.%s", f, "encrypted")

	if err != nil {
		return
	}

	encrypted, err := encrypt(body, password)
	if err != nil {
		return
	}

	file, _ := os.OpenFile(newFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	file.Write(encrypted)

	if remove {
		err = os.Remove(f)
		if err != nil {
			return
		}
	}
	return
}

func generateKey(masterPassword string) []byte {
	key := sha256.Sum256([]byte(masterPassword))
	return key[:]
}

func encrypt(body []byte, masterPassword string) ([]byte, error) {
	c, err := aes.NewCipher(generateKey(masterPassword))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, body, nil), nil
}

func DecryptFile(f string, password string, remove bool) (newFileName string, err error) {
	body, err := os.ReadFile(f)
	if err != nil {
		return
	}
	newFileName = strings.Replace(f, ".encrypted", "", 1)
	decryped, err := decrypt(body, password)
	if err != nil {
		return
	}

	err = os.WriteFile(newFileName, decryped, 0644)
	file, _ := os.OpenFile(newFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	file.Write(decryped)

	if remove {
		err = os.Remove(f)
		if err != nil {
			return
		}
	}
	return
}

func decrypt(ciphertext []byte, masterPassword string) ([]byte, error) {
	c, err := aes.NewCipher(generateKey(masterPassword))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
