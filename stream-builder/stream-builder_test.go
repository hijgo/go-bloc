package stream_builder

import (
	"dev.hijgo.go-bloc/bloc"
	"dev.hijgo.go-bloc/event"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Event struct {
	Data int
}

type AD struct {
	AD string
}

type State struct {
	State int
}

func TestInitStreamBuilder(t *testing.T) {
	ad := AD{}
	check := 0
	initialEvent := Event{
		Data: 1,
	}
	b := bloc.CreateBloC[Event, State, AD](ad, func(E event.Event[Event], AD AD) State { return State{State: E.Data.Data} })
	streamBuilder := InitStreamBuilder[Event, State, AD](b, initialEvent, func(NewState State) {
		check = NewState.State
	})

	if value := reflect.TypeOf(streamBuilder.BloC); value != reflect.TypeOf(b) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(b), value)
	}

	if value := reflect.TypeOf(streamBuilder.BuilderFunc); value != reflect.TypeOf(func(State) {}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), value)
	}

	if value := reflect.TypeOf(streamBuilder.InitialEvent); value != reflect.TypeOf(Event{}) {
		t.Errorf("Expected InitialEvent Of Type '%s' Actual '%s'", reflect.TypeOf(Event{}), value)
	}

	if value := streamBuilder.InitialEvent; value != initialEvent {
		t.Errorf("Expected BloC Of Value '%s' Actual '%s'", fmt.Sprint(initialEvent), fmt.Sprint(value))

	}
	time.Sleep(5 * time.Microsecond)
	if check != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, check)
	}
	streamBuilder.BloC.AddEvent(Event{
		Data: 2,
	})
	time.Sleep(5 * time.Microsecond)
	if check != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, check)
	}
}

func TestStreamBuilder_Dispose(t *testing.T) {
	ad := AD{}
	check := 0
	initialEvent := Event{
		Data: 1,
	}
	b := bloc.CreateBloC[Event, State, AD](ad, func(E event.Event[Event], AD AD) State { return State{State: E.Data.Data} })
	streamBuilder := InitStreamBuilder[Event, State, AD](b, initialEvent, func(NewState State) {
		check = NewState.State
	})
	time.Sleep(5 * time.Microsecond)
	if check != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, check)
	}
	streamBuilder.Dispose()
	streamBuilder.BloC.AddEvent(Event{Data: 2})
	time.Sleep(500 * time.Millisecond)
	if check != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, check)
	}
}
