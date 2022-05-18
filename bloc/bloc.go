package bloc

import (
	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
)

var DefaultMaxQueueSize = 100

type BloC[E any, S any, BD any] struct {
	eventStream     *stream.Stream[event.Event[E]]
	stateStream     *stream.Stream[S]
	BloCData        BD
	mapEventToState func(NewEvent event.Event[E], AdditionalData *BD) S
}

func CreateBloC[E any, S any, BD any](InitialBloCData BD, mapEventToState func(NewEvent event.Event[E], BloCData *BD) S) BloC[E, S, BD] {
	newBloC := BloC[E, S, BD]{
		BloCData:        InitialBloCData,
		mapEventToState: mapEventToState,
	}
	stateStream := stream.CreateStream[S](DefaultMaxQueueSize, func(NewItem S) {})
	eventStream := stream.CreateStream[event.Event[E]](DefaultMaxQueueSize, func(NewEvent event.Event[E]) {
		stateStream.Add(mapEventToState(NewEvent, &newBloC.BloCData))
	})
	newBloC.stateStream = &stateStream
	newBloC.eventStream = &eventStream
	return newBloC
}

func (b *BloC[E, S, AD]) AddEvent(NewEvent E) {
	b.eventStream.Add(event.CreateEvent[E](NewEvent))
}

func (b *BloC[E, S, AD]) ListenOnNewState(OnNewState func(S)) bool {
	if b.stateStream.GetListenStatus() {
		return false
	}
	b.stateStream.OnNewItem = func(NewState S) {
		OnNewState(NewState)
	}
	b.stateStream.Listen()
	return true
}

func (b *BloC[E, S, AD]) StopListenToStateStream() {
	b.stateStream.StopListen()
}

func (b *BloC[E, S, AD]) StartListenToEventStream() bool {
	if b.eventStream.GetListenStatus() {
		return false
	}
	b.eventStream.Listen()
	return true
}

func (b *BloC[E, S, AD]) StopListenToEventStream() {
	b.eventStream.StopListen()
}
