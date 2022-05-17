package stream

import (
	"reflect"
	"sync"
	"testing"
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
	s := CreateStream[int](2, func(int) {})

	s.Add(1)
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
	s.Add(1)
	s.Add(1)
	if value := s.GetQueueSize(); value != 2 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 2, value)
	}
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
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 0, value)
	}

	wg.Add(1)
	s.Listen()
	s.Add(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, value)
	}
}

func TestStream_StopListen(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream[int](2, func(NewItem int) { value += NewItem; wg.Done() })

	s.Listen()
	wg.Add(1)
	s.Add(1)
	wg.Wait()

	if value != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, value)
	}
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
	s.StopListen()
	s.Add(1)
	if value := s.GetQueueSize(); value != 2 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 2, value)
	}
	if value != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, value)
	}

}

func TestStream_ResumeAtQueuePosition(t *testing.T) {
	var wg sync.WaitGroup
	value := 0
	s := CreateStream[int](2, func(NewItem int) { value += NewItem; wg.Done() })

	s.Listen()
	wg.Add(1)
	s.Add(1)
	wg.Wait()
	if value != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, value)
	}
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}

	wg.Add(2)
	s.Add(2)
	s.ResumeAtQueuePosition(1)
	wg.Wait()
	if value != 5 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 5, value)
	}

	wg.Add(1)
	s.Add(2)
	wg.Wait()
	if value != 7 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 7, value)
	}
}

func TestStream_GetListenStatus(t *testing.T) {
	s := CreateStream[int](2, func(NewItem int) {})
	if value := s.GetListenStatus(); value {
		t.Errorf("Expected GetListenStatus To Equal '%t' Actual '%t'", false, value)
	}
	s.Listen()
	if value := s.GetListenStatus(); !value {
		t.Errorf("Expected GetListenStatus To Equal '%t' Actual '%t'", true, value)
	}

}
