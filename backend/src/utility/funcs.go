package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

func GetFieldsMapFromString(fieldsString string) ([]KeyValueSchema, error) {
	fields := make([]KeyValueSchema, 0)
	if fieldsString == "" {
		return fields, nil
	}
	// Each field is separated by `|`
	dbFields := strings.Split(fieldsString, "|")
	for _, field := range dbFields {
		// Each field is mapped `<key>;<value>;<dataType>;<required>`
		fieldKV := strings.Split(field, ";")
		if len(fieldKV) != 4 {
			return nil, errors.New("invalid fields string")
		}
		required, err := strconv.ParseBool(strings.TrimSpace(fieldKV[3]))
		if err != nil {
			return nil, err
		}
		fields = append(fields, KeyValueSchema{
			Key:      strings.TrimSpace(fieldKV[0]),
			Value:    strings.TrimSpace(fieldKV[1]),
			Type:     strings.TrimSpace(fieldKV[2]),
			Required: required,
		})
	}
	return fields, nil
}

func GetStringFromFieldsMap(fieldsMap []KeyValueSchema) string {
	parts := make([]string, 0)
	for _, field := range fieldsMap {
		parts = append(parts, fmt.Sprintf("%s;%s;%s;%s", field.Key, field.Value, field.Type, strconv.FormatBool(field.Required)))
	}
	return strings.Join(parts, "|")
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
