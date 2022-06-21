package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
	mathrand "math/rand"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
	iter       = 1000
	keyLen     = 16
)

// 生成8-32之间的随机数字
func GenRandSalt() ([]byte, error) {
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenCryptKey(pwd []byte, salt []byte) string {
	cryptKey := pbkdf2.Key(pwd, salt, iter, keyLen, sha256.New)
	return base64.StdEncoding.EncodeToString(cryptKey)
}
