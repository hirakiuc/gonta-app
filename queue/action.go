// nolint:dupl
package queue

import (
	"sync"

	"github.com/slack-go/slack"
)

type ActionCallback func(e *slack.InteractionCallback)

type ActionQueue struct {
	queue     chan *slack.InteractionCallback
	callbacks []ActionCallback

	wg   *sync.WaitGroup
	stop chan interface{}
}

func NewActionQueue(size int64) *ActionQueue {
	return &ActionQueue{
		queue:     make(chan *slack.InteractionCallback, size),
		callbacks: []ActionCallback{},
		wg:        &sync.WaitGroup{},
		stop:      make(chan interface{}),
	}
}

func (q *ActionQueue) AddCallback(s ActionCallback) {
	q.callbacks = append(q.callbacks, s)
}

func (q *ActionQueue) Enqueue(e *slack.InteractionCallback) {
	q.queue <- e
}

func (q *ActionQueue) Start() {
	for {
		select {
		case e := <-q.queue:
			q.dispatch(e)
		case <-q.stop:
			return
		}
	}
}

func (q *ActionQueue) Stop() {
	q.stop <- true
}

func (q *ActionQueue) Wait() {
	q.wg.Wait()
}

func (q *ActionQueue) dispatch(e *slack.InteractionCallback) {
	for _, c := range q.callbacks {
		q.wg.Add(1)

		go func(wg *sync.WaitGroup, callback ActionCallback) {
			callback(e)
			wg.Done()
		}(q.wg, c)
	}
}
