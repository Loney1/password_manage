package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"time"
)

// Time-based One-time Password (TOTP) algorithm specified in RFC 6238
// thanks: github.com/gotoolkits/AuthOTP
func computeCode(secret string, value int64) int {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return -1
	}
	hash := hmac.New(sha1.New, key)
	err = binary.Write(hash, binary.BigEndian, value)
	if err != nil {
		return -1
	}
	h := hash.Sum(nil)
	offset := h[19] & 0x0f
	truncated := binary.BigEndian.Uint32(h[offset : offset+4])
	truncated &= 0x7fffffff
	code := truncated % 1000000
	return int(code)
}

func TotpCheck(secret string, code int) bool {
	const windowSize int = 5 // valid range: technically 0..100 or so, but beyond 3-5 is probably bad security
	t0 := int(time.Now().UTC().Unix() / 30)
	minT := t0 - (windowSize / 2)
	maxT := t0 + (windowSize / 2)
	for t := minT; t <= maxT; t++ {
		if computeCode(secret, int64(t)) == code {
			return true
		}
	}
	return false
}
