package httpstreambuilder

import (
	"fmt"
	"net/http"

	"github.com/hijgo/go-bloc/bloc"
)

var (
	defaultHttpHandler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "No State")
	}
)

// Wrap around structure for the Business Logic Component, that will simplify the BloC experience.
// Used when the buildFunction should be exposed under a Http - Endpoint
//
// E : Type of events being emitted into the BloC
//
// S : Type of states being produced by the BloC from incoming events
//
// BD : BloCData Type of data that will be available to function that produces new state, can be used for example
//
// BloC : The BloC structure that should be wrapped
type HttpStreamBuilder[E any, S any, BD any] struct {
	BloC         bloc.BloC[E, S, BD]
	initialEvent *E
	builderFunc  func(S) func(http.ResponseWriter, *http.Request)
	httpHandler  func(http.ResponseWriter, *http.Request)
}

// Function that should be called if a new HttpStreamBuilder is needed.
//
// Will create all necessary values so the HttpStreamBuilder can function properly and then return the new HttpStreamBuilder
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
func CreateHttpStreamBuilder[E any, S any, BD any](BloC bloc.BloC[E, S, BD], InitialEvent *E, BuildFunc func(S) func(http.ResponseWriter, *http.Request)) HttpStreamBuilder[E, S, BD] {
	return HttpStreamBuilder[E, S, BD]{
		BloC:         BloC,
		initialEvent: InitialEvent,
		builderFunc:  BuildFunc,
		httpHandler:  defaultHttpHandler,
	}
}

func (sB *HttpStreamBuilder[E, S, BD]) Init(Pattern string, Mux *http.ServeMux) error {
	err := sB.BloC.StartListenToEventStream()
	if err != nil {
		return err
	}
	err = sB.BloC.ListenOnNewState(func(NewState S) {
		sB.httpHandler = sB.builderFunc(NewState)
	})
	sB.BloC.AddEvent(*sB.initialEvent)

	(*Mux).HandleFunc(Pattern, func(w http.ResponseWriter, r *http.Request) {
		sB.httpHandler(w, r)
	})

	return err
}

// If the StreamBuilder is no longer needed call this function to clear it gracefully
func (sB *HttpStreamBuilder[E, S, AD]) Dispose() {
	err := sB.BloC.StopListenToEventStream()
	if err != nil {
		panic(err)
	}
	err = sB.BloC.StopListenToStateStream()
	if err != nil {
		panic(err)
	}
}
