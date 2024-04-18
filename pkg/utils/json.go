package utils

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func ToJSONString(obj any) (string, error) {
	s, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(s), nil
}

func ToJSONStringPanic(obj any) string {
	s, err := ToJSONString(obj)
	if err != nil {
		panic("failed to convert to JSON string: " + err.Error())
	}
	return s
}
