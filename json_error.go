package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func marshalError(err error) *string {
	if err == nil {
		return nil
	}
	message := err.Error()
	return &message
}

func unmarshalError(raw json.RawMessage) (error, error) {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "null" {
		return nil, nil
	}

	var message string
	if err := json.Unmarshal(raw, &message); err == nil {
		return errors.New(message), nil
	}

	// Older versions encoded the error interface as an empty object. It carries
	// no usable information, but accepting it keeps old event logs readable.
	var legacy map[string]json.RawMessage
	if err := json.Unmarshal(raw, &legacy); err == nil && legacy != nil {
		return nil, nil
	}

	return nil, fmt.Errorf("decode error message %s: expected string or null", string(raw))
}

func (x *Action) MarshalJSON() ([]byte, error) {
	type actionJSON Action
	return json.Marshal(struct {
		*actionJSON
		Err *string `json:"err"`
	}{actionJSON: (*actionJSON)(x), Err: marshalError(x.Err)})
}

func (x *Action) UnmarshalJSON(data []byte) error {
	type actionJSON Action
	var value struct {
		*actionJSON
		Err json.RawMessage `json:"err"`
	}
	value.actionJSON = (*actionJSON)(x)
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	decoded, err := unmarshalError(value.Err)
	if err != nil {
		return err
	}
	x.Err = decoded
	return nil
}

func (x *Event) MarshalJSON() ([]byte, error) {
	type eventJSON Event
	return json.Marshal(struct {
		*eventJSON
		Err *string `json:"err"`
	}{eventJSON: (*eventJSON)(x), Err: marshalError(x.Err)})
}

func (x *Event) UnmarshalJSON(data []byte) error {
	type eventJSON Event
	var value struct {
		*eventJSON
		Err json.RawMessage `json:"err"`
	}
	value.eventJSON = (*eventJSON)(x)
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	decoded, err := unmarshalError(value.Err)
	if err != nil {
		return err
	}
	x.Err = decoded
	return nil
}
