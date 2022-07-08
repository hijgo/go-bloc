package stream_builder

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hijgo/go-bloc/bloc"
	"github.com/hijgo/go-bloc/event"
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

func TestInitStreamBuilder(t *testing.T) {
	var wg sync.WaitGroup
	bd := BD{}
	value := 0
	initialEvent := Event{
		Data: 1,
	}

	wg.Add(1)

	b := bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	streamBuilder, err := InitStreamBuilder(b, &initialEvent, func(NewState State) {
		value = NewState.State
		wg.Done()
	})
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	if value := reflect.TypeOf(streamBuilder.BloC); value != reflect.TypeOf(b) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(b), value)
	}

	if value := reflect.TypeOf(streamBuilder.builderFunc); value != reflect.TypeOf(func(State) {}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), value)
	}

	if value := reflect.TypeOf(*streamBuilder.initialEvent); value != reflect.TypeOf(Event{}) {
		t.Errorf("Expected InitialEvent Of Type '%s' Actual '%s'", reflect.TypeOf(Event{}), value)
	}

	if value := *streamBuilder.initialEvent; value != initialEvent {
		t.Errorf("Expected initialEvent Of Value '%s' Actual '%s'", fmt.Sprint(initialEvent), fmt.Sprint(value))

	}
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, value)
	}

	wg.Add(1)
	streamBuilder.BloC.AddEvent(Event{
		Data: 2,
	})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, value)
	}
}

func TestStreamBuilder_Dispose(t *testing.T) {
	var wgBuild sync.WaitGroup
	bd := BD{}
	value := 0
	initialEvent := Event{
		Data: 1,
	}

	wgBuild.Add(1)

	b := bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	streamBuilder, err := InitStreamBuilder(b, &initialEvent, func(NewState State) {
		value = NewState.State
		wgBuild.Done()
	})
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wgBuild.Wait()
	if value != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, value)
	}

	streamBuilder.Dispose()
	streamBuilder.BloC.AddEvent(Event{Data: 2})
	time.Sleep(1 * time.Second)
	if value != 1 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, value)
	}
}
