package bloc

import (
	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
	"reflect"
	"sync"
	"testing"
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
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		value += NewItem.Data.Data
		wg.Done()
	}

	b.AddEvent(Event{Data: 1})

	if value != 0 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 0, value)
	}

	b.StartListenToEventStream()

	if returnValue := b.StartListenToEventStream(); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	wg.Add(1)
	b.AddEvent(Event{Data: 2})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 2, value)
	}
}

func TestBloC_StopListenToEventStream(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		value += NewItem.Data.Data
		wg.Done()
	}

	b.StartListenToEventStream()
	wg.Add(1)
	b.AddEvent(Event{Data: 1})
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}

	b.StopListenToEventStream()
	b.AddEvent(Event{Data: 2})

	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}
}

func TestBloC_ListenOnNewState(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State {
		return State{State: 2 * E.Data.Data}
	})
	value := 0

	b.StartListenToEventStream()
	b.AddEvent(Event{Data: 1})

	b.ListenOnNewState(func(S State) {
		value = S.State
		wg.Done()
	})

	if returnValue := b.ListenOnNewState(func(S State) { value = S.State }); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	if value != 0 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 0, value)
	}

	wg.Add(1)
	b.AddEvent(Event{Data: 1})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 2, value)
	}
}

func TestBloC_StopListenToStateStream(t *testing.T) {
	var wgState sync.WaitGroup
	var wgEvent sync.WaitGroup
	bd := BD{}
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { defer wgEvent.Done(); return State{State: 2 * E.Data.Data} })
	value := 0

	b.StartListenToEventStream()
	wgEvent.Add(1)
	b.AddEvent(Event{Data: 1})
	wgEvent.Wait()

	b.ListenOnNewState(func(S State) {
		value = S.State
		wgState.Done()
	})

	if returnValue := b.ListenOnNewState(func(S State) { value = S.State }); returnValue {
		t.Errorf("Expected value To Be Of Value '%t' Actual '%t'", false, returnValue)
	}

	if value != 0 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 0, value)
	}

	wgState.Add(1)
	wgEvent.Add(1)
	b.AddEvent(Event{Data: 1})
	wgEvent.Wait()
	wgState.Wait()
	if value != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, value)
	}

	b.StopListenToStateStream()
	wgEvent.Add(1)
	b.AddEvent(Event{Data: 2})
	wgEvent.Wait()
	if value != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, value)
	}
}

func TestBloC_AddEvent(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	b := CreateBloC[Event, State, BD](bd, func(E event.Event[Event], BD *BD) State { defer wg.Done(); value += E.Data.Data; return State{} })
	b.StartListenToEventStream()
	wg.Add(1)
	b.AddEvent(Event{Data: 1})
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}
}
