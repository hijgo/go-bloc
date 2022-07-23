package basicstreambuilder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/hijgo/go-bloc/bloc"
	"github.com/hijgo/go-bloc/event"
	"github.com/hijgo/go-bloc/stream"
	stream_builder "github.com/hijgo/go-bloc/stream-builder"
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

func resetEnv() {
	testBloC = bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	bd = BD{}
	value = 0
	streamBuilder = CreateBasicStreamBuilder(testBloC, func(NewState State) {
		value = NewState.State
		wg.Done()
	})
	initialEvent = Event{
		Data: 1,
	}
}

var (
	wg            sync.WaitGroup
	testBloC      = bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	bd            = BD{}
	value         int
	streamBuilder = CreateBasicStreamBuilder(testBloC, func(NewState State) {
		value = NewState.State
		wg.Done()
	})
	initialEvent = Event{
		Data: 1,
	}
)

func TestBasicStreamBuilder_ShouldImplementStreamBuilder(t *testing.T) {
	if !reflect.TypeOf(streamBuilder).Implements(reflect.TypeOf((*stream_builder.StreamBuilder[Event, State, BD])(nil)).Elem()) {
		t.Errorf("BasicStreamBuilder does not implement StreamBuilder Interface")
	}
}

func TestStreamBuilder_Init(t *testing.T) {
	wg.Add(1)
	streamBuilder.Init(&initialEvent)

	if value := reflect.TypeOf(streamBuilder.BloC); value != reflect.TypeOf(testBloC) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(testBloC), value)
	}

	if value := reflect.TypeOf(streamBuilder.builderFunc); value != reflect.TypeOf(func(State) {}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), value)
	}
	wg.Wait()
	if value != initialEvent.Data {
		t.Errorf("Expected initialEvent altering value to 1")
	}
	wg.Add(1)
	streamBuilder.BloC.AddEvent(Event{Data: 2})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, value)

	}
	t.Cleanup(resetEnv)
}

func TestStreamBuilder_InitOnError(t *testing.T) {
	wg.Add(3)
	streamBuilder.Init(&initialEvent)

	if err := streamBuilder.Init(&initialEvent); err == nil || err != &stream.StartListenErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	streamBuilder.Dispose()
	if err := streamBuilder.Init(&initialEvent); err == nil || err != &stream.AlreadyDisposedErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	t.Cleanup(resetEnv)

}

func TestStreamBuilder_Dispose(t *testing.T) {
	wg.Add(1)
	streamBuilder.Init(&initialEvent)
	streamBuilder.Dispose()

	expected := value
	streamBuilder.BloC.AddEvent(Event{Data: 2})
	if value != expected {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, value)
	}
	t.Cleanup(resetEnv)
}
