package stream_test

import (
	"dev.go-bloc/util/stream"
	"reflect"
	"testing"
	"time"
)

func TestCreateStream(t *testing.T) {
	s := stream.CreateStream[struct{}](10, func(NewItem struct{}) {})

	if value := s.MaxQueueSize; value != 10 {
		t.Errorf("Expected Field MaxQueueSize To Equal '%d' Actual '%d'", 10, value)
	}

	if reflect.TypeOf(s.OnNewItem) != reflect.TypeOf(func(struct{}) {}) {
		t.Errorf("Expected OnNewItem Of Type '%s' Actual '%s'", reflect.TypeOf(func(struct{}) {}), reflect.TypeOf(s.OnNewItem))
	}
}

func TestStream_Add(t *testing.T) {
	s := stream.CreateStream[int](2, func(int) {})

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
	s := stream.CreateStream[int](2, func(int) {})

	if value := s.GetQueueSize(); value != 0 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 0, value)
	}
	s.Add(1)
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
}

func TestStream_Listen(t *testing.T) {
	check := 0
	s := stream.CreateStream[int](1, func(NewItem int) { check += NewItem })
	s.Add(1)
	if check != 0 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 0, check)
	}
	s.Listen()
	s.Add(1)

	if check != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, check)
	}
}

func TestStream_StopListen(t *testing.T) {
	check := 0
	s := stream.CreateStream[int](2, func(NewItem int) { check += NewItem })
	s.Listen()
	s.Add(1)
	if check != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, check)
	}
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
	s.StopListen()
	s.Add(1)
	if value := s.GetQueueSize(); value != 2 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 2, value)
	}
	if check != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, check)
	}

}

func TestStream_ResumeAtQueuePosition(t *testing.T) {
	check := 0
	s := stream.CreateStream[int](2, func(NewItem int) { check += NewItem })
	s.Listen()
	s.Add(1)
	if check != 1 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 1, check)
	}
	if value := s.GetQueueSize(); value != 1 {
		t.Errorf("Expected GetQueueSize To Equal '%d' Actual '%d'", 1, value)
	}
	s.Add(2)
	s.ResumeAtQueuePosition(1)
	if check != 5 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 5, check)
	}
	s.Add(2)
	time.Sleep(1 * time.Second)
	if check != 7 {
		t.Errorf("Expected check To Equal '%d' Actual '%d'", 7, check)
	}
}
