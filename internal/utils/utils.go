package utils

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
