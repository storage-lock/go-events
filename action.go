package events

import (
	"github.com/golang-infrastructure/go-pointer"
	"time"
)

// 用于统计Storage上的方法调用行为
const (
	ActionStorageGetName           = "Storage.GetName"
	ActionStorageInit              = "Storage.Init"
	ActionStorageUpdateWithVersion = "Storage.UpdateWithVersion"
	ActionStorageInsertWithVersion = "Storage.InsertWithVersion"
	ActionStorageDeleteWithVersion = "Storage.DeleteWithVersion"
	ActionStorageGetTime           = "Storage.GetTime"
	ActionStorageGet               = "Storage.Get"
	ActionStorageClose             = "Storage.Close"
	ActionStorageList              = "Storage.List"
)

// Action 一个事件可以添加若干个Action，每个Action都会有一些时间、名称、上下文之类的
type Action struct {

	// Action被创建的时间
	Time *time.Time `json:"time"`

	// Action的名字，上面内置了一些，其它的可以自行定义
	Name string `json:"name"`

	// 执行此Action时发生的错误
	Err error `json:"err"`

	// action可以携带一些自己单独的上下文信息之类的
	Payload string `json:"payload"`
}

func NewAction(name string) *Action {
	return &Action{
		Time: pointer.Now(),
		Name: name,
	}
}

func (x *Action) SetTime(time *time.Time) *Action {
	x.Time = time
	return x
}

func (x *Action) SetName(name string) *Action {
	x.Name = name
	return x
}

func (x *Action) SetErr(err error) *Action {
	x.Err = err
	return x
}

func (x *Action) SetPayload(payload string) *Action {
	x.Payload = payload
	return x
}
