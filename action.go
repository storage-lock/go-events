package events

import (
	"github.com/golang-infrastructure/go-pointer"
	"time"
)

// Action 一个事件可以添加若干个Action，每个Action都会有一些时间、名称、上下文之类的
type Action struct {

	// Action被创建的时间
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	// Action的名字
	Name string `json:"name"`

	// 执行此Action时发生的错误
	Err error `json:"err"`

	// action可以携带一些自己单独的上下文信息之类的
	PayloadMap map[string]any `json:"payload_map"`
}

func NewAction(name string) *Action {
	return &Action{
		Name:      name,
		StartTime: pointer.Now(),
	}
}

func (x *Action) End() *Action {
	x.EndTime = pointer.Now()
	return x
}

// Cost 计算此Action花费的时间
func (x *Action) Cost() time.Duration {
	if x.StartTime == nil || x.EndTime == nil {
		return 0
	}
	return x.EndTime.Sub(*x.StartTime)
}

func (x *Action) SetName(name string) *Action {
	x.Name = name
	return x
}

func (x *Action) SetErr(err error) *Action {
	x.Err = err
	return x
}

func (x *Action) SetPayloadMap(payloadMap map[string]any) *Action {
	x.PayloadMap = payloadMap
	return x
}

func (x *Action) AddPayload(key string, value any) *Action {
	if x.PayloadMap == nil {
		x.PayloadMap = make(map[string]any, 0)
	}

	x.PayloadMap[key] = value

	return x
}
