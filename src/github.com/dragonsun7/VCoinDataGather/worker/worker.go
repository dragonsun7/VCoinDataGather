package worker

import "sync"

type Worker interface {
	Start(wg *sync.WaitGroup)
}

type BaseWorker struct {
	isCancel bool
	err      error
}

func (worker *BaseWorker) Stop() {
	worker.isCancel = true
}

func (worker *BaseWorker) Start(wg *sync.WaitGroup) {
	return
}
