package basicstreambuilder

import (
	"github.com/hijgo/go-bloc/bloc"
)

// Wrap around structure for the Business Logic Component, that will simplify the BloC experience.
//
// E: Type of events being emitted into the BloC
//
// S: Type of states being produced by the BloC from incoming events
//
// BD: BloCData Type of data that will be available to function that produces new state, can be used for example
//
// BloC: The BloC structure that should be wrapped
type StreamBuilder[E any, S any, BD any] struct {
	BloC        bloc.BloC[E, S, BD]
	builderFunc func(S)
}

// Will create all necessary values so the StreamBuilder can function properly and then return the new StreamBuilder
// of type E,S,BD or/and an error.
//
// E : Type of events being emitted into the BloC
//
// S : Type of states being produced by the BloC from incoming events
//
// BD : BloCData Type of data that will be available to function that produces new state, can be used for example
// to store additional data not originating from events or store event specific temporally use later.
//
// BloC : The BloC structure that should be wrapped
//
// InitialEvent : A start event of type E start will kick off things and as a result will create an initial state of type S
//
// BuildFunc : The function that will handle any new produced state
func CreateBasicStreamBuilder[E any, S any, BD any](BloC bloc.BloC[E, S, BD], BuildFunc func(S)) StreamBuilder[E, S, BD] {
	return StreamBuilder[E, S, BD]{
		BloC:        BloC,
		builderFunc: BuildFunc,
	}

}

func (sB *StreamBuilder[E, S, BD]) Init(initialEvent *E) error {
	if err := sB.BloC.StartListenToEventStream(); err != nil {
		return err
	}

	if err := sB.BloC.ListenOnNewState(sB.builderFunc); err != nil {
		return err
	}
	sB.BloC.AddEvent(*initialEvent)
	return nil
}

// If the StreamBuilder is no longer needed call this function to clear it gracefully
func (sB *StreamBuilder[E, S, AD]) Dispose() {
	sB.BloC.Dispose()
}
