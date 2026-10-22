package events

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNewAction(t *testing.T) {
	a := NewAction("foo")
	if a.Name != "foo" {
		t.Fatalf("expect name foo, got %s", a.Name)
	}
	if a.StartTime == nil {
		t.Fatal("expect StartTime not nil")
	}
}

func TestAction_End(t *testing.T) {
	a := NewAction("foo")
	if a.EndTime != nil {
		t.Fatal("expect EndTime nil before End")
	}
	ret := a.End()
	if ret != a {
		t.Fatal("expect End return self")
	}
	if a.EndTime == nil {
		t.Fatal("expect EndTime not nil after End")
	}
}

func TestAction_Cost(t *testing.T) {
	// StartTime 和 EndTime 都缺失时返回 0
	empty := &Action{}
	if got := empty.Cost(); got != 0 {
		t.Fatalf("expect 0, got %v", got)
	}
	// 只有 StartTime 没有 EndTime 时返回 0
	started := NewAction("foo")
	if got := started.Cost(); got != 0 {
		t.Fatalf("expect 0, got %v", got)
	}
	// 正常路径：Cost >= 0
	started.End()
	if got := started.Cost(); got < 0 {
		t.Fatalf("expect Cost >= 0, got %v", got)
	}
	// 只有 EndTime 没有 StartTime 时返回 0
	endOnly := &Action{}
	endOnly.End()
	endOnly.StartTime = nil
	if got := endOnly.Cost(); got != 0 {
		t.Fatalf("expect 0, got %v", got)
	}
	// 已知时间差验证计算正确
	start := time.Now()
	end := start.Add(100 * time.Millisecond)
	fixed := &Action{StartTime: &start, EndTime: &end}
	if got := fixed.Cost(); got != 100*time.Millisecond {
		t.Fatalf("expect 100ms, got %v", got)
	}
}

func TestAction_SetName(t *testing.T) {
	a := NewAction("foo")
	if ret := a.SetName("bar"); ret != a || a.Name != "bar" {
		t.Fatal("SetName failed")
	}
}

func TestAction_Err(t *testing.T) {
	sentinel := errors.New("sentinel")
	a := NewAction("foo")

	// 无错误时
	if msg := a.GetErrMsg(); msg != "" {
		t.Fatalf("expect empty err msg, got %s", msg)
	}
	if a.ErrorIs(sentinel) {
		t.Fatal("expect ErrorIs false when Err is nil")
	}

	// 设置错误
	if ret := a.SetErr(sentinel); ret != a {
		t.Fatal("expect SetErr return self")
	}
	if msg := a.GetErrMsg(); msg != "sentinel" {
		t.Fatalf("expect sentinel, got %s", msg)
	}
	if !a.ErrorIs(sentinel) {
		t.Fatal("expect ErrorIs true for same error")
	}
	if a.ErrorIs(errors.New("other")) {
		t.Fatal("expect ErrorIs false for other error")
	}

	// 包装错误 errors.Is 语义
	wrapped := NewAction("wrapped").SetErr(&wrapError{msg: "wrap", err: sentinel})
	if !wrapped.ErrorIs(sentinel) {
		t.Fatal("expect ErrorIs true for wrapped sentinel")
	}

	// 清空错误
	if ret := a.ClearErr(); ret != a {
		t.Fatal("expect ClearErr return self")
	}
	if a.Err != nil {
		t.Fatal("expect Err nil after ClearErr")
	}
	if msg := a.GetErrMsg(); msg != "" {
		t.Fatalf("expect empty err msg after ClearErr, got %s", msg)
	}
}

type wrapError struct {
	msg string
	err error
}

func (x *wrapError) Error() string { return x.msg }
func (x *wrapError) Unwrap() error { return x.err }

func TestAction_PayloadMap(t *testing.T) {
	a := NewAction("foo")

	// SetPayloadMap
	m := map[string]any{"k": "v"}
	if ret := a.SetPayloadMap(m); ret != a {
		t.Fatal("expect SetPayloadMap return self")
	}
	if got := a.GetPayloadMap(); len(got) != 1 || got["k"] != "v" {
		t.Fatalf("unexpected payload map: %v", got)
	}

	// AddPayload 在已有 map 上追加
	a.AddPayload("k2", 42)
	if got := a.GetPayloadMap(); got["k2"] != 42 {
		t.Fatalf("unexpected payload map after AddPayload: %v", got)
	}

	// ClearPayloadMap
	if ret := a.ClearPayloadMap(); ret != a {
		t.Fatal("expect ClearPayloadMap return self")
	}
	if a.GetPayloadMap() != nil {
		t.Fatal("expect nil payload map after ClearPayloadMap")
	}

	// AddPayload 在 nil map 上自动建 map
	b := NewAction("bar")
	if ret := b.AddPayload("k", "v"); ret != b {
		t.Fatal("expect AddPayload return self")
	}
	if got := b.GetPayloadMap(); got["k"] != "v" {
		t.Fatalf("unexpected payload map: %v", got)
	}
}

func TestAction_GetPayload(t *testing.T) {
	// nil map 路径
	a := NewAction("foo")
	if v, ok := a.GetPayload("missing"); ok || v != nil {
		t.Fatalf("expect (nil,false), got (%v,%v)", v, ok)
	}

	a.AddPayload("s", "hello").AddPayload("i", 7)
	if v, ok := a.GetPayload("s"); !ok || v != "hello" {
		t.Fatalf("expect (hello,true), got (%v,%v)", v, ok)
	}
	if v, ok := a.GetPayload("missing"); ok || v != nil {
		t.Fatalf("expect (nil,false), got (%v,%v)", v, ok)
	}
}

func TestAction_GetPayloadAsString(t *testing.T) {
	a := NewAction("foo")
	tests := []struct {
		key  string
		want string
	}{
		{"missing", ""},
		{"string", "hello"},
		{"int", "7"},
		{"bool", "true"},
		{"float", "1.5"},
		{"unsupported", ""},
		{"nil", ""},
	}
	a.AddPayload("string", "hello").
		AddPayload("int", 7).
		AddPayload("bool", true).
		AddPayload("float", 1.5).
		AddPayload("unsupported", []int{1}).
		AddPayload("nil", nil)

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := a.GetPayloadAsString(tt.key); got != tt.want {
				t.Fatalf("GetPayloadAsString(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestAction_GetPayloadAsInt(t *testing.T) {
	a := NewAction("foo")
	tests := []struct {
		key  string
		want int
	}{
		{"missing", 0},
		{"int", 7},
		{"int64", 8},
		{"uint", 9},
		{"number", 10},
		{"number_float_int", 12},
		{"string", 11},
		{"invalid", 0},
		{"overflow", 0},
		{"float", 0},
		{"unsupported", 0},
		{"nil", 0},
	}
	a.AddPayload("int", 7).
		AddPayload("int64", int64(8)).
		AddPayload("uint", uint(9)).
		AddPayload("number", json.Number("10")).
		AddPayload("number_float_int", json.Number("12.0")).
		AddPayload("string", "11").
		AddPayload("invalid", "not-an-int").
		AddPayload("overflow", json.Number("9223372036854775808")).
		AddPayload("float", 1.5).
		AddPayload("unsupported", []int{1}).
		AddPayload("nil", nil)

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := a.GetPayloadAsInt(tt.key); got != tt.want {
				t.Fatalf("GetPayloadAsInt(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}
