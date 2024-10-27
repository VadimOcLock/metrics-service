package hashutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ComputeHMAC вычисляет HMAC SHA256 для переданного сообщения и ключа
func ComputeHMAC(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}
