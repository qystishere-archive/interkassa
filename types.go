package interkassa

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	Time struct {
		time.Time
	}

	OptionalStringValue *string
	OptionalInt32Value  *int32
	OptionalBoolValue   *bool
	OptionalTimeValue   *Time

	Fields map[string]string
	Form   struct {
		Method string `json:"method"`
		Action string `json:"action"`
		Fields Fields `json:"fields"`
	}
)

func (t *Time) UnmarshalJSON(bytes []byte) error {
	realTime, err := time.Parse(
		"2006-01-02 15:04:05",
		strings.Trim(string(bytes), "\""),
	)
	if err != nil {
		return err
	}
	*t = Time{realTime}
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

func (t *Time) Format(s string) string {
	return t.Time.Format(s)
}