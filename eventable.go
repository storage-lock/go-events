package events

// Eventable
type Eventable interface {
	SetEvent(e *Event)
}
