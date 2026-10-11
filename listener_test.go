package events

import (
	"context"
	"testing"
)

func TestNewListenerWrapper(t *testing.T) {
	called := false
	var gotCtx context.Context
	var gotEvent *Event
	w := NewListenerWrapper("my-listener", func(ctx context.Context, e *Event) {
		called = true
		gotCtx = ctx
		gotEvent = e
	})
	if w.Name() != "my-listener" {
		t.Fatalf("expect my-listener, got %s", w.Name())
	}

	ctx := context.Background()
	e := NewEvent("lock-1")

	// Listener 接口断言
	var _ Listener = w
	w.On(ctx, e)
	if !called {
		t.Fatal("expect listener func called")
	}
	if gotCtx != ctx {
		t.Fatal("expect ctx passed through")
	}
	if gotEvent != e {
		t.Fatal("expect event passed through")
	}
}
