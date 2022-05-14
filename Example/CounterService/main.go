package main

import (
	"dev.hijgo.go-bloc/bloc"
	"dev.hijgo.go-bloc/event"
	stream_builder "dev.hijgo.go-bloc/stream-builder"
	"fmt"
	"net/http"
	"strconv"
)

type CounterEventType uint8

const (
	addEvent      CounterEventType = 0
	resetEvent    CounterEventType = 1
	subtractEvent CounterEventType = 2
	setToEvent    CounterEventType = 3
)

type CounterEvent struct {
	EventType CounterEventType
	Input     int
}

type CounterState int

type CounterBloCData struct {
	CurrentCount int
}

func main() {

	/** First we need to define a function that will accept an incoming CounterEvent and a CounterBloCData ptr and map it
	to a CounterState.
	The logic of this function should be self-explanatory. The function will add, subtract, reset the CurrentCount field
	or set it to a specified number.
	Lastly it will return the modified CurrentCount field as a CounterState.
	*/
	mapCounterEventToCounterState := func(NewEvent event.Event[CounterEvent], BloCData *CounterBloCData) CounterState {
		switch NewEvent.Data.EventType {
		case addEvent:
			BloCData.CurrentCount += NewEvent.Data.Input
		case subtractEvent:
			BloCData.CurrentCount -= NewEvent.Data.Input
		case setToEvent:
			BloCData.CurrentCount = NewEvent.Data.Input
		case resetEvent:
			BloCData.CurrentCount = 0
		}
		return CounterState(BloCData.CurrentCount)
	}

	/** Now we are creating the Bloc, in doing so we need to pass the type of the events that should be accepted, the type
	of state that will be produced and the type of data the bloc will work with.
	In addition to that, we are also passing an initial state of out BloCData type as a starting point for the Bloc, as well
	as the mapCounterEventToCounterState to handle incoming events.
	*/
	counterBloc := bloc.CreateBloC[CounterEvent, CounterState, CounterBloCData](CounterBloCData{CurrentCount: 0}, mapCounterEventToCounterState)

	// Function that will expose an endpoint. Will try to map incoming requests to an event. Not really relevant to the
	// BloC lib.
	handleEventRequest := func(path string, eventType CounterEventType, counterBloc bloc.BloC[CounterEvent, CounterState, CounterBloCData]) {
		http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			if input, err := strconv.Atoi(req.Header.Get("Input")); err == nil {
				// Add an event to the Bloc. Will result in a new state.
				counterBloc.AddEvent(CounterEvent{
					EventType: eventType,
					Input:     input,
				})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		})
	}

	handleEventRequest("/add", addEvent, counterBloc)
	handleEventRequest("/subtract", subtractEvent, counterBloc)
	handleEventRequest("/setTo", setToEvent, counterBloc)

	http.HandleFunc("/reset", func(w http.ResponseWriter, req *http.Request) {
		counterBloc.AddEvent(CounterEvent{
			EventType: resetEvent,
			Input:     0,
		})
	})

	// This function will be called everytime a new state was produced
	buildFromState := func(state CounterState) {
		println(fmt.Sprintf("New State is: '%d'", state))
	}

	/** Initialize the StreamBuilder. We are passing the BloC, an initial event to get things started and the
	buildFromState function.*/
	stream_builder.InitStreamBuilder[CounterEvent, CounterState, CounterBloCData](counterBloc, &CounterEvent{
		EventType: resetEvent,
		Input:     0,
	}, buildFromState)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}
}
