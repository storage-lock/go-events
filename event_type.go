package events

type EventType int

const (

	// EventTypeUnknown 未设置事件类型
	EventTypeUnknown EventType = iota

	// EventTypeCreateLock 创建锁的过程中产生的事件
	EventTypeCreateLock

	// EventTypeLock 获取锁的过程中产生的事件
	EventTypeLock

	// EventTypeUnlock 释放锁的过程中产生的事件
	EventTypeUnlock

	// EventTypeWatchDog 看门狗产生的事件
	EventTypeWatchDog
)
