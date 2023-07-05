package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

// func init() {
// 	injection.Register[injection.Group[displayTaskTool]](func(ctx context.Context) injection.Group[displayTaskTool] {
// 		return injection.AddToGroup[displayTaskTool](ctx, displayTaskTool{
// 			Tool: Show(ctx),
// 		})
// 	})
// 	injection.Register[injection.Group[taskFinderTool]](func(ctx context.Context) injection.Group[taskFinderTool] {
// 		return injection.AddToGroup[taskFinderTool](ctx, taskFinderTool{
// 			Tool: Show(ctx),
// 		})
// 	})
// }

// Show returns a list of task names.
func Show(ctx context.Context) tools.Tool {
	s := injection.Resolve[Store](ctx)
	return tools.Tool{
		Name:        "show",
		Description: "Show/list/display tasks to the user. The response will not be availbe to the agent. The input is the instructions from the user.",
		Args: []string{
			"instructions",
		},
		Examples: []string{
			"show me the task about buying groceries",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			var names []string
			for _, name := range s.TaskNames() {
				names = append(names, fmt.Sprintf("* %s", name))
			}

			fmt.Println(strings.Join(names, "\n"))
			return "Displayed the current list of tasks to the user", nil
		},
	}
}
