package events

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/storage-lock/go-storage"
)

func TestNewEvent(t *testing.T) {
	e := NewEvent("lock-1")
	if e.LockId != "lock-1" {
		t.Fatalf("expect lock-1, got %s", e.LockId)
	}
	if e.ID == "" {
		t.Fatal("expect ID not empty")
	}
	if !strings.HasPrefix(e.ID, EventIDPrefix) {
		t.Fatalf("expect ID prefix %s, got %s", EventIDPrefix, e.ID)
	}
	if e.RootID != e.ID {
		t.Fatalf("expect RootID == ID, got %s vs %s", e.RootID, e.ID)
	}
	if e.StartTime == nil {
		t.Fatal("expect StartTime not nil")
	}
	if !e.IsRootEvent() {
		t.Fatal("expect new event is root event")
	}
}

func TestEvent_End(t *testing.T) {
	e := NewEvent("lock-1")
	if ret := e.End(); ret != e {
		t.Fatal("expect End return self")
	}
	if e.EndTime == nil {
		t.Fatal("expect EndTime not nil after End")
	}
}

func TestEvent_Fork(t *testing.T) {
	parent := NewEvent("lock-1").
		SetStorageName("storage").
		SetOwnerId("owner").
		SetWatchDogId("watchdog").
		SetType(EventTypeLock)
	parent.AddListeners(NewListenerWrapper("l", func(ctx context.Context, e *Event) {}))
	parent.LockInformation = &storage.LockInformation{LockId: "lock-1"}

	child := parent.Fork()
	if child.ID == "" || child.ID == parent.ID {
		t.Fatalf("expect child has new ID, got %s", child.ID)
	}
	if child.RootID != parent.RootID {
		t.Fatalf("expect child RootID %s, got %s", parent.RootID, child.RootID)
	}
	if child.ParentID != parent.ID {
		t.Fatalf("expect child ParentID %s, got %s", parent.ID, child.ParentID)
	}
	if child.Parent != parent {
		t.Fatal("expect child Parent point to parent")
	}
	if child.LockId != parent.LockId {
		t.Fatal("expect child inherit LockId")
	}
	if child.StorageName != parent.StorageName {
		t.Fatal("expect child inherit StorageName")
	}
	if child.OwnerId != parent.OwnerId {
		t.Fatal("expect child inherit OwnerId")
	}
	if child.WatchDogId != parent.WatchDogId {
		t.Fatal("expect child inherit WatchDogId")
	}
	if child.EventType != parent.EventType {
		t.Fatal("expect child inherit EventType")
	}
	if child.LockInformation != parent.LockInformation {
		t.Fatal("expect child inherit LockInformation")
	}
	if len(child.Listeners) != len(parent.Listeners) {
		t.Fatal("expect child inherit Listeners")
	}
	if child.StartTime == nil {
		t.Fatal("expect child StartTime not nil")
	}
}

func TestEvent_Setters(t *testing.T) {
	e := NewEvent("lock-1")
	if ret := e.SetRootID("root"); ret != e || e.RootID != "root" {
		t.Fatal("SetRootID failed")
	}
	if ret := e.SetStorageName("s"); ret != e || e.StorageName != "s" {
		t.Fatal("SetStorageName failed")
	}
	if ret := e.SetLockId("l"); ret != e || e.LockId != "l" {
		t.Fatal("SetLockId failed")
	}
	if ret := e.SetOwnerId("o"); ret != e || e.OwnerId != "o" {
		t.Fatal("SetOwnerId failed")
	}
	if ret := e.SetType(EventTypeUnlock); ret != e || e.EventType != EventTypeUnlock {
		t.Fatal("SetType failed")
	}
	sentinel := errors.New("boom")
	if ret := e.SetErr(sentinel); ret != e || e.Err != sentinel {
		t.Fatal("SetErr failed")
	}
	if ret := e.SetWatchDogId("w"); ret != e || e.WatchDogId != "w" {
		t.Fatal("SetWatchDogId failed")
	}
	info := &storage.LockInformation{LockId: "l"}
	if ret := e.SetLockInformation(info); ret != e || e.LockInformation != info {
		t.Fatal("SetLockInformation failed")
	}
}

func TestEvent_AddAction(t *testing.T) {
	e := NewEvent("lock-1")
	a := NewAction("a1")
	if ret := e.AddAction(a); ret != e {
		t.Fatal("expect AddAction return self")
	}
	if len(e.Actions) != 1 || e.Actions[0] != a {
		t.Fatal("AddAction failed")
	}
	if ret := e.AddActionByName("a2"); ret != e {
		t.Fatal("expect AddActionByName return self")
	}
	if len(e.Actions) != 2 || e.Actions[1].Name != "a2" {
		t.Fatal("AddActionByName failed")
	}
}

func TestEvent_Listeners(t *testing.T) {
	e := NewEvent("lock-1")
	l1 := NewListenerWrapper("l1", func(ctx context.Context, e *Event) {})
	l2 := NewListenerWrapper("l2", func(ctx context.Context, e *Event) {})

	if ret := e.SetListeners([]Listener{l1}); ret != e {
		t.Fatal("expect SetListeners return self")
	}
	if len(e.Listeners) != 1 {
		t.Fatal("SetListeners failed")
	}
	if ret := e.AddListeners(l2); ret != e {
		t.Fatal("expect AddListeners return self")
	}
	if len(e.Listeners) != 2 {
		t.Fatal("AddListeners failed")
	}
	if ret := e.ClearListeners(); ret != e {
		t.Fatal("expect ClearListeners return self")
	}
	if e.Listeners != nil {
		t.Fatal("expect Listeners nil after ClearListeners")
	}
}

