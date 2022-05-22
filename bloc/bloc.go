package bloc

import (
	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
)

var DefaultMaxQueueSize = 100

// BloC
// Business Logic Component, that will provide event and state streams for further use.
//
// E - Type of events being emitted into the BloC
//
// S - Type of states being produced by the BloC from incoming events
//
// BD - BloCData Type of data that will be available to function that produces new state, can be used for example
// to store additional data not originating from events or store event specific temporally use later.
type BloC[E any, S any, BD any] struct {
	eventStream     *stream.Stream[event.Event[E]]
	stateStream     *stream.Stream[S]
	BloCData        BD
	mapEventToState func(NewEvent event.Event[E], AdditionalData *BD) S
}

// CreateBloC
// Function that should be called if a new BloC is needed.
// Will create all necessary values so the BloC can function properly and then return the new BloC of type E,S,BD.
//
// E - Type of events being emitted into the BloC
//
// S - Type of states being produced by the BloC from incoming events
//
// BD - BloCData Type of data that will be available to function that produces new state, can be used for example
// to store additional data not originating from events or store event specific temporally use later
//
// InitialBloCData - The initial BloCData struct being used by the bloc
//
// mapEventToState - Function that accepts an Event of Type event.Event[E] and BD ptr and maps the event to a new state
// of type S. This Function will be called everytime when a new event it added to the event stream
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

// AddEvent
// Should be called when a new Event should be passed to the event stream.
// Will result ultimately in a new state.
//
// NewEvent - The event of type E that should be passed to the event stream.
func (b *BloC[E, S, AD]) AddEvent(NewEvent E) {
	b.eventStream.Add(event.CreateEvent[E](NewEvent))
}

// ListenOnNewState
// Start listening to the state stream by calling the function.
//
// OnNewState - Function that must accept a new state of type S
//
// Will return an error if for example the state stream is already being listened to.
func (b *BloC[E, S, AD]) ListenOnNewState(OnNewState func(S)) error {
	b.stateStream.OnNewItem = func(NewState S) {
		OnNewState(NewState)
	}
	return b.stateStream.Listen()
}

// StopListenToStateStream
// Call to stop listen to the state stream.
//
// Will return an error if for example the stream wasn't listened to.
func (b *BloC[E, S, AD]) StopListenToStateStream() error {
	return b.stateStream.StopListen()
}

// StartListenToEventStream
// Start listening to the event stream by calling the function.
//
// When called will produce a new state for every new event passed.
//
// Will return an error if for example the state stream is already being listened to.
func (b *BloC[E, S, AD]) StartListenToEventStream() error {
	return b.eventStream.Listen()
}

// StopListenToEventStream
// Call to stop listen to the event stream.
//
// Will return an error if for example the stream wasn't listened to.
func (b *BloC[E, S, AD]) StopListenToEventStream() error {
	return b.eventStream.StopListen()
}

// Dispose
// If the BloC is no longer needed call this function to clear it gracefully
func (b *BloC[E, S, AD]) Dispose() {
	b.stateStream.Dispose()
	b.eventStream.Dispose()
}
