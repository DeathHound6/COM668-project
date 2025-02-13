package utility

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func SliceHasElement[T comparable](slice []T, element T) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func GenerateRandomUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func ReadJSONStruct[T any](bytes []byte) (*T, error) {
	var data T
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func RemoveElementFromSlice[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func Pointer[T any](value T) *T {
	return &value
}

func MapToSlice(params map[string]any) []string {
	parts := make([]string, 0)
	for key, value := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return parts
}
