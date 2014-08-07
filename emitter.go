package goble

import (
	"log"
)

const (
	ALL = "__allEvents__"
)

// Event generated by blued, with associated data
type Event struct {
	Name        string
	State       string
	DeviceUUID  UUID
	ServiceUuid string
	Peripheral  Peripheral
}

// The event handler function.
// Return true to terminate
type EventHandlerFunc func(Event) bool

// Emitter is an object to emit and handle Event(s)
type Emitter struct {
	handlers map[string]EventHandlerFunc
	event    chan Event
}

// Init initialize the emitter and start a goroutine to execute the event handlers
func (e *Emitter) Init() {
	e.handlers = make(map[string]EventHandlerFunc)
	e.event = make(chan Event)

	// event handler
	go func() {
		for {
			ev := <-e.event

			if fn, ok := e.handlers[ev.Name]; ok {
				if fn(ev) {
					break
				}
			} else if fn, ok := e.handlers[ALL]; ok {
				if fn(ev) {
					break
				}
			} else {
				log.Println("unhandled Emit", ev)
			}
		}

		close(e.event)
	}()
}

// Emit sends the event on the 'event' channel
func (e *Emitter) Emit(ev Event) {
	e.event <- ev
}

// On(event, cb) registers an handler for the specified event
func (e *Emitter) On(event string, fn EventHandlerFunc) {
	if fn == nil {
		delete(e.handlers, event)
	} else {
		e.handlers[event] = fn
	}
}
