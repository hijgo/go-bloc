package stream

import (
	"sync"
)

type Stream[T any] struct {
	MaxQueueSize int
	OnNewItem    func(NewItem T)
	sink         chan T
	isListenedTo bool
	pauseListen  chan bool
	stopListen   chan struct{}
	queue        []*T
}

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

func (s *Stream[_]) GetListenStatus() bool {
	return s.isListenedTo
}

func (s *Stream[_]) StopListen() {
	if !s.isListenedTo {
		return
	}
	s.isListenedTo = false
	s.stopListen <- struct{}{}
}

func (s *Stream[T]) Listen() {
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
						select {
						case isPaused := <-s.pauseListen:
							if !isPaused {
								wg.Done()
								return
							}
						}
					}()
					wg.Wait()
				}
			case <-s.stopListen:
				return
			}
		}
	}()
}

func (s Stream[_]) GetQueueSize() int {
	return len(s.queue)
}

func (s Stream[T]) ResumeAtQueuePosition(Position int) {
	s.pauseListen <- true
	s.OnNewItem(*s.queue[Position])

	defer func() {
		s.pauseListen <- false
		s.queue = s.queue[:Position+1]
	}()
}

func (s *Stream[T]) Add(NewItem T) {
	if len(s.queue) >= s.MaxQueueSize {
		s.queue = s.queue[1:s.MaxQueueSize]
		s.queue = append(s.queue, &NewItem)
	} else {
		s.queue = append(s.queue, &NewItem)
	}
	if s.isListenedTo {
		s.sink <- NewItem
	}
}
