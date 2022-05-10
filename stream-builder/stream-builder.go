package stream_builder

import (
	"dev.hijgo.go-bloc/bloc"
)

type StreamBuilder[E any, S any, AD any] struct {
	BloC         bloc.BloC[E, S, AD]
	InitialEvent E
	BuilderFunc  func(S)
}

func InitStreamBuilder[E any, S any, AD any](bloc bloc.BloC[E, S, AD], InitialEvent E, buildFunc func(S)) StreamBuilder[E, S, AD] {
	streamBuilder := StreamBuilder[E, S, AD]{
		BloC:         bloc,
		InitialEvent: InitialEvent,
		BuilderFunc:  buildFunc,
	}
	streamBuilder.BloC.StartListenToEventStream()
	streamBuilder.BloC.ListenOnNewState(buildFunc)
	streamBuilder.BloC.AddEvent(InitialEvent)
	return streamBuilder
}

func (sB *StreamBuilder[E, S, AD]) Dispose() {
	sB.BloC.StartListenToEventStream()
	sB.BloC.StopListenToStateStream()
}