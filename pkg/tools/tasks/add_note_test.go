package tasks_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/poy/assistant/pkg/tools/tasks"
	"github.com/poy/go-dependency-injection/pkg/injection"
	injectiontesting "github.com/poy/go-dependency-injection/pkg/injection/testing"
)

func TestAddNote(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		input  string
		setup  func(f *fakeTaskFinder)
		assert func(t *testing.T, val string, err error, f *fakeTaskFinder)
	}{
		{
			name:  "adds note",
			input: "task 1",
			setup: func(f *fakeTaskFinder) {
				f.Add("task 1", "")
			},
			assert: func(t *testing.T, val string, err error, f *fakeTaskFinder) {
				if err != nil {
					t.Fatal(err)
				}

				if actual, expected := len(f.GetTask("task 1").Notes()), 1; actual != expected {
					t.Fatalf("expected %d, got %d", expected, actual)
				}
				if actual, expected := val, "done adding note to task 1"; actual != expected {
					t.Fatalf("expected %q, got %q", expected, actual)
				}
			},
		},
		{
			name:  "unknown task",
			input: "task 1",
			assert: func(t *testing.T, val string, err error, f *fakeTaskFinder) {
				if err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name:  "task finder returns an error",
			input: "task 1",
			setup: func(f *fakeTaskFinder) {
				f.Add("task 1", "")
				f.AddErr("task 1", errors.New("some-error"))
			},
			assert: func(t *testing.T, val string, err error, f *fakeTaskFinder) {
				if err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name:  "too few arguments",
			input: "",
			assert: func(t *testing.T, val string, err error, f *fakeTaskFinder) {
				if actual, expected := fmt.Sprint(err), "wrong number of arguments"; actual != expected {
					t.Fatalf("expected %q, got %q", actual, expected)
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

			f := injection.Resolve[tasks.TaskFinder](ctx).(*fakeTaskFinder)
			if tc.setup != nil {
				tc.setup(f)
			}
			result, err := tasks.AddNote(ctx).Run(context.Background(), tc.input)
			tc.assert(t, result, err, f)
		})
	}
}
