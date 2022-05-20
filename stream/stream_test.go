package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCreateStream(t *testing.T) {
	s := CreateStream[struct{}](10, func(NewItem struct{}) {})

	if value := s.MaxQueueSize; value != 10 {
		t.Errorf("Expected Field MaxQueueSize To Equal '%d' Actual '%d'", 10, value)
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

	if value := reflect.TypeOf(s.queue); value != reflect.TypeOf(make([]*struct{}, 0)) {
		t.Errorf("Expected queue Of Type '%s' Actual '%s'", reflect.TypeOf(reflect.TypeOf(make([]*struct{}, 0))), value)
	}

	if value := len(s.queue); value != 0 {
		t.Errorf("Expected len(queue) Of Value '%d' Actual '%d'", 0, value)
	}

	if value := cap(s.queue); value != 10 {
		t.Errorf("Expected len(queue) Of Value '%d' Actual '%d'", 10, value)
	}

}

func TestStream_Add(t *testing.T) {
	var wg sync.WaitGroup
	s := CreateStream[int](2, func(int) { wg.Done() })
	val1, val2, val3 := 1, 2, 3

	s.Add(val1)
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}

	s.Add(val2)
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Add(1)
	s.Add(val3)
	wg.Wait()
	if value := s.GetQueueSize(); value != 2 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 2, value)
	}

	expectedValues := []*int{&val2, &val3}
	for i, v := range s.queue {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected queue To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}

	defer s.Dispose()
}

func TestStream_GetQueueSize(t *testing.T) {
	s := CreateStream[int](2, func(int) {})

	if value := s.GetQueueSize(); value != 0 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 0, value)
	}
	s.Add(1)
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
}

func TestStream_Listen(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream[int](1, func(NewItem int) { value += NewItem; wg.Done() })

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

	defer s.Dispose()
}

func TestStream_StopListen(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream[int](2, func(NewItem int) { value += NewItem; wg.Done() })
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

	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}

	err = s.StopListen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	s.Add(val2)
	if value := s.GetQueueSize(); value != 2 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 2, value)
	}

	if value != val1 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1, value)
	}

	expectedValues := []*int{&val1, &val2}
	for i, v := range s.queue {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected queue To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}

}

func TestStream_ResumeAtQueuePosition(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream[int](2, func(NewItem int) { time.Sleep(50 * time.Millisecond); value += NewItem; wg.Done() })
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
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}

	wg.Add(2)
	err = s.ResumeAtQueuePosition(0)
	s.Add(val2)
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	wg.Wait()
	if value != val1*2+val2 {
		t.Errorf("Expected value To Equal '%d' Actual '%d'", val1*2+val2, value)
	}

	expectedValues := []*int{&val1, &val2}
	for i, v := range s.queue {
		if value := *expectedValues[i]; value != *v {
			t.Errorf("Expected queue To Equal '%d' At Position '%d' Actual '%d'", value, i, *v)
		}
	}
	defer s.Dispose()

}

func TestStream_GetListenStatus(t *testing.T) {
	s := CreateStream[int](2, func(NewItem int) {})
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
	defer s.Dispose()
}

func TestStream_Dispose(t *testing.T) {
	s := CreateStream[int](1, func(NewItem int) {})

	s.Dispose()
	if value := s.wasDisposed; !value {
		t.Errorf("Expected wasDisposed To Equal '%t' Actual '%t'", false, value)
	}
}

func TestStream_ListenShouldReturnErrorWhenAlreadyDisposed(t *testing.T) {
	s := CreateStream[int](1, func(NewItem int) {})

	s.Dispose()

	err := s.Listen()
	if err == nil {
		t.Errorf("Expected Listen To Return Error When Already Disposed")
	} else if wantedErr := errors.New("stream was disposed"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}
}

func TestStream_ListenShouldReturnErrorWhenAlreadyListenedTo(t *testing.T) {
	s := CreateStream[int](1, func(NewItem int) {})

	err := s.Listen()
	err = s.Listen()
	if err == nil {
		t.Errorf("Expected Listen To Return Error When Already Listend To")
	} else if wantedErr := errors.New("stream already listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}

	defer s.Dispose()
}

func TestStream_ResumeAtQueuePositionShouldReturnErrorWhenPositionNotInRange(t *testing.T) {
	s := CreateStream[int](1, func(NewItem int) {})
	err := s.Listen()
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err.Error())
	}

	err = s.ResumeAtQueuePosition(0)
	s.waitForResumeAtPositionCompletion.Wait()
	if err == nil {
		t.Errorf("Expected ResumeAtQueuePosition To Return Error When Position Out Of Range")
	} else if wantedErr := errors.New(fmt.Sprintf("position '%d' out of range '%d'", 0, len(s.queue))); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}

	err = s.ResumeAtQueuePosition(-2)
	s.waitForResumeAtPositionCompletion.Wait()
	if err == nil {
		t.Errorf("Expected ResumeAtQueuePosition To Return Error When Position Out Of Range")
	} else if wantedErr := errors.New(fmt.Sprintf("position '%d' out of range '%d'", -2, len(s.queue))); err.Error() != wantedErr.Error() {
		t.Errorf("Expected Listen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}

	defer s.Dispose()
}

func TestStream_StopListenShouldReturnErrorWhenAlreadyStopped(t *testing.T) {
	s := CreateStream[int](1, func(NewItem int) {})

	err := s.StopListen()
	if err == nil {
		t.Errorf("Expected StopListen To Return Error When Not Listened To")
	} else if wantedErr := errors.New("stream isn't listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected StopListen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
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
	} else if wantedErr := errors.New("stream isn't listened to"); err.Error() != wantedErr.Error() {
		t.Errorf("Expected StopListen To Return Error With Message '%s ' Actual '%s'", wantedErr.Error(), err.Error())
	}
}
