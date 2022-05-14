package bloc

import (
	"dev.hijgo.go-bloc/event"
	"dev.hijgo.go-bloc/stream"
	"reflect"
	"testing"
	"time"
)

type Event struct {
	Data int
}

type BD struct {
	BD string
}

type State struct {
	State int
}

func TestCreateBloC(t *testing.T) {
	bd := BD{}
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{} })

	if value := reflect.TypeOf(b.eventStream); value != reflect.TypeOf(&stream.Stream[event.Event[Event]]{}) {
		t.Errorf("Expected eventStream Of Type '%s' Actual '%s'", reflect.TypeOf(&stream.Stream[event.Event[Event]]{}), value)
	}

	if value := reflect.TypeOf(b.stateStream); value != reflect.TypeOf(&stream.Stream[State]{}) {
		t.Errorf("Expected stateStream Of Type '%s' Actual '%s'", reflect.TypeOf(&stream.Stream[State]{}), value)
	}

	if value := reflect.TypeOf(b.mapEventToState); value != reflect.TypeOf(func(event.Event[Event], *BD) State { return State{} }) {
		t.Errorf("Expected mapEventToState Of Type '%s' Actual '%s'", reflect.TypeOf(func(event.Event[Event], BD) {}), value)
	}

	if value := reflect.TypeOf(b.BloCData); value != reflect.TypeOf(BD{}) {
		t.Errorf("Expected AdditionalData Of Type '%s' Actual '%s'", reflect.TypeOf(BD{}), value)
	}

	if value := b.BloCData; value != bd {
		t.Errorf("Expected AdditionalData Of Value '%s' Actual '%s'", bd, value)
	}
}

func TestBloC_StartListenToEventStream(t *testing.T) {
	bd := BD{}
	value := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		value += NewItem.Data.Data
	}

	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))

	if value != 0 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 0, value)
	}

	b.StartListenToEventStream()

	if returnValue := b.StartListenToEventStream(); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 2}))

	if value != 2 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 2, value)
	}
}

func TestBloC_StopListenToEventStream(t *testing.T) {
	bd := BD{}
	check := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		check += NewItem.Data.Data
	}

	b.StartListenToEventStream()
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))

	if check != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, check)
	}

	b.StopListenToEventStream()
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 2}))

	if check != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, check)
	}
}

func TestBloC_ListenOnNewState(t *testing.T) {
	bd := BD{}
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{State: 2 * E.Data.Data} })
	check := 0

	b.StartListenToEventStream()
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))

	b.ListenOnNewState(func(S State) {
		check = S.State
	})

	if returnValue := b.ListenOnNewState(func(S State) { check = S.State }); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	if check != 0 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 0, check)
	}

	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))
	time.Sleep(1 * time.Millisecond)
	if check != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, check)
	}
}

func TestBloC_StopListenToStateStream(t *testing.T) {
	bd := BD{}
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{State: 2 * E.Data.Data} })
	check := 0

	b.StartListenToEventStream()
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))

	b.ListenOnNewState(func(S State) {
		check = S.State
	})

	if returnValue := b.ListenOnNewState(func(S State) { check = S.State }); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	if check != 0 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 0, check)
	}

	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 1}))
	time.Sleep(1 * time.Millisecond)
	if check != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, check)
	}
	b.StopListenToStateStream()
	b.eventStream.Add(event.CreateEvent[Event](Event{Data: 2}))
	time.Sleep(1 * time.Millisecond)
	if check != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, check)
	}
}

func TestBloC_AddEvent(t *testing.T) {
	bd := BD{}
	check := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { check += E.Data.Data; return State{} })
	b.StartListenToEventStream()
	b.AddEvent(Event{Data: 1})
	if check != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, check)
	}
}
