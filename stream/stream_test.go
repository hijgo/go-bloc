package stream

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCreateStream(t *testing.T) {
	s := CreateStream(10, func(NewItem struct{}) {})

	if value := s.MaxHistorySize; value != 10 {
		t.Errorf("Expected Field MaxHistorySize To Equal '%d' Actual '%d'", 10, value)
	}

	if value := reflect.TypeOf(s.OnNewItem); value != reflect.TypeOf(func(struct{}) {}) {
		t.Errorf("Expected OnNewItem Of Type '%s' Actual '%s'", reflect.TypeOf(s.OnNewItem), value)
	}

	if value := reflect.TypeOf(s.sink); value != reflect.TypeOf(make(chan struct{})) {
		t.Errorf("Expected sink Of Type '%s' Actual '%s'", reflect.TypeOf(reflect.TypeOf(make(chan struct{}))), value)
	}

	if value := reflect.TypeOf(s.pauseListen); value != reflect.TypeOf(make(chan bool)) {
		t.Errorf("Expected pauseListen Of Type '%s' Actual '%s'", reflect.TypeOf(reflect.TypeOf(make(chan bool))), value)
	}

	if value := reflect.TypeOf(s.stopListen); value != reflect.TypeOf(make(chan struct{})) {
		t.Errorf("Expected stopListen Of Type '%s' Actual '%s'", reflect.TypeOf(reflect.TypeOf(make(chan struct{}))), value)
	}

	if value := reflect.TypeOf(s.history); value != reflect.TypeOf(make([]*struct{}, 0)) {
		t.Errorf("Expected history Of Type '%s' Actual '%s'", reflect.TypeOf(reflect.TypeOf(make([]*struct{}, 0))), value)
	}

	if value := len(s.history); value != 0 {
		t.Errorf("Expected len(history) Of Value '%d' Actual '%d'", 0, value)
	}

	if value := cap(s.history); value != 10 {
		t.Errorf("Expected len(history) Of Value '%d' Actual '%d'", 10, value)
	}
}

func TestStream_Add(t *testing.T) {
	var wg sync.WaitGroup
	s := CreateStream(2, func(int) { wg.Done() })
	val1, val2, val3 := 1, 2, 3

	s.Add(val1)
	if value := s.GetHistorySize(); value != 1 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 1, value)
	}

	s.Add(val2)
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	s.Add(val3)
	wg.Wait()
	if value := s.GetHistorySize(); value != 2 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 2, value)
	}

	expectedValues := []*int{&val2, &val3}
	for i, v := range s.history {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected history To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}

	defer func() {
		if err := s.Dispose(); err != nil {
			t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
		}
	}()
}

func TestStream_GetHistorySize(t *testing.T) {
	s := CreateStream(2, func(int) {})

	if value := s.GetHistorySize(); value != 0 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 0, value)
	}
	s.Add(1)
	if value := s.GetHistorySize(); value != 1 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 1, value)
	}
}

func TestStream_Listen(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream(1, func(NewItem int) { value += NewItem; wg.Done() })

	s.Add(1)
	if value != 0 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", 0, value)
	}

	wg.Add(1)
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	s.Add(1)
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", 1, value)
	}

	defer func() {
		if err := s.Dispose(); err != nil {
			t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
		}
	}()
}

func TestStream_StopListen(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream(2, func(NewItem int) { value += NewItem; wg.Done() })
	val1, val2 := 1, 2

	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	s.Add(val1)
	wg.Wait()
	if value != val1 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1, value)
	}

	if value := s.GetHistorySize(); value != 1 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 1, value)
	}

	err = s.StopListen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	s.Add(val2)
	if value := s.GetHistorySize(); value != 2 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 2, value)
	}

	if value != val1 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1, value)
	}

	expectedValues := []*int{&val1, &val2}
	for i, v := range s.history {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected history To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}

}

func TestStream_ResumeAtHistoryPosition(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream(2, func(NewItem int) { time.Sleep(50 * time.Millisecond); value += NewItem; wg.Done() })
	val1, val2 := 1, 2

	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	s.Add(val1)
	wg.Wait()
	if value != val1 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1, value)
	}
	if value := s.GetHistorySize(); value != 1 {
		t.Errorf("Expected GetHistorySize To Equal '%d' Actual '%d'", 1, value)
	}

	wg.Add(2)
	err = s.ResumeAtHistoryPosition(0)
	s.Add(val2)
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Wait()
	if value != val1*2+val2 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1*2+val2, value)
	}

	expectedValues := []*int{&val1, &val2}
	for i, v := range s.history {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected history To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}
	defer func() {
		if err := s.Dispose(); err != nil {
			t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
		}
	}()
}

