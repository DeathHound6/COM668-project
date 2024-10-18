package utility

import (
	"encoding/json"
	"fmt"
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

func GetFieldsMapFromString(fieldsString string) []KeyValueSchema {
	fields := make([]KeyValueSchema, 0)
	// Each field is separated by `|`
	dbFields := strings.Split(fieldsString, "|")
	for _, field := range dbFields {
		// Each field is mapped `<key>;<value>;<dataType>`
		fieldKV := strings.Split(field, ";")
		if len(fieldKV) != 3 {
			continue
		}
		fields = append(fields, KeyValueSchema{
			Key:   strings.TrimSpace(fieldKV[0]),
			Value: strings.TrimSpace(fieldKV[1]),
			Type:  strings.TrimSpace(fieldKV[2]),
		})
	}
	return fields
}

func GetStringFromFieldsMap(fieldsMap []KeyValueSchema) string {
	parts := make([]string, 0)
	for _, field := range fieldsMap {
		parts = append(parts, fmt.Sprintf("%s;%s;%s", field.Key, field.Value, field.Type))
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
