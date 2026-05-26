package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"

	"metavida/backend/core"
)

// 18 raw bytes encodes to exactly 24 base64 chars (no padding).
const passwordHashBytes = 18

func HashPassword(password string) (string, error) {
	return computePasswordHMAC(password), nil
}

func VerifyPassword(password string, storedHash string) bool {
	return subtle.ConstantTimeCompare(
		[]byte(computePasswordHMAC(password)),
		[]byte(storedHash),
	) == 1
}

// HMAC-SHA256 keyed by the app secret acts as a pepper, preventing offline
// attacks without the key. The first 18 bytes give 24 base64 chars.
func computePasswordHMAC(password string) string {
	mac := hmac.New(sha256.New, []byte(core.GetAuthSecret()))
	mac.Write([]byte(password))
	return base64.RawStdEncoding.EncodeToString(mac.Sum(nil)[:passwordHashBytes])
}
