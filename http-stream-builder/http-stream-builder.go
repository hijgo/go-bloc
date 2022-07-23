package httpstreambuilder

import (
	"fmt"
	"net/http"

	"github.com/hijgo/go-bloc/bloc"
)

func defaultHttpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "No state was created yet")
	w.WriteHeader(204)
}

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
	BloC        bloc.BloC[E, S, BD]
	Mux         *http.ServeMux
	pattern     string
	builderFunc func(S) func(http.ResponseWriter, *http.Request)
	httpHandler func(http.ResponseWriter, *http.Request)
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
// BuildFunc : The function that will handle any new produced state
func CreateHttpStreamBuilder[E any, S any, BD any](BloC bloc.BloC[E, S, BD], BuildFunc func(S) func(http.ResponseWriter, *http.Request), Mux *http.ServeMux, Pattern string) HttpStreamBuilder[E, S, BD] {
	return HttpStreamBuilder[E, S, BD]{
		BloC:        BloC,
		builderFunc: BuildFunc,
		httpHandler: defaultHttpHandler,
		Mux:         Mux,
		pattern:     Pattern,
	}
}

func (sB HttpStreamBuilder[E, S, BD]) Init(InitialEvent *E) error {

	if err := sB.BloC.StartListenToEventStream(); err != nil {
		return err
	}

	if err := sB.BloC.ListenOnNewState(func(NewState S) {
		sB.httpHandler = sB.builderFunc(NewState)
	}); err != nil {
		return err
	}

	sB.BloC.AddEvent(*InitialEvent)

	if sB.Mux != nil {
		(*sB.Mux).HandleFunc(sB.pattern, func(w http.ResponseWriter, r *http.Request) {
			sB.httpHandler(w, r)
		})
	} else {
		http.DefaultServeMux.HandleFunc(sB.pattern, func(w http.ResponseWriter, r *http.Request) {
			sB.httpHandler(w, r)
		})
	}

	return nil
}

// If the StreamBuilder is no longer needed call this function to clear it gracefully
func (sB HttpStreamBuilder[E, S, AD]) Dispose() error {
	return sB.BloC.Dispose()
}
