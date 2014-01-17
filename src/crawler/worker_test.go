package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkerChan(t *testing.T) {
	workerChan := make(workerChan)
	worker := NewWorker(workerChan, nil)

	worker.Run()
	w := <-workerChan
	assert.Equal(t, w, worker)
}

func TestWorkerProcessing(t *testing.T) {
	respChan := make(responseChan)
	workerChan := make(workerChan)
	worker := NewWorker(workerChan, respChan)

	worker.Run()
	w := <-workerChan
	w.Process("http://nox73.ru/")

	response := <-respChan
	assert.True(t, response.success)
}
