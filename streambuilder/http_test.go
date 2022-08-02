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

func httpTestResetEnv() {
	httpTestBloC = bloc.CreateBloC(httpTestBd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	httpTestMux = http.NewServeMux()
	httpTestBd = BD{}
	httpTestValue = 0
	httpTestStreamBuilder = CreateHttp(httpTestBloC, func(NewState State, out **int) {
		httpTestValue = NewState.State
		httpTestWg.Done()
		(*out) = &NewState.State
	}, httpTestMux, defaultPath)
	httpTestInitialEvent = Event{
		Data: 1,
	}
	if httpTestServer != nil {
		httpTestServer.Close()
		httpTestServer = nil
	}
}

var (
	httpTestWg            sync.WaitGroup
	httpTestMux           = http.NewServeMux()
	httpTestBloC          = bloc.CreateBloC(httpTestBd, func(E event.Event[Event], BD *BD) State { return State{State: E.Data.Data} })
	httpTestBd            = BD{}
	httpTestServer        *httptest.Server
	httpTestValue         int
	httpTestStreamBuilder = CreateHttp(httpTestBloC, func(NewState State, out **int) {
		httpTestValue = NewState.State
		httpTestWg.Done()
		(*out) = &NewState.State
	}, httpTestMux, defaultPath)
	httpTestInitialEvent = Event{
		Data: 1,
	}
)

func TestHttp_ShouldImplementAdvancedHttp(t *testing.T) {
	if !reflect.TypeOf(httpTestStreamBuilder).Implements(reflect.TypeOf((*StreamBuilder[Event, State, BD])(nil)).Elem()) {
		t.Errorf("BasicHttp does not implement AdvancedHttp Interface")
	}
}

func TestHttp_Init(t *testing.T) {
	httpTestWg.Add(1)
	if err := httpTestStreamBuilder.Init(&httpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}

	httpTestServer = httptest.NewServer(httpTestMux)
	httpTestWg.Wait()
	if httpTestValue := reflect.TypeOf(httpTestStreamBuilder.BloC); httpTestValue != reflect.TypeOf(httpTestBloC) {
		t.Errorf("Expected BloC Of Type '%s' Actual '%s'", reflect.TypeOf(httpTestBloC), httpTestValue)
	}

	if httpTestValue := reflect.TypeOf(httpTestStreamBuilder.builderFunc); httpTestValue != reflect.TypeOf(func(NewState State, a **int) {}) {
		t.Errorf("Expected BuilderFunc Of Type '%s' Actual '%s'", reflect.TypeOf(func(State, **int) {}), httpTestValue)
	}

	if httpTestValue != httpTestInitialEvent.Data {
		t.Errorf("Expected httpTestInitialEvent altering httpTestValue to 1")
	}

	httpTestWg.Add(1)
	httpTestBloC.AddEvent(Event{Data: 2})
	httpTestWg.Wait()
	if httpTestValue != 2 {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 2, httpTestValue)
	}

	req, err := http.Get(httpTestServer.URL + defaultPath)
	if err != nil {
		t.Errorf("Error was returned while trying to connect with AdvancedHttp-endpoint. Err: %s", err.Error())
	}

	var bytes []byte
	if bytes, err = io.ReadAll(req.Body); err != nil {
		t.Errorf("Error was returned while trying to read the conent of AdvancedHttp-endpoint. Err: %s", err.Error())
	}

	if string(bytes) != fmt.Sprint(httpTestValue) {
		t.Errorf("Expected Http-endpoint to mirror the internal state. Expected: %s Actual: %s", fmt.Sprint(httpTestValue), string(bytes))
	}
	t.Cleanup(httpTestResetEnv)
}

func TestHttp_InitOnError(t *testing.T) {
	httpTestWg.Add(3)
	if err := httpTestStreamBuilder.Init(&httpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}
	httpTestStreamBuilder.pattern = defaultPath + "1"
	if err := httpTestStreamBuilder.Init(&httpTestInitialEvent); err == nil || err != &stream.StartListenErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.StartListenErr), reflect.TypeOf(err))
	}
	if err := httpTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing AdvancedHttp", err.Error())
	}
	httpTestStreamBuilder.pattern = defaultPath + "2"
	if err := httpTestStreamBuilder.Init(&httpTestInitialEvent); err == nil || err != &stream.AlreadyDisposedErr {
		t.Errorf("Expected StartListenToEventStream error of type '%s' actual was '%s'", reflect.TypeOf(stream.AlreadyDisposedErr), reflect.TypeOf(err))
	}
	t.Cleanup(httpTestResetEnv)
}

func TestHttp_Dispose(t *testing.T) {
	httpTestWg.Add(1)
	if err := httpTestStreamBuilder.Init(&httpTestInitialEvent); err != nil {
		t.Errorf("Unexpected Error: '%s' while initializing AdvancedHttp", err.Error())
	}
	if err := httpTestStreamBuilder.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing AdvancedHttp", err.Error())
	}
	expected := httpTestValue
	httpTestStreamBuilder.BloC.AddEvent(Event{Data: 2})

	if httpTestValue != expected {
		t.Errorf("Expected check Of Value '%d' Actual '%d'", 1, httpTestValue)
	}
	t.Cleanup(httpTestResetEnv)
}
