package bloc

import (
	"dev.hijgo.go-bloc/event"
	"dev.hijgo.go-bloc/stream"
)

var DefaultMaxQueueSize = 100

type BloC[E any, S any, AD any] struct {
	eventStream     *stream.Stream[event.Event[E]]
	stateStream     *stream.Stream[S]
	AdditionalData  AD
	mapEventToState func(NewEvent event.Event[E], AdditionalData AD) S
}

func CreateBloC[E any, S any, AD any](InitAdditionalData AD, mapEventToState func(NewEvent event.Event[E], AdditionalData AD) S) BloC[E, S, AD] {
	newBloC := BloC[E, S, AD]{
		AdditionalData:  InitAdditionalData,
		mapEventToState: mapEventToState,
	}
	stateStream := stream.CreateStream[S](DefaultMaxQueueSize, func(NewItem S) {})
	eventStream := stream.CreateStream[event.Event[E]](DefaultMaxQueueSize, func(NewEvent event.Event[E]) {
		stateStream.Add(mapEventToState(NewEvent, newBloC.AdditionalData))
	})
	newBloC.stateStream = &stateStream
	newBloC.eventStream = &eventStream
	return newBloC
}

func (b *BloC[E, S, D]) ListenOnNewState(OnNewState func(S, D)) bool {
	if b.stateStream.GetListenStatus() {
		return false
	}
	b.stateStream.OnNewItem = func(NewState S) {
		OnNewState(NewState, b.AdditionalData)
	}
	b.stateStream.Listen()
	return true
}

func (b *BloC[E, S, D]) StopListenToStateStream() {
	b.stateStream.StopListen()
}

func (b *BloC[E, S, D]) StartListenToEventStream() bool {
	if b.eventStream.GetListenStatus() {
		return false
	}
	b.eventStream.Listen()
	return true
}

func (b *BloC[E, S, D]) StopListenToEventStream() {
	b.eventStream.StopListen()
}
