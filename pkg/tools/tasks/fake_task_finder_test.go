package tasks_test

import (
	"context"

	"github.com/poy/assistant/pkg/tools/tasks"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[tasks.TaskFinder](
		func(ctx context.Context) tasks.TaskFinder {
			return &fakeTaskFinder{
				errs: make(map[string]error),
				s:    injection.Resolve[tasks.Store](ctx),
			}
		},
	)
}

type fakeTaskFinder struct {
	errs map[string]error
	s    tasks.Store
}

func (f *fakeTaskFinder) FindTask(ctx context.Context, name string) (*tasks.Task, error) {
	return f.GetTask(name), f.errs[name]
}

func (f *fakeTaskFinder) Add(name, description string) {
	f.s.Add(name, description)
}

func (f *fakeTaskFinder) AddErr(name string, err error) {
	f.errs[name] = err
}

func (f *fakeTaskFinder) GetTask(name string) *tasks.Task {
	return f.s.GetTask(name)
}
