package tasks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-react/pkg/llms"
	llmstesting "github.com/google/go-react/pkg/llms/testing"
	"github.com/google/go-react/pkg/llms/vertex"
	"github.com/poy/assistant/pkg/tools/tasks"
	"github.com/poy/go-dependency-injection/pkg/injection"
	injectiontesting "github.com/poy/go-dependency-injection/pkg/injection/testing"

	// Register the fake LLM.
	_ "github.com/poy/assistant/pkg/testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		input  string
		setup  func(llm *llmstesting.Fake[vertex.Params])
		assert func(t *testing.T, val string, err error, s tasks.Store)
	}{
		{
			name:  "adds task",
			input: "some_name some_description",
			setup: func(llm *llmstesting.Fake[vertex.Params]) {
				llm.AlwaysText = "some-llm-output"
			},
			assert: func(t *testing.T, val string, err error, s tasks.Store) {
				if err != nil {
					t.Fatal(err)
				}

				if expected, actual := 1, len(s.TaskNames()); actual != expected {
					t.Fatalf("expected %d, got %d", expected, actual)
				}
				if expected, actual := "some-llm-output", s.TaskNames()[0]; actual != expected {
					t.Fatalf("expected %q, got %q", expected, actual)
				}
				if expected, actual := "some-llm-output", s.GetTask("some-llm-output").Name(); actual != expected {
					t.Fatalf("expected %q, got %q", expected, actual)
				}
				if expected, actual := "some-llm-output", s.GetTask("some-llm-output").Description(); actual != expected {
					t.Fatalf("expected %q, got %q", expected, actual)
				}
				if expected, actual := "Added task \"some-llm-output\" - some-llm-output", val; actual != expected {
					t.Fatalf("expected %q, got %q", expected, actual)
				}
			},
		},
		{
			name:  "too few arguments",
			input: "",
			assert: func(t *testing.T, val string, err error, s tasks.Store) {
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

			llm := injection.Resolve[llms.LLM[vertex.Params]](ctx).(*llmstesting.Fake[vertex.Params])
			s := injection.Resolve[tasks.Store](ctx)

			if tc.setup != nil {
				tc.setup(llm)
			}

			result, err := tasks.Add(ctx).Run(context.Background(), tc.input)
			tc.assert(t, result, err, s)
		})
	}
}
