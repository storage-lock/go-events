# Events

# 一、这是什么

Storage Lock底层的事件机制，分布式锁算法的大多数操作都会触发事件，通过监听事件可以实现锁的可观测性等等功能，这里只是进行事件机制的基础定义，具体事件触发是在[Storage Lock](https://github.com/storage-lock/go-storage-lock)的分布式锁的各个步骤里。

# 二、安装依赖

```bash
go get -u github.com/storage-lock/go-events
```

# 三、组件介绍

## Event

用于封装每个被触发的事件：

```go
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
```

## Action

一个事件可能会携带多个Action：

```go
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
```

## Listener

实现Listener在创建锁的时候设置上则可以实现对锁事件的监听： 

```go
// Listener 事件监听器，当有任何事件的时候都会触发此监听器把事件传给它处理
type Listener interface {

	// Name 当有多个事件监听器的时候用于区分彼此
	Name() string

	// On 每次事件被触发的时候此函数被执行一次
	On(ctx context.Context, e *Event)
}
```

# 四、实现示例

基于事件机制，你可以实现很多有意思的功能，这里是一些基于事件机制实现的示例：

- https://github.com/storage-lock/go-event-listener-stdout

