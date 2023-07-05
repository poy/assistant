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
			Tool: Read(ctx),
		})
	})
}

// Read returns a task.
func Read(ctx context.Context) tools.Tool {
	f := injection.Resolve[TaskFinder](ctx)
	return tools.Tool{
		Name:        "read",
		Description: "Read a task by its name. If you don't know the name, I can guess it.",
		Args: []string{
			"name",
		},
		Examples: []string{
			"buy groceries",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			t, err := f.FindTask(ctx, input)
			if err != nil {
				return "", fmt.Errorf("failed to find task: %w", err)
			}

			result := fmt.Sprintf(`
Name: %s
Description: %s
`, t.Name(), t.Description())

			if t.Completed() {
				result = fmt.Sprintf("%s\nCompleted At: %s", result, t.CompletedAt())
			}

			if len(t.Notes()) > 0 {
				var notes []string
				for _, n := range t.Notes() {
					notes = append(notes, fmt.Sprintf("* [%s] %s", n.Datetime(), n.Note()))
				}
				result = fmt.Sprintf("%s\nNotes:\n%s", result, strings.Join(notes, "\n"))
			}

			fmt.Println(result)

			return fmt.Sprintf("Details of the task %s were shown to the user", input), nil
		},
	}
}
