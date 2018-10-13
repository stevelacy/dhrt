package dhrt

import (
	"crypto/sha256"
	"encoding/base64"
)

func Hash(data string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func HashToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func StringToHash(data string) []byte {
	hash, _ := base64.StdEncoding.DecodeString(data)
	return hash
}
