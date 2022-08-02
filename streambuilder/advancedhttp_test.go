package streambuilder

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

func advHttpTestResetEnv() {
	advHttpTestBloC = bloc.CreateBloC(advHttpTestBd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	advHttpTestMux = http.NewServeMux()
	advHttpTestBd = BD{}
	advHttpTestValue = 0
	advHttpTestStreamBuilder = CreateAdvancedHttp(advHttpTestBloC, func(NewState State) func(w http.ResponseWriter, r *http.Request) {
		advHttpTestValue = NewState.State
		advHttpTestWg.Done()

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, NewState.State)
		}
	}, advHttpTestMux, defaultPath)
	advHttpTestInitialEvent = Event{
		Data: 1,
	}
	if testServer != nil {
		testServer.Close()
		testServer = nil
	}
}

var (
	advHttpTestWg            sync.WaitGroup
	advHttpTestMux           = http.NewServeMux()
	advHttpTestBloC          = bloc.CreateBloC(advHttpTestBd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	advHttpTestBd            = BD{}
	testServer               *httptest.Server
	advHttpTestValue         int
	advHttpTestStreamBuilder = CreateAdvancedHttp(advHttpTestBloC, func(NewState State) func(http.ResponseWriter, *http.Request) {
		advHttpTestValue = NewState.State
		advHttpTestWg.Done()

		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, NewState.State)
		}
	}, advHttpTestMux, defaultPath)
	advHttpTestInitialEvent = Event{
		Data: 1,
	}
)

func TestAdvancedHttp_ShouldImplementAdvancedHttp(t *testing.T) {
	if !reflect.TypeOf(advHttpTestStreamBuilder).Implements(reflect.TypeOf((*StreamBuilder[Event, State, BD])(nil)).Elem()) {
		t.Errorf("AdvancedHttp does not implement AdvancedHttp Interface")
	}
}

func TestAdvancedHttp_Init(t *testing.T) {
	advHttpTestWg.Add(1)
	if err := advHttpTestStreamBuilder.Init(&advHttpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}
	testServer = httptest.NewServer(advHttpTestMux)
	advHttpTestWg.Wait()

	if advHttpTestValue := reflect.TypeOf(advHttpTestStreamBuilder.BloC); advHttpTestValue != reflect.TypeOf(advHttpTestBloC) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(advHttpTestBloC), advHttpTestValue)
	}

	if advHttpTestValue := reflect.TypeOf(advHttpTestStreamBuilder.builderFunc); advHttpTestValue != reflect.TypeOf(func(NewState State) func(http.ResponseWriter, *http.Request) {
		return func(http.ResponseWriter, *http.Request) {}
	}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State) {}), advHttpTestValue)
	}
	if advHttpTestValue != advHttpTestInitialEvent.Data {
		t.Errorf("Expected advHttpTestInitialEvent altering advHttpTestValue to 1")
	}
	advHttpTestWg.Add(1)
	advHttpTestBloC.AddEvent(Event{Data: 2})
	advHttpTestWg.Wait()
	if advHttpTestValue != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, advHttpTestValue)
	}

	req, err := http.Get(testServer.URL + defaultPath)
	if err != nil {
		t.Errorf("Error was returned while trying to connect with AdvancedHttp-endpoint. Err: %s", err.Error())
	}

	var bytes []byte
	if bytes, err = io.ReadAll(req.Body); err != nil {
		t.Errorf("Error was returned while trying to read the conent of AdvancedHttp-endpoint. Err: %s", err.Error())
	}

	if string(bytes) != fmt.Sprint(advHttpTestValue) {
		t.Errorf("Expected AdvancedHttp-endpoint to mirror the internal state. Expected: %s Actual: %s", fmt.Sprint(advHttpTestValue), string(bytes))
	}
	t.Cleanup(advHttpTestResetEnv)
}

func TestAdvancedHttp_InitOnError(t *testing.T) {
	advHttpTestWg.Add(3)
	if err := advHttpTestStreamBuilder.Init(&advHttpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}
	advHttpTestStreamBuilder.pattern = defaultPath + "1"
	if err := advHttpTestStreamBuilder.Init(&advHttpTestInitialEvent); err == nil || err != &stream.StartListenErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	if err := advHttpTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing AdvancedHttp", err.Error())
	}
	advHttpTestStreamBuilder.pattern = defaultPath + "2"
	if err := advHttpTestStreamBuilder.Init(&advHttpTestInitialEvent); err == nil || err != &stream.AlreadyDisposedErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.AlreadyDisposedErr), reflect.TypeOf(err))
	}
	t.Cleanup(advHttpTestResetEnv)
}

func TestAdvancedHttp_Dispose(t *testing.T) {
	advHttpTestWg.Add(1)
	if err := advHttpTestStreamBuilder.Init(&advHttpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}
	if err := advHttpTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing AdvancedHttp", err.Error())
	}
	expected := advHttpTestValue
	advHttpTestStreamBuilder.BloC.AddEvent(Event{Data: 2})

	if advHttpTestValue != expected {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, advHttpTestValue)
	}
	t.Cleanup(advHttpTestResetEnv)
}
