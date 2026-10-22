package events

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestEventJSONPreservesErrors(t *testing.T) {
	eventErr := errors.New("event failed")
	actionErr := errors.New("action failed")
	e := NewEvent("lock-1").
		SetType(EventTypeLost).
		SetErr(eventErr).
		AddAction(NewAction("renew").SetErr(actionErr))

	jsonString, err := e.ToJsonStringE()
	if err != nil {
		t.Fatalf("ToJsonStringE() error = %v", err)
	}
	if !strings.Contains(jsonString, `"err":"event failed"`) || !strings.Contains(jsonString, `"err":"action failed"`) {
		t.Fatalf("JSON does not preserve error messages: %s", jsonString)
	}

	restored, err := EventFromJsonStringE(jsonString)
	if err != nil {
		t.Fatalf("EventFromJsonStringE() error = %v", err)
	}
	if restored.Err == nil || restored.Err.Error() != eventErr.Error() {
		t.Fatalf("restored event error = %v, want %q", restored.Err, eventErr)
	}
	if len(restored.Actions) != 1 || restored.Actions[0].Err == nil || restored.Actions[0].Err.Error() != actionErr.Error() {
		t.Fatalf("restored action error = %+v, want %q", restored.Actions, actionErr)
	}
	if restored.EventType != EventTypeLost {
		t.Fatalf("restored event type = %v, want %v", restored.EventType, EventTypeLost)
	}
}

func TestEventJSONAcceptsLegacyObjectErrors(t *testing.T) {
	const legacyJSON = `{"id":"event-1","root_id":"event-1","lock_id":"lock-1","err":{},"actions":[{"name":"action","err":{}}]}`

	e, err := EventFromJsonStringE(legacyJSON)
	if err != nil {
		t.Fatalf("EventFromJsonStringE() rejected legacy JSON: %v", err)
	}
	if e.Err != nil || len(e.Actions) != 1 || e.Actions[0].Err != nil {
		t.Fatalf("legacy object errors should decode as nil: %+v", e)
	}
}

func TestEventPublishProvidesErrorForListener(t *testing.T) {
	want := errors.New("listener-visible")
	var got error
	listener := NewListenerWrapper("errors", func(_ context.Context, e *Event) {
		got = e.Err
	})

	NewEvent("lock-1").SetErr(want).Publish(context.Background(), listener)
	if !errors.Is(got, want) {
		t.Fatalf("listener error = %v, want %v", got, want)
	}
}
