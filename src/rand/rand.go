package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(numBytes int) ([]byte, error) {
	b := make([]byte, numBytes)
	numRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if numRead < numBytes {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}
	return b, nil
}

func String(numBytes int) (string, error) {
	b, err := Bytes(numBytes)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
