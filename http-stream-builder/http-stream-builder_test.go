package httpstreambuilder

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/hijgo/go-bloc/bloc"
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

const (
	defaultPath = "/"
)

func resetEnv() {
	testBloC = bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	mux = http.NewServeMux()
	bd = BD{}
	value = 0
	streamBuilder = CreateHttpStreamBuilder(testBloC, func(NewState State) func(w http.ResponseWriter, r *http.Request) {
		value = NewState.State
		wg.Done()

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, NewState.State)
		}
	})
	initialEvent = Event{
		Data: 1,
	}
	if testServer != nil {
		testServer.Close()
		testServer = nil
	}
}

var (
	wg            sync.WaitGroup
	mux           = http.NewServeMux()
	testBloC      = bloc.CreateBloC(bd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	bd            = BD{}
	testServer    *httptest.Server
	value         int
	streamBuilder = CreateHttpStreamBuilder(testBloC, func(NewState State) func(http.ResponseWriter, *http.Request) {
		value = NewState.State
		wg.Done()

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, NewState.State)
		}
	})
	initialEvent = Event{
		Data: 1,
	}
)

func TestStreamBuilder_Init(t *testing.T) {
	wg.Add(1)
	streamBuilder.Init(defaultPath, nil, initialEvent)
	testServer = httptest.NewServer(http.DefaultServeMux)
	wg.Wait()

	if value := reflect.TypeOf(streamBuilder.BloC); value != reflect.TypeOf(testBloC) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(testBloC), value)
	}

	if value := reflect.TypeOf(streamBuilder.builderFunc); value != reflect.TypeOf(func(NewState State) func(http.ResponseWriter, *http.Request) {
		return func(http.ResponseWriter, *http.Request) {}
	}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), value)
	}
	if value != initialEvent.Data {
		t.Errorf("Expected initialEvent altering value to 1")
	}
	wg.Add(1)
	streamBuilder.BloC.AddEvent(Event{Data: 2})
	wg.Wait()
	if value != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, value)
	}

	req, err := http.Get(testServer.URL + defaultPath)
	if err != nil {
		t.Errorf("Error was returned while trying to connect with httpstreambuilder-endpoint. Err: %s", err.Error())
	}

	var bytes []byte
	if bytes, err = io.ReadAll(req.Body); err != nil {
		t.Errorf("Error was returned while trying to read the conent of httpstreambuilder-endpoint. Err: %s", err.Error())
	}

	if string(bytes) != fmt.Sprint(value) {
		t.Errorf("Expected httpstreambuilder-endpoint to mirror the internal state. Expected: %s Actual: %s", fmt.Sprint(value), string(bytes))
	}
	t.Cleanup(resetEnv)
}

func TestStreamBuilder_InitOnError(t *testing.T) {
	wg.Add(3)
	streamBuilder.Init(defaultPath, mux, initialEvent)

	if err := streamBuilder.Init(defaultPath+"1", mux, initialEvent); err == nil || err != &stream.StartListenErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	streamBuilder.Dispose()

	if err := streamBuilder.Init(defaultPath+"2", mux, initialEvent); err == nil || err != &stream.AlreadyDisposedErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.AlreadyDisposedErr), reflect.TypeOf(err))
	}
	t.Cleanup(resetEnv)
}

func TestStreamBuilder_Dispose(t *testing.T) {
	wg.Add(1)
	streamBuilder.Init(defaultPath, mux, initialEvent)
	streamBuilder.Dispose()

	expected := value
	streamBuilder.BloC.AddEvent(Event{Data: 2})

	if value != expected {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, value)
	}
	t.Cleanup(resetEnv)
}
