package zutils

import (
	"context"

	"github.com/panjf2000/ants/v2"

	"github.com/Bishoptylaor/go-toolbox/zlog"
)

// Pool alias of ants.Pool
type Pool struct {
	ap *ants.Pool
}

// NewWorkerPoolWithOptions 创建Worker池, 支持传入ants Option.
func NewWorkerPoolWithOptions(size int, options ...ants.Option) (*Pool, error) {
	options = append(options, ants.WithLogger(&WorkerLogger{}))
	ap, err := ants.NewPool(size, options...)
	if err != nil {
		return nil, err
	}
	return &Pool{ap: ap}, nil
}

// NewWorkerPool constructor of Pool
func NewWorkerPool(size int) (*Pool, error) {
	return NewWorkerPoolWithOptions(size)
}

// Tune 调整Worker池容量
func (p *Pool) Tune(size int) {
	p.ap.Tune(size)
}

// Release close pool and release resources
func (p *Pool) Release() {
	p.ap.Release()
}

// Submit submit a task
func (p *Pool) Submit(task func()) error {
	return p.ap.Submit(task)
}

// Running return goroutines of runnning
func (p *Pool) Running() int {
	return p.ap.Running()
}

// WorkerLogger log handler
type WorkerLogger struct {
}

// Printf implements ants.Logger
func (WorkerLogger) Printf(format string, args ...interface{}) {
	zlog.Infof(context.Background(), format, args...)
}
