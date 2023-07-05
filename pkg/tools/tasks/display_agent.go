package tasks

import (
	"context"
	"fmt"

	"github.com/google/go-react/pkg/agents"
	"github.com/google/go-react/pkg/tools"
	assistanttools "github.com/poy/assistant/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[injection.Group[taskTool]](func(ctx context.Context) injection.Group[taskTool] {
		return injection.AddToGroup[taskTool](ctx, taskTool{
			Tool: newDisplayTaskTool(ctx),
		})
	})
}

type displayTaskTool struct {
	tools.Tool
}

func newDisplayTaskTool(ctx context.Context) tools.Tool {
	// finder := injection.Resolve[TaskFinder](ctx)
	agent := newDisplayTaskAgent(ctx)
	return tools.Tool{
		Name:        "display",
		Description: "Display information about the tasks. The data will be showed to the user instead of being returned to the agent. You don't need to know anything about the data, let me handle it.",
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

			return agent.Run(ctx, input)
		},
	}
}

func newDisplayTaskAgent(ctx context.Context) agents.Agent[string] {
	return assistanttools.AgentBuilder[string, displayTaskTool](
		ctx,
		assistanttools.WithName[string]("TaskDisplayer"),
		assistanttools.WithPreamble[string](displayAgentPreamble),
		assistanttools.WithExamples[string](displayAgentExamples),
	)
}

const (
	displayAgentPreamble = "Help the user figure out what tool to use based on the question."
)

var (
	displayAgentExamples = []agents.PromptDataExample[string]{
		{
			Question: "What tasks do I have?",
			Output: agents.Reasoning[string]{
				Thought: "I should use the list tool",
				Action:  "list",
				Input:   "",
			},
		},
		{
			Question: "Show me the details of the grocery store task",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the read tool",
						Action:  "read",
						Input:   "display the details of the grocery store task",
					},
					Observation: "Displayed of the task pick up groceries at the store were shown to the user",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "Displayed of the task pick up groceries at the store were shown to the user",
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
