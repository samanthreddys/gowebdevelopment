package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	//RememberTokenBytes constant
	RememberTokenBytes = 512
)

//Bytes is used to generate n random bytes. uses crypto/rand package
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//String used to generate a byte slice of size n
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//RememberToken used to generate a byte slice
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
