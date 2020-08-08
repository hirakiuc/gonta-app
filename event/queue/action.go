package queue

import (
	"fmt"
	"sync"

	"github.com/slack-go/slack"
)

type ActionCallback func(e *slack.InteractionCallback) error

type ActionQueue struct {
	queue     chan *slack.InteractionCallback
	callbacks map[string][]ActionCallback

	wg   *sync.WaitGroup
	stop chan interface{}
}

func NewActionQueue(size int64) *ActionQueue {
	return &ActionQueue{
		queue:     make(chan *slack.InteractionCallback, size),
		callbacks: map[string][]ActionCallback{},
		wg:        &sync.WaitGroup{},
		stop:      make(chan interface{}),
	}
}

func (q *ActionQueue) eventKey(eventType slack.InteractionType, blockID string) string {
	return fmt.Sprintf("%s:block-%s", eventType, blockID)
}

func (q *ActionQueue) AddBlockActionCallback(eventType slack.InteractionType, blockID string, c ActionCallback) {
	key := q.eventKey(eventType, blockID)

	v, ok := q.callbacks[key]
	if !ok {
		v = []ActionCallback{c}
	} else {
		v = append(v, c)
	}

	q.callbacks[key] = v
}

func (q *ActionQueue) Enqueue(e *slack.InteractionCallback) {
	q.queue <- e
}

func (q *ActionQueue) Start() {
	for {
		select {
		case e := <-q.queue:
			// nolint:exhaustive
			switch e.Type {
			case slack.InteractionTypeBlockActions:
				q.dispatchBlockAction(e)
			default:
				return
			}
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

func (q *ActionQueue) dispatchBlockAction(e *slack.InteractionCallback) {
	action := e.ActionCallback.BlockActions[0]
	key := q.eventKey(e.Type, action.BlockID)

	callbacks, ok := q.callbacks[key]
	if !ok {
		return
	}

	q.dispatch(e, callbacks)
}

func (q *ActionQueue) dispatch(e *slack.InteractionCallback, callbacks []ActionCallback) {
	if len(callbacks) == 0 {
		return
	}

	for _, c := range callbacks {
		q.wg.Add(1)

		go func(wg *sync.WaitGroup, callback ActionCallback) {
			_ = callback(e)

			wg.Done()
		}(q.wg, c)
	}
}
