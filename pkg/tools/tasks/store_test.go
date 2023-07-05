package tasks_test

import (
	"testing"

	"github.com/poy/assistant/pkg/tools/tasks"
	"github.com/poy/go-dependency-injection/pkg/injection"
	injectiontesting "github.com/poy/go-dependency-injection/pkg/injection/testing"
)

func TestStore(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		setup  func(tasks.Store)
		assert func(*testing.T, tasks.Store)
	}{
		{
			name: "no tasks",
			assert: func(t *testing.T, s tasks.Store) {
				if len(s.TaskNames()) != 0 {
					t.Errorf("expected 0 tasks, got %d", len(s.TaskNames()))
				}
			},
		},
		{
			name: "one task",
			setup: func(s tasks.Store) {
				s.Add("some-task", "some-description")
			},
			assert: func(t *testing.T, s tasks.Store) {
				if len(s.TaskNames()) != 1 {
					t.Errorf("expected 1 task, got %d", len(s.TaskNames()))
				}
			},
		},
		{
			name: "removes task",
			setup: func(s tasks.Store) {
				s.Add("some-task", "some-description")
				s.Add("some-other-task", "some-description")
			},
			assert: func(t *testing.T, s tasks.Store) {
				s.Remove("some-TASK")
				if len(s.TaskNames()) != 1 {
					t.Errorf("expected 1 task, got %d", len(s.TaskNames()))
				}
				if actual, expected := s.TaskNames()[0], "some-other-task"; actual != expected {
					t.Errorf("expected %q, got %q", expected, actual)
				}
			},
		},
		{
			name: "add duplicate task",
			setup: func(s tasks.Store) {
				s.Add("some-TASK", "some-description")
			},
			assert: func(t *testing.T, s tasks.Store) {
				s.Add("some-task", "some-other-description")
				if actual, expected := s.GetTask("some-task").Description(), "some-description"; actual != expected {
					t.Errorf("expected %q, got %q", expected, actual)
				}
			},
		},
		{
			name: "complete task",
			setup: func(s tasks.Store) {
				s.Add("some-task", "some-description")
				s.Add("some-other-task", "some-description")
				s.GetTask("some-task").Complete()
			},
			assert: func(t *testing.T, s tasks.Store) {
				if actual, expected := s.GetTask("some-task").CompletedAt().IsZero(), false; actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
				if actual, expected := s.GetTask("some-other-task").CompletedAt().IsZero(), true; actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			},
		},
		{
			name: "add notes",
			setup: func(s tasks.Store) {
				s.Add("some-task", "some-description")
				s.GetTask("some-task").AddNotes("some-note", "some-other-note")
			},
			assert: func(t *testing.T, s tasks.Store) {
				if actual, expected := len(s.GetTask("some-task").Notes()), 2; actual != expected {
					t.Errorf("expected %d, got %d", expected, actual)
				}
				if actual, expected := s.GetTask("some-task").Notes()[0].Note(), "some-note"; actual != expected {
					t.Errorf("expected %q, got %q", expected, actual)
				}
				if actual, expected := s.GetTask("some-task").Notes()[0].Datetime().IsZero(), false; actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
				if actual, expected := s.GetTask("some-task").Notes()[1].Datetime().IsZero(), false; actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			},
		},
	}

	for _, tc := range testCases {
		// Avoid issues with closure.
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := injectiontesting.WithTesting(t)

			s := injection.Resolve[tasks.Store](ctx)
			if tc.setup != nil {
				tc.setup(s)
			}
			tc.assert(t, s)
		})
	}
}
