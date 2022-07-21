package stream_builder

// An interface to unify the interactions with all Stream-Builders
type StreamBuilder[E any, S any, BD any] interface {
	// A function that sould shutdown the StreamBuilder gracefuly
	//
	// Will return an error when not successful
	Dispose() error
	// For initializing the structure.
	//
	// After this function was called, the bloc should be able to receive and process incomming and build a state from it.
	//
	// Will return an error when not successful
	Init(initialEvent *E) error
}
