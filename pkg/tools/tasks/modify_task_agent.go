package tasks

import (
	"context"
	"fmt"

	"github.com/google/go-react/pkg/agents"
	"github.com/google/go-react/pkg/tools"
	assistanttools "github.com/poy/assistant/pkg/tools"
	"github.com/poy/assistant/pkg/tools/userinput"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	// Ensure the modify toolset has the user input tool.
	userinput.Register[modifyTaskTool]()

	injection.Register[injection.Group[taskTool]](func(ctx context.Context) injection.Group[taskTool] {
		return injection.AddToGroup[taskTool](ctx, taskTool{
			Tool: newModifyTaskTool(ctx),
		})
	})
}

func newModifyTaskTool(ctx context.Context) tools.Tool {
	finder := injection.Resolve[TaskFinder](ctx)
	agent := newModifyTaskAgent(ctx)
	return tools.Tool{
		Name:        "modify",
		Description: "Modify the tasks. The tool just wants the raw instructions, it doesn't need you to get the task name.",
		Args: []string{
			"instructions",
		},
		Examples: []string{
			"add a note to the grocery store task to buy milk",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			if input == "" {
				return "", fmt.Errorf("no instructions provided")
			}

			// Figure out the task name from the instructions.
			task, err := finder.FindTask(ctx, input)
			if err != nil {
				return "", fmt.Errorf("finding task: %w", err)
			}

			return agent.Run(ctx, fmt.Sprintf("for the task %s, do the following: %s", task.Name(), input))
		},
	}
}

func newModifyTaskAgent(ctx context.Context) agents.Agent[string] {
	return assistanttools.AgentBuilder[string, modifyTaskTool](
		ctx,
		assistanttools.WithName[string]("TaskModifier"),
		assistanttools.WithPreamble[string](modifyTaskPreamble),
		assistanttools.WithExamples[string](modifyTaskExamples),
	)
}

type modifyTaskParams struct {
	instructions string
	taskName     string
}

type modifyTaskTool struct {
	tools.Tool
}

const (
	modifyTaskPreamble = `Using the given tools, help the user modify the given task. You never need to confirm things with the user, just do it.`
)

var (
	modifyTaskExamples = []agents.PromptDataExample[string]{
		{
			Question: "for the task change the tires, do the following: add a note that we need a spare too",
			Output: agents.Reasoning[string]{
				Thought: "I should use the add-note tool",
				Action:  "add-note",
				Input:   "add a note to the change the tires task that a spare is needed too",
			},
		},
		{
			Question: "for the task change the tires, do the following: add a note",
			Output: agents.Reasoning[string]{
				Thought: "I should use the user-input tool",
				Action:  "user-input",
				Input:   "What should the note include?",
			},
		},
		{
			Question: "for the task change the tires, do the following: add a note",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the user-input tool",
						Action:  "user-input",
						Input:   "What should the note include?",
					},
					Observation: "the note should include that we need a spare",
				},
			},
			Output: agents.Reasoning[string]{
				Thought: "I should use the add-note tool",
				Action:  "add-note",
				Input:   "add a note to the change the tires task that a spare is needed too",
			},
		},
		{
			Question: "for the task change the tires, do the following: add a note",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the user-input tool",
						Action:  "user-input",
						Input:   "What should the note include?",
					},
					Observation: "the note should include that we need a spare",
				},
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the add-note tool",
						Action:  "add-note",
						Input:   "add a note to the change the tires task that a spare is needed too",
					},
					Observation: "added the note",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "added the note to the change the tires task that a spare is needed too",
			},
		},
		{
			Question: "for the task change the tires, do the following: add a note",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the user-input tool",
						Action:  "user-input",
						Input:   "What should the note include?",
					},
					Observation: "nevermind, I want to add a new task to buy some drinks",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "Instead of adding a note, the user now wants to add a new task to buy drinks",
			},
		},
	}
)
