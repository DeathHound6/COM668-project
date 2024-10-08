package utility

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func SliceHasElement(slice []reflect.Type, element reflect.Type) bool {
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

func GenerateRandomUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
