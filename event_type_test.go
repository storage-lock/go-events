package events

import "testing"

func TestEventTypeValuesAndStrings(t *testing.T) {
	tests := []struct {
		name  string
		typ   EventType
		value EventType
		text  string
	}{
		{"unknown", EventTypeUnknown, 0, "unknown"},
		{"create_lock", EventTypeCreateLock, 1, "create"},
		{"lock", EventTypeLock, 2, "lock"},
		{"unlock", EventTypeUnlock, 3, "unlock"},
		{"watchdog", EventTypeWatchDog, 4, "watchdog"},
		{"expired", EventTypeExpired, 5, "expired"},
		{"busy", EventTypeBusy, 6, "busy"},
		{"renewed", EventTypeRenewed, 7, "renewed"},
		{"lost", EventTypeLost, 8, "lost"},
		{"wait_changed", EventTypeWaitChanged, 9, "wait_changed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.typ != tt.value {
				t.Errorf("value = %d, want %d", tt.typ, tt.value)
			}
			if got := tt.typ.String(); got != tt.text {
				t.Errorf("String() = %q, want %q", got, tt.text)
			}
		})
	}
}

func TestEventTypeStringUnknownValue(t *testing.T) {
	if got := EventType(999).String(); got != "unknown" {
		t.Errorf("String() = %q, want %q", got, "unknown")
	}
}
