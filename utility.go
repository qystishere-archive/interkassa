package interkassa

import (
	"encoding/json"
	"time"
)

// OptionalString формирует значение опционального параметра - строки.
func OptionalString(value string) OptionalStringValue {
	return &value
}

// OptionalInt32 формирует значение опционального параметра - целого числа.
func OptionalInt32(value int32) OptionalInt32Value {
	return &value
}

// OptionalBool формирует значение опционального параметра - булевого значения (да/нет).
func OptionalBool(value bool) OptionalBoolValue {
	return &value
}

// OptionalTime формирует значение опционального параметра - времени.
func OptionalTime(value time.Time) OptionalTimeValue {
	return &Time{value}
}

func bind(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, to); err != nil {
		return err
	}
	return nil
}