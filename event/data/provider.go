package data

import (
	"net/http"
	"sync"
)

type Generator func(e *ExternalDataRequest) ([]byte, error)

type Provider struct {
	providers map[string]Generator
	wg        *sync.WaitGroup
}

func NewProvider() *Provider {
	return &Provider{
		providers: map[string]Generator{},
		wg:        &sync.WaitGroup{},
	}
}

func (q *Provider) AddProvider(blockID string, generator Generator) {
	q.providers[blockID] = generator
}

func (q *Provider) Wait() {
	q.wg.Wait()
}

func (q *Provider) Process(e *ExternalDataRequest, w http.ResponseWriter) error {
	generator, ok := q.providers[e.BlockID]
	if !ok {
		return nil
	}

	q.wg.Add(1)
	defer q.wg.Done()

	bytes, err := generator(e)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
