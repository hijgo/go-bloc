package streambuilder

import (
	"fmt"
	"net/http"

	"github.com/hijgo/go-bloc/bloc"
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
type Http[E any, S any, BD any, Out any] struct {
	BloC        bloc.BloC[E, S, BD]
	Mux         *http.ServeMux
	pattern     string
	builderFunc func(S, **Out)
}

// Function that should be called if a new HttpHttp is needed.
//
// Will create all necessary values so the HttpHttp can function properly and then return the new HttpHttp
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
func CreateHttp[E any, S any, BD any, Out any](BloC bloc.BloC[E, S, BD], BuildFunc func(S, **Out), Mux *http.ServeMux, Pattern string) Http[E, S, BD, Out] {
	return Http[E, S, BD, Out]{
		BloC:        BloC,
		builderFunc: BuildFunc,
		Mux: func() *http.ServeMux {
			if Mux != nil {
				return Mux
			}
			return http.DefaultServeMux
		}(),
		pattern: Pattern,
	}
}

func (sB Http[E, S, BD, Out]) Init(InitialEvent *E) error {

	if err := sB.BloC.StartListenToEventStream(); err != nil {
		return err
	}

	var body *Out

	if err := sB.BloC.ListenOnNewState(func(NewState S) {
		sB.builderFunc(NewState, &body)
	}); err != nil {
		return err
	}

	(*sB.Mux).HandleFunc(sB.pattern, func(w http.ResponseWriter, r *http.Request) {
		if body == nil {
			defaultHttpHandler(w, r)
			return
		}
		fmt.Fprint(w, *body)
	})

	sB.BloC.AddEvent(*InitialEvent)

	return nil
}

// If the Http - StreamBuilder is no longer needed call this function to clear it gracefully
func (sB Http[E, S, BD, Out]) Dispose() error {
	return sB.BloC.Dispose()
}
