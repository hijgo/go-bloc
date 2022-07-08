package bloc

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
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
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

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
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		value += NewItem.Data.Data
		wg.Done()
	}

	b.AddEvent(Event{Data: 1})

	if value != 0 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 0, value)
	}

	if err := b.StartListenToEventStream(); err != nil {
		t.Errorf("Expected StartListenToEventStream To Return No Error Actual '%s'", err.Error())
	}

	wg.Add(1)
	b.AddEvent(Event{Data: 2})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected value To Be Of Value '%d' Actual '%d'", 2, value)
	}

	defer b.Dispose()
}

func TestBloC_StartListenToEventStreamShouldReturnErrorWhenAlreadyListenedTo(t *testing.T) {
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

	err := b.eventStream.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}
	err = b.StartListenToEventStream()
	if err == nil {
		t.Errorf("Expected StartListenToEventStream To Return Error When Already Listened To")
	} else if wantedErr := errors.New("stream already listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected StartListenToEventStream To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}
	defer b.Dispose()
}

func TestBloC_StopListenToEventStream(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

	b.eventStream.OnNewItem = func(NewItem event.Event[Event]) {
		value += NewItem.Data.Data
		wg.Done()
	}

	err := b.StartListenToEventStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	b.AddEvent(Event{Data: 1})
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}

	err = b.StopListenToEventStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	b.AddEvent(Event{Data: 2})
	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}
	defer b.Dispose()
}

func TestBloC_StopListenToEventStreamShouldReturnErrorWhenNotListenedTo(t *testing.T) {
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

	err := b.StopListenToEventStream()
	if err == nil {
		t.Errorf("Expected StopListenToEventStream To Return Error When Already Listened To")
	} else if wantedErr := errors.New("stream isn't listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected StopListenToEventStream To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}
	defer b.Dispose()

}

func TestBloC_ListenOnNewState(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State {
		return State{State: 2 * E.Data.Data}
	})
	value := 0

	err := b.StartListenToEventStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}
	b.AddEvent(Event{Data: 1})

	err = b.ListenOnNewState(func(S State) {
		value = S.State
		wg.Done()
	})
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
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
	defer b.Dispose()
}

func TestBloC_ListenOnNewStateShouldReturnErrorWhenAlreadyListenedTo(t *testing.T) {
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State {
		return State{State: 2 * E.Data.Data}
	})

	err := b.ListenOnNewState(func(state State) {})
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	err = b.ListenOnNewState(func(state State) {})
	if err == nil {
		t.Errorf("Expected ListenOnNewState To Return Error When Already Listened To")
	} else if wantedErr := errors.New("stream already listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected ListenOnNewState To Return Error With Message '%s' Actual '%s'", wantedErr.Error(), err.Error())
	}
	defer b.Dispose()
}

func TestBloC_StopListenToStateStream(t *testing.T) {
	var wgState sync.WaitGroup
	var wgEvent sync.WaitGroup
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { defer wgEvent.Done(); return State{State: 2 * E.Data.Data} })
	value := 0

	err := b.StartListenToEventStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wgEvent.Add(1)
	b.AddEvent(Event{Data: 1})
	wgEvent.Wait()

	err = b.ListenOnNewState(func(S State) {
		value = S.State
		wgState.Done()
	})
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
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

	err = b.StopListenToStateStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}
	wgEvent.Add(1)
	b.AddEvent(Event{Data: 2})
	wgEvent.Wait()
	if value != 2 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 2, value)
	}
	defer b.Dispose()
}

func TestBloC_StopListenToStateStreamShouldReturnErrorWhenNotListenedTo(t *testing.T) {
	bd := BD{}
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{} })

	err := b.StopListenToStateStream()
	if err == nil {
		t.Errorf("Expected StopListenToStateStream To Return Error When Not Listened To")
	} else if wantedErr := errors.New("stream isn't listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected StopListenToStateStream To Return Error With Message '%s' Actual '%s'", wantedErr.Error(), err.Error())
	}
	defer b.Dispose()
}

func TestBloC_AddEvent(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	b := CreateBloC(bd, func(E event.Event[Event], BD *BD) State { defer wg.Done(); value += E.Data.Data; return State{} })

	err := b.StartListenToEventStream()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	b.AddEvent(Event{Data: 1})
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check To Be Of Value '%d' Actual '%d'", 1, value)
	}
}
