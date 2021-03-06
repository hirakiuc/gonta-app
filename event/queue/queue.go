package queue

import (
	"sync"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

type Queue struct {
	events   *EventQueue
	actions  *ActionQueue
	commands *CommandQueue

	log *zap.Logger
}

func New(size int64, log *zap.Logger) *Queue {
	return &Queue{
		events:   NewEventQueue(size),
		actions:  NewActionQueue(size),
		commands: NewCommandQueue(size),
		log:      log,
	}
}

func (q *Queue) Start() {
	go q.events.Start()
	go q.actions.Start()
	go q.commands.Start()

	q.log.Info("Started queue...")
}

func (q *Queue) Stop() {
	q.events.Stop()
	q.actions.Stop()
	q.commands.Stop()

	q.log.Info("Sent stop signals to queues...")
}

func (q *Queue) AddEventCallback(eventType string, c EventCallback) {
	q.events.AddCallback(eventType, c)
}

func (q *Queue) AddBlockActionCallback(eventType slack.InteractionType, blockID string, c ActionCallback) {
	q.actions.AddBlockActionCallback(eventType, blockID, c)
}

func (q *Queue) AddCommandCallback(c CommandCallback) {
	q.commands.AddCallback(c)
}

func (q *Queue) EnqueueEvent(e *slackevents.EventsAPIEvent) {
	q.events.Enqueue(e)
}

func (q *Queue) EnqueueAction(e *slack.InteractionCallback) {
	q.actions.Enqueue(e)
}

func (q *Queue) EnqueueCommand(e *slack.SlashCommand) {
	q.commands.Enqueue(e)
}

func (q *Queue) WaitUntilFinish() {
	wg := &sync.WaitGroup{}

	// Wait events queue
	wg.Add(1)

	go func() {
		q.events.Wait()
		wg.Done()
	}()

	// Wait actions queue
	wg.Add(1)

	go func() {
		q.actions.Wait()
		wg.Done()
	}()

	// Wait commands queue
	wg.Add(1)

	go func() {
		q.commands.Wait()
		wg.Done()
	}()

	q.log.Info("Waiting for queue processes...")
	wg.Wait()
	q.log.Info("Waiting for queue processes -> Done")
}
