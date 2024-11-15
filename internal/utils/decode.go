// internal/utils/decode.go
package utils

import (
	"encoding/base64"
)

func Base64Decode(s string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
