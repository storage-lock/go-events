package events

import "context"

// Listener 事件监听器，当有任何事件的时候都会触发此监听器把事件传给它处理
type Listener interface {
	On(ctx context.Context, e *Event)
}
