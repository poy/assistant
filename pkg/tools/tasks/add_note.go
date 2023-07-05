package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/predictors"
	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[injection.Group[modifyTaskTool]](func(ctx context.Context) injection.Group[modifyTaskTool] {
		return injection.AddToGroup[modifyTaskTool](ctx, modifyTaskTool{
			Tool: AddNote(ctx),
		})
	})
}

// AddNote adds a note to a task.
func AddNote(ctx context.Context) tools.Tool {
	f := injection.Resolve[TaskFinder](ctx)
	taskRewriter := injection.Resolve[predictors.Predictor[rewriteTaskParams, string]](ctx)

	return tools.Tool{
		Name:        "add-note",
		Description: "Add a note to the task. Provide the instructions from the user.",
		Args: []string{
			"instructions",
		},
		Examples: []string{
			"add a note to the grocery store task to buy eggs",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			fields := strings.Fields(input)
			if len(fields) == 0 {
				return "", errors.New("wrong number of arguments")
			}
			instructions := strings.Join(fields, " ")

			t, err := f.FindTask(ctx, instructions)
			if err != nil {
				return "", fmt.Errorf("failed to find task: %w", err)
			}

			if t == nil {
				return "", fmt.Errorf("task %q not found. Try listing the tasks to find the right one.", fields[0])
			}

			instructions, err = taskRewriter.Predict(ctx, rewriteTaskParams{
				Description: instructions,
			})
			if err != nil {
				return "", fmt.Errorf("failed to rewrite task description: %w", err)
			}

			t.AddNotes(instructions)

			return fmt.Sprintf("done adding note to %s", t.Name()), nil
		},
	}
}
