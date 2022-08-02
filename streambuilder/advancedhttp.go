package streambuilder

import (
	"net/http"

	"github.com/hijgo/go-bloc/bloc"
)

// Wrap around structure for the Business Logic Component, that will simplify the BloC experience.
// Used when the buildFunction should be exposed under a Http-Endpoint
//
// AdvancedHttp provides detailed behaviour configuration of the Http-Endpoint.
//
// E : Type of events being emitted into the BloC
//
// S : Type of states being produced by the BloC from incoming events
//
// BD : BloCData Type of data that will be available to function that produces new state, can be used for example
//
// BloC : The BloC structure that should be wrapped
type AdvancedHttp[E any, S any, BD any] struct {
	BloC        bloc.BloC[E, S, BD]
	Mux         *http.ServeMux
	pattern     string
	builderFunc func(S) func(http.ResponseWriter, *http.Request)
	httpHandler func(http.ResponseWriter, *http.Request)
}

// Function that should be called if a new AdvancedHttp is needed.
//
// Will create all necessary values so the AdvancedHttp can function properly and then return the new AdvancedHttp
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
func CreateAdvancedHttp[E any, S any, BD any](BloC bloc.BloC[E, S, BD], BuildFunc func(S) func(http.ResponseWriter, *http.Request), Mux *http.ServeMux, Pattern string) AdvancedHttp[E, S, BD] {
	return AdvancedHttp[E, S, BD]{
		BloC:        BloC,
		builderFunc: BuildFunc,
		httpHandler: defaultHttpHandler,
		Mux: func() *http.ServeMux {
			if Mux != nil {
				return Mux
			}
			return http.DefaultServeMux
		}(),
		pattern: Pattern,
	}
}

func (sB AdvancedHttp[E, S, BD]) Init(InitialEvent *E) error {

	if err := sB.BloC.StartListenToEventStream(); err != nil {
		return err
	}

	if err := sB.BloC.ListenOnNewState(func(NewState S) {
		sB.httpHandler = sB.builderFunc(NewState)
	}); err != nil {
		return err
	}

	(*sB.Mux).HandleFunc(sB.pattern, func(w http.ResponseWriter, r *http.Request) {
		sB.httpHandler(w, r)
	})

	sB.BloC.AddEvent(*InitialEvent)

	return nil
}

// If the AdvancedHttp - StreamBuilder is no longer needed call this function to clear it gracefully
func (sB AdvancedHttp[E, S, AD]) Dispose() error {
	return sB.BloC.Dispose()
}
