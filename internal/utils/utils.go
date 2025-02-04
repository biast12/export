package utils

import (
	"encoding/base64"
	"log/slog"
	"math/rand"
	"strings"
)

type Map map[string]interface{}

func Ptr[T any](v T) *T {
	return &v
}

func Contains[T comparable](arr []T, val T) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}

	return false
}

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(bytes)
}

func ParseLogLevel(level string, defaultValue slog.Level) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return defaultValue
	}
}
