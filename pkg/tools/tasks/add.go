package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-react/pkg/llms"
	"github.com/google/go-react/pkg/llms/vertex"
	"github.com/google/go-react/pkg/parsers"
	"github.com/google/go-react/pkg/predictors"
	"github.com/google/go-react/pkg/prompters"
	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[injection.Group[taskTool]](func(ctx context.Context) injection.Group[taskTool] {
		return injection.AddToGroup[taskTool](ctx, taskTool{
			Tool: Add(ctx),
		})
	})
	setupTaskTitleGenerator()
	setupTaskRewriter()
}

// Add returns a tool that adds tasks.
func Add(ctx context.Context) tools.Tool {
	s := injection.Resolve[Store](ctx)

	taskTitlePredictor := injection.Resolve[predictors.Predictor[generateTaskTitleParams, string]](ctx)
	taskRewriter := injection.Resolve[predictors.Predictor[rewriteTaskParams, string]](ctx)

	return tools.Tool{
		Name:        "add",
		Description: "Add a task. The argument is the instructions from the user on the task. This tool takes care of figuring out the name, description, etc, so you don't have to. Just pass the instructions through to this tool.",
		Args: []string{
			"instructions",
		},
		Examples: []string{
			"buy things for dinner for the next few days",
		},
		Run: func(ctx context.Context, input string) (string, error) {
			fields := strings.Fields(input)
			if len(fields) == 0 {
				return "", errors.New("wrong number of arguments")
			}

			description := strings.Join(fields, " ")

			title, err := taskTitlePredictor.Predict(ctx, generateTaskTitleParams{
				Description: description,
			})
			if err != nil {
				return "", err
			}

			description, err = taskRewriter.Predict(ctx, rewriteTaskParams{
				Description: description,
			})
			if err != nil {
				return "", err
			}

			s.Add(title, description)

			return fmt.Sprintf("Added task %q - %s", title, description), nil
		},
	}
}

const (
	generateTaskTitlePromptTempl = `Given the description of the task, create a good title for it:

Task description: {{.Description}}
Output: `
)

type generateTaskTitleParams struct {
	Description string
}

const (
	rewriteTaskPromptTempl = `Given the description of the task, rewrite it to be more concise:

Task description: {{.Description}}
Output: `
)

type rewriteTaskParams struct {
	Description string
}

func setupTaskTitleGenerator() {
	injection.Register[predictors.Predictor[generateTaskTitleParams, string]](
		func(ctx context.Context) predictors.Predictor[generateTaskTitleParams, string] {
			llm := injection.Resolve[llms.LLM[vertex.Params]](ctx)
			params := injection.Resolve[vertex.Params](ctx)

			prompter := prompters.NewTextTemplate[generateTaskTitleParams, vertex.Params](
				generateTaskTitlePromptTempl,
				params,
			)
			parser := parsers.NewTextParser()
			predictor := predictors.New(llm, prompter, parser)
			predictor = predictors.NewRetrier(predictor)
			return predictor
		},
	)
}

func setupTaskRewriter() {
	injection.Register[predictors.Predictor[rewriteTaskParams, string]](
		func(ctx context.Context) predictors.Predictor[rewriteTaskParams, string] {
			llm := injection.Resolve[llms.LLM[vertex.Params]](ctx)
			params := injection.Resolve[vertex.Params](ctx)

			prompter := prompters.NewTextTemplate[rewriteTaskParams, vertex.Params](
				rewriteTaskPromptTempl,
				params,
			)
			parser := parsers.NewTextParser()
			predictor := predictors.New(llm, prompter, parser)
			predictor = predictors.NewRetrier(predictor)
			return predictor
		},
	)
}
