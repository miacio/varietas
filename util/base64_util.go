package util

import "encoding/base64"

// Base64Encode
func Base64Encode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

// Base64Decode
func Base64Decode(data string) string {
	b, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(b)
}
