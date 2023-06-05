package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	HandleEvent(event EventInterface, wg *sync.WaitGroup) error
}

type EventDispatcherInterface interface {
	RegisterHandler(eventName string, handler EventHandlerInterface) error
	DispatchEvent(event EventInterface) error
	RemoveHandler(event EventInterface, handler EventHandlerInterface) error
	HasHandler(eventName string, handler EventHandlerInterface) bool
	ClearHandlers()
}
