package events

import "context"

// ------------------------------------------------- --------------------------------------------------------------------

// Listener 事件监听器，当有任何事件的时候都会触发此监听器把事件传给它处理
type Listener interface {

	// Name 当有多个事件监听器的时候用于区分彼此
	Name() string

	// On 每次事件被触发的时候此函数被执行一次
	On(ctx context.Context, e *Event)
}

// ------------------------------------------------- --------------------------------------------------------------------

// ListenerWrapper 如果不想声明struct实现Listener接口的话，可以使用这个struct来对函数来包裹函数实现
type ListenerWrapper struct {
	name         string
	listenerFunc func(ctx context.Context, e *Event)
}

var _ Listener = &ListenerWrapper{}

func NewListenerWrapper(name string, listenerFunc func(ctx context.Context, e *Event)) *ListenerWrapper {
	return &ListenerWrapper{
		name:         name,
		listenerFunc: listenerFunc,
	}
}

func (x *ListenerWrapper) Name() string {
	return x.name
}

func (x *ListenerWrapper) On(ctx context.Context, e *Event) {
	x.listenerFunc(ctx, e)
}

// ------------------------------------------------- --------------------------------------------------------------------
