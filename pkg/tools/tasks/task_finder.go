package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/agents"
	"github.com/google/go-react/pkg/tools"
	assistanttools "github.com/poy/assistant/pkg/tools"
	"github.com/poy/assistant/pkg/tools/userinput"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	// Ensure the TaskFinder toolset has the user input tool.
	userinput.Register[taskFinderTool]()

	injection.Register[TaskFinder](
		func(ctx context.Context) TaskFinder {
			return newTaskFinder(ctx)
		},
	)
}

// TaskFinder finds a task by name using the LLM.
type TaskFinder interface {
	// FindTask finds a task by name using the LLM.
	FindTask(ctx context.Context, taskName string) (*Task, error)
}

type taskFinder struct {
	agent agents.Agent[string]
	s     Store
}

type taskFinderTool struct {
	tools.Tool
}

// newTaskFinder creates a new TaskFinder.
func newTaskFinder(ctx context.Context) TaskFinder {
	s := injection.Resolve[Store](ctx)

	agent := assistanttools.AgentBuilder[string, taskFinderTool](
		ctx,
		assistanttools.WithName[string]("TaskFinder"),
		assistanttools.WithPreamble[string](taskFinderPreamble),
		assistanttools.WithExamples[string](taskFinderExamples),
	)
	return &taskFinder{
		agent: agent,
		s:     s,
	}
}

// FindTask finds a task by name using the LLM.
func (t *taskFinder) FindTask(ctx context.Context, taskName string) (*Task, error) {
	for i := 0; i < 3; i++ {
		name, err := t.agent.Run(
			ctx,
			fmt.Sprintf(
				"Which of the tasks (%s) do you think the user is looking for when they say: %s ",
				strings.Join(t.s.TaskNames(), ", "),
				taskName,
			),
		)

		if err != nil {
			return nil, fmt.Errorf("failed to find task: %w", err)
		}

		task := t.s.GetTask(name)
		if task == nil {
			continue
		}
		return task, nil
	}
	return nil, fmt.Errorf("could not find task %q", taskName)
}

func taskFinderToolSet(ctx context.Context, userInput tools.Tool) []tools.Tool {
	return []tools.Tool{
		userInput,
		Show(ctx),
	}
}

const (
	taskFinderPreamble = "Using the given tools, help the root agent figure out what task the user is talking about."
)

var (
	taskFinderExamples = []agents.PromptDataExample[string]{
		{
			Question: "Which of the tasks (pickup clothes, hire nanny, buy groceries) do you think the user is looking for when they say: xyz",
			Output: agents.Reasoning[string]{
				Thought: "I should show the user the tasks and then ask for them to be more clear",
				Action:  "list",
				Input:   "",
			},
		},
		{
			Question: "Which of the tasks (pickup clothes, hire nanny, buy groceries) do you think the user is looking for when they say: food shopping",
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "buy groceries",
			},
		},
		{
			Question: "Which of the tasks (pickup clothes, hire nanny, buy groceries) do you think the user is looking for when they say: xyz",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should show the user the tasks and then ask for them to be more clear",
						Action:  "list",
						Input:   "",
					},
					Observation: "The user has been shown the tasks",
				},
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should now ask the user to be more clear",
						Action:  "user-input",
						Input:   "Which task were you referring to?",
					},
					Observation: "the nanny one",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "hire nanny",
			},
		},
		{
			Question: "Build a spaceship",
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "I can't build a spaceship",
			},
		},
	}
)
