package queue

import (
	"sync"

	"github.com/slack-go/slack"
)

type CommandCallback func(e *slack.SlashCommand)

type CommandQueue struct {
	queue     chan *slack.SlashCommand
	callbacks []CommandCallback

	wg   *sync.WaitGroup
	stop chan interface{}
}

func NewCommandQueue(size int64) *CommandQueue {
	return &CommandQueue{
		queue:     make(chan *slack.SlashCommand, size),
		callbacks: []CommandCallback{},
		wg:        &sync.WaitGroup{},
		stop:      make(chan interface{}),
	}
}

func (q *CommandQueue) AddCallback(s CommandCallback) {
	q.callbacks = append(q.callbacks, s)
}

func (q *CommandQueue) Enqueue(e *slack.SlashCommand) {
	q.queue <- e
}

func (q *CommandQueue) Start() {
	for {
		select {
		case e := <-q.queue:
			q.dispatch(e)
		case <-q.stop:
			return
		}
	}
}

func (q *CommandQueue) Stop() {
	q.stop <- true
}

func (q *CommandQueue) Wait() {
	q.wg.Wait()
}

func (q *CommandQueue) dispatch(e *slack.SlashCommand) {
	for _, c := range q.callbacks {
		q.wg.Add(1)

		go func(wg *sync.WaitGroup, callback CommandCallback) {
			callback(e)
			wg.Done()
		}(q.wg, c)
	}
}
