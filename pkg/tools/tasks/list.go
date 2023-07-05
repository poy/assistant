package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[injection.Group[displayTaskTool]](func(ctx context.Context) injection.Group[displayTaskTool] {
		return injection.AddToGroup[displayTaskTool](ctx, displayTaskTool{
			Tool: List(ctx),
		})
	})
	injection.Register[injection.Group[taskFinderTool]](func(ctx context.Context) injection.Group[taskFinderTool] {
		return injection.AddToGroup[taskFinderTool](ctx, taskFinderTool{
			Tool: List(ctx),
		})
	})
}

// List returns a list of task names.
func List(ctx context.Context) tools.Tool {
	s := injection.Resolve[Store](ctx)
	return tools.Tool{
		Name:        "list",
		Description: "List all task names.",
		Run: func(ctx context.Context, input string) (string, error) {
			var names []string
			for _, name := range s.TaskNames() {
				names = append(names, fmt.Sprintf("* %s", name))
			}

			if len(names) > 0 {
				fmt.Println(strings.Join(names, "\n"))
			} else {
				fmt.Println("You don't have any tasks yet...")
			}
			return "Displayed the current list of tasks to the user", nil
		},
	}
}
