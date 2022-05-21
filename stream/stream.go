package stream

import (
	"fmt"
	err "github.com/hijgo/go-bloc/error"
	"sync"
)

// Stream
// A structure that defines the operating values of a stream.
// Such as type of data being processed, behaviour when a new item is being passed down the stream
// and the size of the queue used as history.
//
// T - Type of the data that will be processed
// MaxQueueSize - The capacity of the history being saved
// OnNewItem - A Function that will be called everytime a new item is being passed to the stream
type Stream[T any] struct {
	MaxQueueSize                      int
	OnNewItem                         func(NewItem T)
	sink                              chan T
	isListenedTo                      bool
	pauseListen                       chan bool
	stopListen                        chan struct{}
	queue                             []*T
	wasDisposed                       bool
	waitForResumeAtPositionCompletion sync.WaitGroup
}

// CreateStream
// Function that should be called if a new stream is needed.
// Will create all necessary values so the stream can function properly and then return the new Stream of type T.
//
// T - Type of the data that will be processed
// MaxQueueSize - The capacity of the history being saved
// OnNewItem - A Function that will be called everytime a new item is being passed to the stream
func CreateStream[T any](MaxQueueSize int, OnNewItem func(NewItem T)) Stream[T] {
	return Stream[T]{
		MaxQueueSize: MaxQueueSize,
		OnNewItem:    OnNewItem,
		sink:         make(chan T),
		pauseListen:  make(chan bool),
		stopListen:   make(chan struct{}),
		queue:        make([]*T, 0, MaxQueueSize),
	}
}

// GetListenStatus
// Returns true if the stream is currently listened to, if not returns false.
func (s *Stream[_]) GetListenStatus() bool {
	return s.isListenedTo
}

// StopListen
// Will stop listening to the stream of incoming items. Any new item being passed into the stream will not be processed
// by the OnNewItem function, but will be stored in the queue.
// If the stream is not listened to, will return an error.
func (s *Stream[_]) StopListen() error {

	if !s.isListenedTo {
		defer func() {}()
		return &err.Error{
			Context: "Cannot call stop listening when stream isn't listened to!",
			Err:     fmt.Errorf("stream isn't listened to"),
		}
	}
	s.isListenedTo = false
	defer func() { s.stopListen <- struct{}{} }()

	return nil
}

// Listen
// Called to start listening to a stream of items. If the stream is already listened to, will return an error.
// Else will set the listening status to true and start processing new items with the OnNewItem function.
func (s *Stream[T]) Listen() error {
	if s.isListenedTo {
		return &err.Error{
			Context: "Cannot listen to stream, already being listened to!",
			Err:     fmt.Errorf("stream already listened to"),
		}
	} else if s.wasDisposed {
		return &err.Error{
			Context: "Cannot listen to stream, stream was disposed!",
			Err:     fmt.Errorf("stream was disposed"),
		}
	}

	s.isListenedTo = true
	go func() {
		for {
			select {
			case newItem := <-s.sink:
				s.OnNewItem(newItem)
			case isPaused := <-s.pauseListen:
				if isPaused {
					var wg sync.WaitGroup
					wg.Add(1)
					go func() {
						isPaused := <-s.pauseListen
						if !isPaused {
							wg.Done()
							return
						}
					}()
					wg.Wait()
				}
			case <-s.stopListen:
				return
			}
		}
	}()
	return nil
}

// GetQueueSize
// Returning the current length of the queue.
func (s *Stream[_]) GetQueueSize() int {
	return len(s.queue)
}

// ResumeAtQueuePosition
// Will temporally pause listening to stream to allow going back to a previous event. When paused new items will be
// processed when listening is resumed. All items in the queue before the given position will be dropped.
// Use with caution.
//
// Position - The position in the queue from where the queue should be resumed
// Will return an error when the given position is not inside the queue range.
func (s *Stream[T]) ResumeAtQueuePosition(Position int) error {
	s.waitForResumeAtPositionCompletion.Wait()
	s.waitForResumeAtPositionCompletion.Add(1)
	s.pauseListen <- true

	if queueLength := len(s.queue); Position < 0 || Position > queueLength || 0 == queueLength {
		defer func() {
			s.pauseListen <- false
			s.waitForResumeAtPositionCompletion.Done()
		}()
		return &err.Error{
			Context: "Wanted Position not in range of queue",
			Err:     fmt.Errorf("position '%d' out of range '%d'", Position, len(s.queue)),
		}
	}

	s.OnNewItem(*s.queue[Position])

	defer func() {
		s.pauseListen <- false
		s.queue = s.queue[:Position+1]
		s.waitForResumeAtPositionCompletion.Done()
	}()
	return nil
}

// Add
// Pass a NewItem into the stream
// The NewItem will only be processed if the stream is currently listened to.
// New Item will always be added to the queue.
func (s *Stream[T]) Add(NewItem T) {
	if s.isListenedTo {
		s.sink <- NewItem
	}
	if len(s.queue) >= s.MaxQueueSize {
		s.queue = s.queue[1:s.MaxQueueSize]
		s.queue = append(s.queue, &NewItem)
	} else {
		s.queue = append(s.queue, &NewItem)
	}
}

// Dispose
// Will close all channels used by the stream, in addition to that will also stop listening to stream.
// After disposing the stream cannot be listened to ever again.
func (s *Stream[_]) Dispose() {
	if s.GetListenStatus() {
		stopError := s.StopListen()
		if stopError != nil {
			panic(stopError)
		}
	}

	defer func() {
		close(s.sink)
		close(s.pauseListen)
		close(s.stopListen)
		s.wasDisposed = true
	}()
}
