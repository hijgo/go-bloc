package streambuilder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/hijgo/go-bloc/bloc"
	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
)

func resetEnv() {
	basicTestBloC = bloc.CreateBloC(basicTestBd, func(E event.Event[Event], bd *BD) State { return State{State: E.Data.Data} })
	basicTestBd = BD{}
	basicTestValue = 0
	basicTestStreamBuilder = CreateBasicStreamBuilder(basicTestBloC, func(NewState State) {
		basicTestValue = NewState.State
		basicTestWg.Done()
	})
	basicTestInitialEvent = Event{
		Data: 1,
	}
}

var (
	basicTestWg            sync.WaitGroup
	basicTestBloC          = bloc.CreateBloC(basicTestBd, func(E event.Event[Event], bd *BD) State { return State{State: E.Data.Data} })
	basicTestBd            = BD{}
	basicTestValue         int
	basicTestStreamBuilder = CreateBasicStreamBuilder(basicTestBloC, func(NewState State) {
		basicTestValue = NewState.State
		basicTestWg.Done()
	})
	basicTestInitialEvent = Event{
		Data: 1,
	}
)

func TestBasicStreamBuilder_ShouldImplementStreamBuilder(t *testing.T) {
	if !reflect.TypeOf(basicTestStreamBuilder).Implements(reflect.TypeOf((*StreamBuilder[Event, State, BD])(nil)).Elem()) {
		t.Errorf("BasicStreamBuilder does not implement StreamBuilder Interface")
	}
}

func TestStreamBuilder_Init(t *testing.T) {
	basicTestWg.Add(1)
	if err := basicTestStreamBuilder.Init(&basicTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing StreamBuilder", err.Error())
	}

	if value := reflect.TypeOf(basicTestStreamBuilder.BloC); value != reflect.TypeOf(basicTestBloC) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(basicTestBloC), value)
	}

	if value := reflect.TypeOf(basicTestStreamBuilder.builderFunc); value != reflect.TypeOf(func(State) {}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), value)
	}
	basicTestWg.Wait()
	if basicTestValue != basicTestInitialEvent.Data {
		t.Errorf("Expected initialEvent altering value to 1")
	}
	basicTestWg.Add(1)
	basicTestStreamBuilder.BloC.AddEvent(Event{Data: 2})
	basicTestWg.Wait()
	if basicTestValue != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, basicTestValue)

	}
	t.Cleanup(resetEnv)
}

func TestStreamBuilder_InitOnError(t *testing.T) {
	basicTestWg.Add(3)
	if err := basicTestStreamBuilder.Init(&basicTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing StreamBuilder", err.Error())
	}

	if err := basicTestStreamBuilder.Init(&basicTestInitialEvent); err == nil || err != &stream.StartListenErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	if err := basicTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing StreamBuilder", err.Error())
	}

	if err := basicTestStreamBuilder.Init(&basicTestInitialEvent); err == nil || err != &stream.AlreadyDisposedErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	t.Cleanup(resetEnv)
}

func TestStreamBuilder_Dispose(t *testing.T) {
	basicTestWg.Add(1)
	if err := basicTestStreamBuilder.Init(&basicTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing StreamBuilder", err.Error())
	}
	if err := basicTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing StreamBuilder", err.Error())
	}

	expected := basicTestValue
	basicTestStreamBuilder.BloC.AddEvent(Event{Data: 2})
	if basicTestValue != expected {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, basicTestValue)
	}
	t.Cleanup(resetEnv)
}
