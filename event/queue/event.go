package queue

import (
	"sync"

	"github.com/slack-go/slack/slackevents"
)

type EventCallback func(e *slackevents.EventsAPIEvent) error

type EventQueue struct {
	queue     chan *slackevents.EventsAPIEvent
	callbacks map[string][]EventCallback

	wg   *sync.WaitGroup
	stop chan interface{}
}

func NewEventQueue(size int64) *EventQueue {
	return &EventQueue{
		queue:     make(chan *slackevents.EventsAPIEvent, size),
		callbacks: map[string][]EventCallback{},
		wg:        &sync.WaitGroup{},
		stop:      make(chan interface{}),
	}
}

func (q *EventQueue) AddCallback(eventType string, c EventCallback) {
	v, ok := q.callbacks[eventType]
	if !ok {
		v = []EventCallback{c}
	} else {
		v = append(v, c)
	}

	q.callbacks[eventType] = v
}

func (q *EventQueue) Enqueue(e *slackevents.EventsAPIEvent) {
	q.queue <- e
}

func (q *EventQueue) Start() {
	for {
		select {
		case e := <-q.queue:
			q.dispatch(e)
		case <-q.stop:
			return
		}
	}
}

func (q *EventQueue) Stop() {
	q.stop <- true
}

func (q *EventQueue) Wait() {
	q.wg.Wait()
}

func (q *EventQueue) dispatch(e *slackevents.EventsAPIEvent) {
	eventType := e.InnerEvent.Type

	callbacks, ok := q.callbacks[eventType]
	if !ok {
		return
	}

	for _, c := range callbacks {
		q.wg.Add(1)

		go func(wg *sync.WaitGroup, callback EventCallback) {
			_ = callback(e)

			wg.Done()
		}(q.wg, c)
	}
}
