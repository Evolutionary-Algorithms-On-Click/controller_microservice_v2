package models

import (
	"encoding/json"
	"github.com/google/uuid"
)

// StringUUID is a custom type that allows for marshalling and unmarshalling of UUIDs as strings.
type StringUUID uuid.UUID

// MarshalJSON implements the json.Marshaler interface.
func (su StringUUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(su).String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (su *StringUUID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*su = StringUUID(id)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (su StringUUID) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(su).String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (su *StringUUID) UnmarshalText(text []byte) error {
	id, err := uuid.Parse(string(text))
	if err != nil {
		return err
	}
	*su = StringUUID(id)
	return nil
}

// ToUUID converts a StringUUID to a uuid.UUID.
func (su StringUUID) ToUUID() uuid.UUID {
	return uuid.UUID(su)
}
