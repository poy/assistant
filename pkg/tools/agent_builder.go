package tools

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/google/go-react/pkg/agents"
	"github.com/google/go-react/pkg/llms"
	"github.com/google/go-react/pkg/llms/vertex"
	"github.com/google/go-react/pkg/predictors"
	"github.com/google/go-react/pkg/prompters"
	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

// Option is an option for the agent.
type Option[TOut any] func(*agentBuilder[TOut])

type agentBuilder[TOut any] struct {
	name     string
	preamble string
	rules    []string
	examples []agents.PromptDataExample[TOut]
}

// AgentBuilder builds an agent.
func AgentBuilder[TOut, TToolGroup any](
	ctx context.Context,
	opts ...Option[TOut],
) agents.Agent[TOut] {
	b := &agentBuilder[TOut]{
		name: "Agent",
	}
	for _, o := range opts {
		o(b)
	}

	params := injection.Resolve[vertex.Params](ctx)
	llm := injection.Resolve[llms.LLM[vertex.Params]](ctx)

	var promptOpts []prompters.Option[agents.PromptData[TOut]]
	if b.preamble != "" {
		promptOpts = append(promptOpts, agents.WithPreamble[vertex.Params, TOut](b.preamble))
	}
	if len(b.rules) > 0 {
		promptOpts = append(promptOpts, agents.WithRules[vertex.Params, TOut](b.rules...))
	}
	if len(b.examples) > 0 {
		promptOpts = append(promptOpts, agents.WithExamples[vertex.Params, TOut](b.examples...))
	}

	prompt := agents.NewDefaultPrompt[vertex.Params, TOut](
		params,
		promptOpts...,
	)
	prompt = prompters.NewLogger(prompt, os.Stderr)
	parser := agents.NewDefaultParser[TOut]()
	predictor := predictors.New(llm, prompt, parser)
	predictor = predictors.NewRetrier(predictor)
	predictor = agents.NewCLILogger(predictor, os.Stderr, agents.WithCLILoggerPrefix[TOut](b.name))
	ts := injection.Resolve[injection.Group[TToolGroup]](ctx).Vals()

	var toolSet []tools.Tool

	for _, t := range ts {
		val := reflect.ValueOf(&t).Elem()
		toolVal := val.FieldByName("Tool")
		tool, ok := toolVal.Addr().Interface().(*tools.Tool)
		if !ok {
			panic(fmt.Sprintf(`type %T does not have a Tool field. You should use an embedded type:
type MyTool struct {
  tools.Tool
}`, t))
		}
		toolSet = append(toolSet, *tool)
	}

	return agents.NewAgent(predictor, toolSet...)
}

// WithName sets the name for the agent.
func WithName[TOut any](name string) Option[TOut] {
	return func(b *agentBuilder[TOut]) {
		b.name = name
	}
}

// WithPreamble sets the preamble for the agent.
func WithPreamble[TOut any](preamble string) Option[TOut] {
	return func(b *agentBuilder[TOut]) {
		b.preamble = preamble
	}
}

// WithRules sets the rules for the agent.
func WithRules[TOut any](rules []string) Option[TOut] {
	return func(b *agentBuilder[TOut]) {
		b.rules = rules
	}
}

// WithExamples sets the examples for the agent.
func WithExamples[TOut any](examples []agents.PromptDataExample[TOut]) Option[TOut] {
	return func(b *agentBuilder[TOut]) {
		b.examples = examples
	}
}
