package queue

import (
	"sync"

	"github.com/slack-go/slack/slackevents"
)

type EventCallback func(e *slackevents.EventsAPIEvent)

type EventQueue struct {
	queue     chan *slackevents.EventsAPIEvent
	callbacks []EventCallback

	wg   *sync.WaitGroup
	stop chan interface{}
}

func NewEventQueue(size int64) *EventQueue {
	return &EventQueue{
		queue:     make(chan *slackevents.EventsAPIEvent, size),
		callbacks: []EventCallback{},
		wg:        &sync.WaitGroup{},
		stop:      make(chan interface{}),
	}
}

func (q *EventQueue) AddCallback(s EventCallback) {
	q.callbacks = append(q.callbacks, s)
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
	for _, c := range q.callbacks {
		q.wg.Add(1)

		go func(wg *sync.WaitGroup, callback EventCallback) {
			callback(e)
			wg.Done()
		}(q.wg, c)
	}
}