func TestStream_GetListenStatus(t *testing.T) {
	s := CreateStream(2, func(NewItem int) {})
	if value := s.GetListenStatus(); value {
		t.Errorf("Expected GetListenStatus To Equal '%t' Actual '%t'", false, value)
	}
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}
	if value := s.GetListenStatus(); !value {
		t.Errorf("Expected GetListenStatus To Equal '%t' Actual '%t'", true, value)
	}
	defer func() {
		if err := s.Dispose(); err != nil {
			t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
		}
	}()
}

func TestStream_Dispose(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})

	if err := s.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
	}

	if value := s.wasDisposed; !value {
		t.Errorf("Expected wasDisposed To Equal '%t' Actual '%t'", true, value)
	}
}

func TestStream_DisposeWhileAlreadyListening(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})
	if err := s.Listen(); err != nil {
		t.Errorf("Unexpected Error: '%s' while trying listen Stream", err.Error())
	}
	if err := s.Dispose(); err != nil {
		t.Errorf("Unexpected Error: '%s' while disposing Stream", err.Error())
	}

	if value := s.wasDisposed; !value {
		t.Errorf("Expected wasDisposed To Equal '%t' Actual '%t'", true, value)
	}
}

func TestStream_ListenShouldReturnErrorWhenAlreadyDisposed(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})

	s.Dispose()

	err := s.Listen()
	if err == nil {
		t.Errorf("Expected Listen To Return Error When Already Disposed")
	} else if err != &AlreadyDisposedErr {
		t.Errorf("Expected Listen To Return Error With Message '%s' Actual '%s'", AlreadyDisposedErr.Error(), err.Error())
	}
}

func TestStream_ListenShouldReturnErrorWhenAlreadyListenedTo(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})

	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	err = s.Listen()
	if err == nil {
		t.Errorf("Expected Listen To Return Error When Already Listend To")
	} else if err != &StartListenErr {
		t.Errorf("Expected Listen To Return Error With Message '%s' Actual '%s'", StartListenErr.Error(), err.Error())
	}

	defer s.Dispose()
}

func TestStream_ResumeAtHistoryPositionShouldReturnErrorWhenPositionNotInRange(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	err = s.ResumeAtHistoryPosition(0)
	s.waitForResumeAtPositionCompletion.Wait()
	if err == nil {
		t.Errorf("Expected ResumeAtHistoryPosition To Return Error When Position Out Of Range")
	} else if wantedErr := PosOutOfRangeErr(0, 0); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s' Actual '%s'", wantedErr.Error(), err.Error())
	}

	err = s.ResumeAtHistoryPosition(-2)
	s.waitForResumeAtPositionCompletion.Wait()
	if err == nil {
		t.Errorf("Expected ResumeAtHistoryPosition To Return Error When Position Out Of Range")
	} else if wantedErr := PosOutOfRangeErr(-2, 0); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s' Actual '%s'", wantedErr.Error(), err.Error())
	}

	defer s.Dispose()
}

func TestStream_StopListenShouldReturnErrorWhenAlreadyStopped(t *testing.T) {
	s := CreateStream(1, func(NewItem int) {})

	err := s.StopListen()
	if err == nil {
		t.Errorf("Expected StopListen To Return Error When Not Listened To")
	} else if err != StopListenErr {
		t.Errorf("Expected StopListen To Return Error With Message '%s ' Actual '%s'", StopListenErr, err.Error())
	}

	err = s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	err = s.StopListen()
	if err != nil {
		t.Errorf("Expected StopListen To Return No Error But Returned '%s'", err.Error())
	}

	err = s.StopListen()
	if err == nil {
		t.Errorf("Expected StopListen To Return Error When Not Listened To")
	} else if err != StopListenErr {
		t.Errorf("Expected StopListen To Return Error With Message '%s ' Actual '%s'", StopListenErr.Error(), err.Error())
	}
}
