package event

import "testing"

func TestCreateEvent(t *testing.T) {
	event := CreateEvent[int](1)
	if value := event.Data; value != 1 {
		t.Errorf("Expected Field Data To Equal '%d' Actual '%d'", 1, value)
	}
}
