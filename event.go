package events

import (
	"context"
	"encoding/json"
	"github.com/golang-infrastructure/go-pointer"
	"github.com/storage-lock/go-storage"
	"github.com/storage-lock/go-utils"
	"time"
)

// Event 表示一个事件
type Event struct {

	// 事件的唯一ID
	ID string `json:"id"`

	// 此事件树的树根ID，用于对事件进行聚合分组
	RootID string `json:"root_id"`

	// 如果此事件拥有父事件，则保留一个指向父事件的引用
	ParentID string `json:"parent_id"`
	Parent   *Event `json:"-"`

	// 事件可以绑定到一个锁
	LockId string `json:"lock_id"`

	// 当前的事件是由谁产生的
	OwnerId string `json:"owner_id"`

	// 事件可以绑定到一个Storage
	StorageName string `json:"storage_name"`

	// 事件拥有开始时间和结束时间，可以通过这个来借助事件进行一些性能统计之类的
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	// 事件的类型
	EventType EventType `json:"event_type"`

	// 事件持续过程中的一些行为记录
	Actions []*Action `json:"actions"`

	// 事件可以绑定一个WatchDog
	WatchDogId string `json:"watch_dog_id"`

	// 事件可以绑定到锁信息
	LockInformation *storage.LockInformation `json:"lock_information"`

	// 如果有事件级别的错误，则可以放在这里，不过还是推荐优先使用Action级别的错误对Action进行精准绑定
	Err error `json:"err"`

	// 订阅了此时间的Listener都有哪些
	Listeners []Listener `json:"-"`
}

// EventIDPrefix 事件的ID的前缀，用于方便区分是否是事件的ID
const EventIDPrefix = "storage-lock-event-"

func NewEvent(lockId string) *Event {
	id := utils.RandomID(EventIDPrefix)
	return &Event{
		ID: id,
		// 根事件的RootID是它自己
		RootID:    id,
		LockId:    lockId,
		StartTime: pointer.Now(),
	}
}

// End 用于设置事件的完成时间为当前
func (x *Event) End() *Event {
	x.EndTime = pointer.Now()
	return x
}

// Fork 从当前事件创建一个子事件，子事件会继承父事件的一些属性
func (x *Event) Fork() *Event {
	return &Event{
		ID:              utils.RandomID(EventIDPrefix),
		RootID:          x.RootID,
		ParentID:        x.ID,
		Parent:          x,
		LockId:          x.LockId,
		StorageName:     x.StorageName,
		StartTime:       pointer.Now(),
		EventType:       x.EventType,
		LockInformation: x.LockInformation,
		Listeners:       x.Listeners,
	}
}

func (x *Event) SetRootID(rootID string) *Event {
	x.RootID = rootID
	return x
}

func (x *Event) SetStorageName(storageName string) *Event {
	x.StorageName = storageName
	return x
}

func (x *Event) SetLockId(lockId string) *Event {
	x.LockId = lockId
	return x
}

func (x *Event) SetOwnerId(ownerId string) *Event {
	x.OwnerId = ownerId
	return x
}

func (x *Event) SetType(eventType EventType) *Event {
	x.EventType = eventType
	return x
}

func (x *Event) SetErr(err error) *Event {
	x.Err = err
	return x
}

// AppendAction 往事件中追加一个action
func (x *Event) AppendAction(action *Action) *Event {
	x.Actions = append(x.Actions, action)
	return x
}

// AppendActionByName 添加Action名，一般用于只需要记录名称和时间不需要记录其它字段的时候
func (x *Event) AppendActionByName(actionName string) *Event {
	x.Actions = append(x.Actions, NewAction(actionName))
	return x
}

func (x *Event) SetListeners(listeners []Listener) *Event {
	x.Listeners = listeners
	return x
}

func (x *Event) AddListeners(listener Listener) *Event {
	x.Listeners = append(x.Listeners, listener)
	return x
}

func (x *Event) ClearListeners() *Event {
	x.Listeners = nil
	return x
}

func (x *Event) SetWatchDogId(watchDogId string) *Event {
	x.WatchDogId = watchDogId
	return x
}

func (x *Event) SetLockInformation(lockInformation *storage.LockInformation) *Event {
	x.LockInformation = lockInformation
	return x
}

//// Pub Publish的短名称
//func (x *Event) Pub(ctx context.Context, listeners ...Listener) {
//	x.Publish(ctx, listeners...)
//}

// Publish 把当前的事件发布到多个Listener上
func (x *Event) Publish(ctx context.Context, listeners ...Listener) {

	// 如果要发布的时候没有设置过结束时间，则自动设置
	if x.EndTime.IsZero() {
		x.End()
	}

	if len(x.Listeners) != 0 {
		for _, listener := range x.Listeners {
			listener.On(ctx, x)
		}
	}

	if len(listeners) != 0 {
		for _, listener := range listeners {
			listener.On(ctx, x)
		}
	}

}

// IsRootEvent 此事件是否是最底层的事件，整个事件的fork关系可以看做是一颗树
func (x *Event) IsRootEvent() bool {
	return x.Parent == nil && x.ParentID == ""
}

func (x *Event) SetParent(event *Event) *Event {
	x.ParentID = event.ID
	x.Parent = event
	return x
}

func (x *Event) GetParentID() string {

	if x.ParentID != "" {
		return x.ParentID
	}

	if x.Parent != nil {
		return x.Parent.ID
	}

	return ""
}

func (x *Event) ToJsonStringE() (string, error) {
	bytes, err := json.Marshal(x)
	return string(bytes), err
}

func (x *Event) ToJsonString() string {
	s, _ := x.ToJsonStringE()
	return s
}

func EventFromJsonStringE(eventJsonString string) (*Event, error) {
	r := &Event{}
	err := json.Unmarshal([]byte(eventJsonString), &r)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
