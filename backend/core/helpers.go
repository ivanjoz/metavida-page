package core

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

func SUnixTime() int32 {
	return int32((time.Now().Unix() - 1e9) / 2)
}

func BytesToBase64(source []byte, useUrlEncoded ...bool) string {
	contentString := base64.StdEncoding.EncodeToString(source)
	if len(useUrlEncoded) == 1 && useUrlEncoded[0] {
		contentString = MakeB64UrlEncode(contentString)
	}
	return contentString
}

func Base64ToBytes(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}

func MakeB64UrlEncode(contentString string) string {
	contentString = strings.ReplaceAll(contentString, "+", "-")
	contentString = strings.ReplaceAll(contentString, "/", "_")
	contentString = strings.ReplaceAll(contentString, "=", "~")
	return contentString
}

func MakeB64UrlDecode(contentString string) string {
	contentString = strings.ReplaceAll(contentString, "_", "/")
	contentString = strings.ReplaceAll(contentString, "-", "+")
	contentString = strings.ReplaceAll(contentString, "~", "=")
	return contentString
}

func RandomTokenID() string {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes)
}