func TestEvent_Publish(t *testing.T) {
	ctx := context.Background()

	// 无 listener：仅验证自动设置 EndTime 不 panic
	e := NewEvent("lock-1")
	if e.EndTime != nil {
		t.Fatal("expect EndTime nil before Publish")
	}
	e.Publish(ctx)
	if e.EndTime == nil {
		t.Fatal("expect EndTime set by Publish")
	}

	// 已设置 EndTime 时 Publish 不应覆盖
	fixed := NewEvent("lock-1").End()
	endTime := *fixed.EndTime
	fixed.Publish(ctx)
	if !fixed.EndTime.Equal(endTime) {
		t.Fatal("expect Publish keep existing EndTime")
	}

	// 通过 Listeners 字段订阅的 listener 会被触发
	var fieldCalls int
	fl := NewListenerWrapper("field", func(ctx context.Context, e *Event) { fieldCalls++ })
	e2 := NewEvent("lock-1").SetListeners([]Listener{fl})
	e2.Publish(ctx)
	if fieldCalls != 1 {
		t.Fatalf("expect field listener called once, got %d", fieldCalls)
	}

	// 通过可变参数传入的 listener 会被触发
	var argCalls int
	al := NewListenerWrapper("arg", func(ctx context.Context, e *Event) { argCalls++ })
	e3 := NewEvent("lock-1")
	e3.Publish(ctx, al)
	if argCalls != 1 {
		t.Fatalf("expect arg listener called once, got %d", argCalls)
	}

	// 两者同时存在时都会被触发
	var c1, c2 int
	l1 := NewListenerWrapper("l1", func(ctx context.Context, e *Event) { c1++ })
	l2 := NewListenerWrapper("l2", func(ctx context.Context, e *Event) { c2++ })
	e4 := NewEvent("lock-1").SetListeners([]Listener{l1})
	e4.Publish(ctx, l2)
	if c1 != 1 || c2 != 1 {
		t.Fatalf("expect both listeners called, got %d %d", c1, c2)
	}
}

func TestEvent_IsRootEvent(t *testing.T) {
	if !NewEvent("l").IsRootEvent() {
		t.Fatal("expect new event is root")
	}
	// RootID 为空
	e := NewEvent("l")
	e.RootID = ""
	if e.IsRootEvent() {
		t.Fatal("expect false when RootID empty")
	}
	// ID 为空
	e = NewEvent("l")
	e.ID = ""
	if e.IsRootEvent() {
		t.Fatal("expect false when ID empty")
	}
	// RootID 与 ID 不一致（fork 出来的子事件）
	parent := NewEvent("l")
	child := parent.Fork()
	if child.IsRootEvent() {
		t.Fatal("expect forked child not root")
	}
	// 两个都为空
	empty := &Event{}
	if empty.IsRootEvent() {
		t.Fatal("expect false when both empty")
	}
}

func TestEvent_Parent(t *testing.T) {
	parent := NewEvent("parent-lock")
	child := NewEvent("child-lock")
	if ret := child.SetParent(parent); ret != child {
		t.Fatal("expect SetParent return self")
	}
	if child.Parent != parent {
		t.Fatal("expect Parent set")
	}
	if child.ParentID != parent.ID {
		t.Fatal("expect ParentID set")
	}
	if got := child.GetParentID(); got != parent.ID {
		t.Fatalf("expect %s, got %s", parent.ID, got)
	}

	// 只有 Parent 引用、ParentID 为空时走 Parent.ID 分支
	onlyRef := NewEvent("l")
	onlyRef.Parent = parent
	onlyRef.ParentID = ""
	if got := onlyRef.GetParentID(); got != parent.ID {
		t.Fatalf("expect %s, got %s", parent.ID, got)
	}

	// 两者都为空返回空字符串
	orphan := NewEvent("l")
	orphan.ParentID = ""
	orphan.Parent = nil
	if got := orphan.GetParentID(); got != "" {
		t.Fatalf("expect empty, got %s", got)
	}
}

func TestEvent_Json(t *testing.T) {
	e := NewEvent("lock-1").
		SetStorageName("s").
		SetOwnerId("o").
		SetType(EventTypeLock).
		AddActionByName("a1")

	s, err := e.ToJsonStringE()
	if err != nil {
		t.Fatalf("ToJsonStringE error: %v", err)
	}
	if !strings.Contains(s, "lock-1") {
		t.Fatalf("expect json contains lock id, got %s", s)
	}
	// ToJsonString 与 ToJsonStringE 结果一致
	if got := e.ToJsonString(); got != s {
		t.Fatal("expect ToJsonString equals ToJsonStringE")
	}

	restored, err := EventFromJsonStringE(s)
	if err != nil {
		t.Fatalf("EventFromJsonStringE error: %v", err)
	}
	if restored.LockId != "lock-1" || restored.ID != e.ID || restored.RootID != e.RootID {
		t.Fatalf("unexpected restored event: %+v", restored)
	}
	if restored.EventType != EventTypeLock {
		t.Fatalf("expect event type lock, got %v", restored.EventType)
	}
	if len(restored.Actions) != 1 || restored.Actions[0].Name != "a1" {
		t.Fatalf("unexpected restored actions: %+v", restored.Actions)
	}

	// 非法 JSON 返回错误
	if _, err := EventFromJsonStringE("{invalid"); err == nil {
		t.Fatal("expect error for invalid json")
	}
}
