package streambuilder

import (
	"net/http"
)

const (
	defaultPath = "/"
)

// Interface to unify the interactions with all Stream-Builders.
// Where a Stream-Builder should be a structure that wraps a already
// existing BloC-Structure, listens to it's state-stream and
// further uses the output of stream.
type StreamBuilder[E any, S any, BD any] interface {
	// Shutdown the StreamBuilder and it's underlying
	// BloC-Structure gracefully. After disposing the BloC
	// cannot be used again.
	//
	// Will return error when unsuccessful.
	Dispose() error
	// Set the starting Event of the bloc`s event-stream.
	// Also will try to listen to the bloc`s state-/event-stream.
	//
	// If not succesfull will return error.
	Init(initialEvent *E) error
}

func defaultHttpHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/NoState.html")
}