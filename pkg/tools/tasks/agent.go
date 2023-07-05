package tasks

import (
	"context"

	"github.com/google/go-react/pkg/agents"
	"github.com/google/go-react/pkg/tools"
	assistanttools "github.com/poy/assistant/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

// TaskAgent is an agent that can be used to interact with tasks.
type TaskAgent struct {
	agents.Agent[string]
}

func init() {
	injection.Register[TaskAgent](func(ctx context.Context) TaskAgent {
		return TaskAgent{
			Agent: assistanttools.AgentBuilder[string, taskTool](
				ctx,
				assistanttools.WithName[string]("Tasks"),
				assistanttools.WithPreamble[string](rootAgentPreamble),
				assistanttools.WithExamples[string](rootAgentExamples),
			),
		}
	})
}

// taskTool wraps a tools.Tool so it can be grouped and used by injection.
type taskTool struct {
	tools.Tool
}

const (
	rootAgentPreamble = "Using the given tools, help the user manage their tasks:"
)

var (
	rootAgentExamples = []agents.PromptDataExample[string]{
		{
			Question: "What tasks do I have?",
			Output: agents.Reasoning[string]{
				Thought: "I should use the display tool",
				Action:  "display",
				Input:   "show the user the tasks",
			},
		},
		{
			Question: "Show me the details of the grocery store task",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the display tool",
						Action:  "display",
						Input:   "display the details of the grocery store task",
					},
					Observation: "Displayed the details of the pick up groceries at the store task",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "The details of the grocery store task were displayed",
			},
		},
		{
			Question: "Add a task",
			Output: agents.Reasoning[string]{
				Thought: "I should use the add tool",
				Action:  "add",
				Input:   "add a task",
			},
		},
		{
			Question: "Remove a task",
			Output: agents.Reasoning[string]{
				Thought: "I should use the remove tool",
				Action:  "remove",
				Input:   "remove a task",
			},
		},
		{
			Question: "add a note to the grocery store task that I need milk and eggs.",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the modify tool",
						Action:  "modify",
						Input:   "Add a note to the grocery store task to purchase milk and eggs",
					},
					Observation: "The task was modified",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "The task was modified",
			},
		},
		{
			Question: "add two tasks, one to pick up the dry cleaning and one to pick up the groceries",
			PreviousContext: []agents.ThoughtIteration[string]{
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "I should use the add tool twice. First to pick up the dry cleaning, and again to pick up the groceries",
						Action:  "add",
						Input:   "add a task to pick up the dry cleaning",
					},
					Observation: "The task was added",
				},
				{
					Reasoning: agents.Reasoning[string]{
						Thought: "Now I should use the add tool again to add the task to pick up the groceries",
						Action:  "add",
						Input:   "add a task to pick up the groceries",
					},
					Observation: "The task was added",
				},
			},
			Output: agents.Reasoning[string]{
				Thought:     "I know the answer",
				FinalAnswer: "The tasks were added",
			},
		},
		{
			Question: "Build a spaceship",
			Output: agents.Reasoning[string]{
				Thought:     "I don't know how to build a spaceship. Better just let the user know.",
				FinalAnswer: "I can't build a spaceship",
			},
		},
	}
)
