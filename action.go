package events

import (
	"encoding/json"
	"errors"
	"github.com/golang-infrastructure/go-pointer"
	"math"
	"strconv"
	"time"
)

// Action 一个事件可以添加若干个Action，每个Action都会有一些时间、名称、上下文之类的
type Action struct {

	// Action被创建的时间
	StartTime *time.Time `json:"start_time"`

	// Action结束的时间，有时候会用不到
	EndTime *time.Time `json:"end_time"`

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

// End 用于标记action结束了，更新其结束时间
func (x *Action) End() *Action {
	x.EndTime = pointer.Now()
	return x
}

// Cost 计算此Action花费的时间，必须startTime和endTime同时存在的时候才有效
func (x *Action) Cost() time.Duration {
	if x.StartTime == nil || x.EndTime == nil {
		return 0
	}
	return x.EndTime.Sub(*x.StartTime)
}

// SetName 设置action的名字
func (x *Action) SetName(name string) *Action {
	x.Name = name
	return x
}

// SetErr 设置action绑定的错误
func (x *Action) SetErr(err error) *Action {
	x.Err = err
	return x
}

// ClearErr 清空action绑定的错误
func (x *Action) ClearErr() *Action {
	x.Err = nil
	return x
}

// GetErrMsg 获取action绑定的错误的信息，如果没有绑定错误的话返回空字符串
func (x *Action) GetErrMsg() string {
	if x.Err != nil {
		return x.Err.Error()
	} else {
		return ""
	}
}

// ErrorIs 判断action绑定的错误是否是给定类型
func (x *Action) ErrorIs(err error) bool {
	if x.Err == nil {
		return false
	}
	return errors.Is(x.Err, err)
}

// SetPayloadMap 设置action的payload map
func (x *Action) SetPayloadMap(payloadMap map[string]any) *Action {
	x.PayloadMap = payloadMap
	return x
}

// AddPayload 往action上增加payload
func (x *Action) AddPayload(key string, value any) *Action {
	if x.PayloadMap == nil {
		x.PayloadMap = make(map[string]any, 0)
	}

	x.PayloadMap[key] = value

	return x
}

// ClearPayloadMap 清空payloadMap
func (x *Action) ClearPayloadMap() *Action {
	x.PayloadMap = nil
	return x
}

func (x *Action) GetPayloadMap() map[string]any {
	return x.PayloadMap
}

// GetPayload 从payloadMap中获取value
func (x *Action) GetPayload(key string) (any, bool) {
	if x.PayloadMap == nil {
		return nil, false
	}
	v, exists := x.PayloadMap[key]
	return v, exists
}

// GetPayloadAsString 从payloadMap中获取value并以string返回。
// 对于序列化后常见的字符串、数字和布尔值会转换为其文本表示；
// 缺少key或值类型不支持时返回空字符串。
func (x *Action) GetPayloadAsString(key string) string {
	payload, ok := x.GetPayload(key)
	if !ok || payload == nil {
		return ""
	}
	switch value := payload.(type) {
	case string:
		return value
	case json.Number:
		return value.String()
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64)
	default:
		return ""
	}
}

// GetPayloadAsInt 获取payload作为int返回。
// 支持整数类型和可解析的十进制数字字符串；缺少key、溢出或类型不支持时返回0。
func (x *Action) GetPayloadAsInt(key string) int {
	payload, ok := x.GetPayload(key)
	if !ok || payload == nil {
		return 0
	}
	switch value := payload.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		if int64(int(value)) != int64(value) {
			return 0
		}
		return int(value)
	case int64:
		if int64(int(value)) != value {
			return 0
		}
		return int(value)
	case uint:
		if uint(int(value)) != value {
			return 0
		}
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		if uint64(int(value)) != uint64(value) {
			return 0
		}
		return int(value)
	case uint64:
		if int64(int(value)) < 0 || uint64(int(value)) != value {
			return 0
		}
		return int(value)
	case json.Number:
		parsed, err := value.Int64()
		if err == nil {
			if int64(int(parsed)) != parsed {
				return 0
			}
			return int(parsed)
		}
		floating, err := value.Float64()
		if err != nil || math.IsNaN(floating) || math.IsInf(floating, 0) || math.Trunc(floating) != floating {
			return 0
		}
		intLimit := math.Ldexp(1, strconv.IntSize-1)
		if floating < -intLimit || floating >= intLimit {
			return 0
		}
		return int(floating)
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}
