package engine

import (
	"time"
)

// EventHandler declares the basic event interface
/*type EventHandler interface {
	Timer
}*/

// Timer declares the timer interface
type EventHandler interface {
	Time() time.Time
	SetTime(time.Time)
}



// Event is the implementation of the basic event interface.
type Event struct {
	timestamp time.Time
}

// Time returns the timestamp of an event
func (e Event) Time() time.Time {
	return e.timestamp
}

// SetTime returns the timestamp of an event
func (e *Event) SetTime(t time.Time) {
	e.timestamp = t
}



// SignalEvent declares the signal event interface.
type SignalEvent interface {
	EventHandler
	Directioner
}

// Directioner defines a direction interface
type Directioner interface {
	Direction() Direction
	//SetDirection(Direction)
	SetOrderType(orderType OrderType)
	GetOrderType()(orderType OrderType)
}


// Quantifier defines a Qty interface.
type Quantifier interface {
	Qty() int64
	SetQty(int64)
}

// IDer declares setting and retrieving of an Id.
type IDer interface {
	ID() int
	SetID(int)
}


