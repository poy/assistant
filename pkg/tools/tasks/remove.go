package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[injection.Group[taskTool]](func(ctx context.Context) injection.Group[taskTool] {
		return injection.AddToGroup[taskTool](ctx, taskTool{
			Tool: Remove(ctx),
		})
	})
}

func Remove(ctx context.Context) tools.Tool {
	s := injection.Resolve[Store](ctx)
	f := injection.Resolve[TaskFinder](ctx)
	return tools.Tool{
		Name:        "remove",
		Description: "Remove a task. The argument is the instructions from the user on the task. This tool takes care of figuring out the name, description, etc, so you don't have to. Just pass the instructions through to this tool.",
		Args: []string{
			"instructions",
		},
		Examples: []string{
			"remove the task to buy things for dinner for the next few days",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			fields := strings.Fields(input)
			if len(fields) == 0 {
				return "", errors.New("wrong number of arguments")
			}

			t, err := f.FindTask(ctx, input)
			if err != nil {
				return "", fmt.Errorf("failed to find task: %w", err)
			}
			s.Remove(t.Name())

			return fmt.Sprintf("Removed task %s", t.Name()), nil
		},
	}
}
