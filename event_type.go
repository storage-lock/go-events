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

	// EventTypeExpired 锁因租约过期被抢占（记录是墓碑或超时后他人接管）。见 ActionLockExpired。
	EventTypeExpired

	// EventTypeBusy 竞争锁时发现锁被他人持有而等待重试。见 ActionLockBusy。
	EventTypeBusy

	// EventTypeRenewed 锁租约被看门狗或主动 Renew 成功续期。见 ActionWatchDogRefreshSuccess。
	EventTypeRenewed

	// EventTypeLost 看门狗确认持有权丢失（他人抢占、锁删除或租约过期），业务侧应停止使用该锁。
	// 对应 OnLockLost 回调与 ActionNotLockOwner/ErrLockLost 路径。
	EventTypeLost

	// EventTypeWaitChanged 锁记录上的等待队列发生变化（有人排队、离开或队首接管）。
	EventTypeWaitChanged
)

// String 返回事件类型的可读名称，供日志、审计和指标使用。
// 未知类型返回 "unknown"。
func (x EventType) String() string {
	switch x {
	case EventTypeUnknown:
		return "unknown"
	case EventTypeCreateLock:
		return "create"
	case EventTypeLock:
		return "lock"
	case EventTypeUnlock:
		return "unlock"
	case EventTypeWatchDog:
		return "watchdog"
	case EventTypeExpired:
		return "expired"
	case EventTypeBusy:
		return "busy"
	case EventTypeRenewed:
		return "renewed"
	case EventTypeLost:
		return "lost"
	case EventTypeWaitChanged:
		return "wait_changed"
	default:
		return "unknown"
	}
}
