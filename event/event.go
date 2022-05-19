package event

import "time"

// Event
// Structure that will be passed into a sink of a stream.
// T - The type of data carried with by the event struct.
// TimeStamp - Time of creation
// Data - Additional data associated with the event.
type Event[T any] struct {
	TimeStamp int64
	Data      T
}

// CreateEvent
// Function that will create a new Event[T] and then return it.
// T - The type of data carried with by the event struct.#
// Data - Additional data associated with the new event.
func CreateEvent[T any](Data T) Event[T] {
	return Event[T]{
		TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
		Data:      Data,
	}
}
