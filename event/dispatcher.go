package event

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/hirakiuc/gonta-app/handler"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

const JobTimeout = 3 * time.Second

type Dispatcher struct {
	log      *zap.Logger
	handlers map[string][]handler.Handler
}

func NewDispatcher(logger *zap.Logger) *Dispatcher {
	d := &Dispatcher{
		log:      logger,
		handlers: map[string][]handler.Handler{},
	}

	handlerMap := handler.GenerateHandlerMap()

	for eventType, handlers := range handlerMap {
		for _, h := range handlers {
			d.Register(eventType, h)
		}
	}

	return d
}

func (d *Dispatcher) Register(eventType string, h handler.Handler) {
	d.log.Debug("Register handler", zap.String("event", eventType), zap.String("handler", reflect.TypeOf(h).String()))

	handlers, ok := d.handlers[eventType]
	if ok {
		d.handlers[eventType] = append(handlers, h)
	} else {
		d.handlers[eventType] = []handler.Handler{h}
	}
}

// Dispatch invoke handlers with this event.
func (d *Dispatcher) Dispatch(ctx context.Context, e *slackevents.EventsAPIEvent) *sync.WaitGroup {
	log := d.log
	wg := &sync.WaitGroup{}
	innerEvent := e.InnerEvent

	handlers, ok := d.handlers[innerEvent.Type]
	if !ok {
		return wg
	}

	for _, h := range handlers {
		wg.Add(1)

		h.SetLogger(log)

		go d.invokeHandler(ctx, h, e, func(err error) {
			if err != nil {
				log.Error("failed to handle event", zap.Error(err))
			}

			wg.Done()
		})
	}

	return wg
}

func (d *Dispatcher) invokeHandler(
	ctx context.Context, hdl handler.Handler, e *slackevents.EventsAPIEvent, callback func(err error),
) {
	// context with timeout,  for this handler
	cx, cancel := context.WithTimeout(ctx, JobTimeout)
	defer cancel()

	// channel to receive the result from handler
	ch := make(chan error, 1)
	defer close(ch)

	go func() {
		ch <- hdl.Handle(cx, e)
	}()

	select {
	case err := <-ch: // finish handler before timeout
		// ignore error if receives because the error should be logged in handler code
		callback(err)
	case <-cx.Done(): // finish with reason(cancel/deadlineexceeded)
		err := cx.Err()
		callback(err)
	}
}
