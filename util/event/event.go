package event

import "time"

type Event[T any] struct {
	TimeStamp int64
	Data      T
}

func CreateEvent[T any](Data T) Event[T] {
	return Event[T]{
		TimeStamp: time.Now().UnixNano() / int64(time.Millisecond),
		Data:      Data,
	}
}
