package event

import (
	"reflect"
	"sync"

	"github.com/hirakiuc/gonta-app/handler"

	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"
)

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
	d.log.Debug("Register handler", zap.String("event", eventType), zap.String("handler", reflect.TypeOf(h).Name()))

	handlers, ok := d.handlers[eventType]
	if ok {
		d.handlers[eventType] = append(handlers, h)
	} else {
		d.handlers[eventType] = []handler.Handler{h}
	}
}

func (d *Dispatcher) Dispatch(e *slackevents.EventsAPIEvent) *sync.WaitGroup {
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

		go func(hdl handler.Handler, evt *slackevents.EventsAPIEvent) {
			if err := hdl.Handle(evt); err != nil {
				log.Error("Failed to handle event:")
			}

			wg.Done()
		}(h, e)
	}

	return wg
}
