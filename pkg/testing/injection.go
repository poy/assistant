package testing

import (
	"context"

	"github.com/google/go-react/pkg/llms"
	llmstesting "github.com/google/go-react/pkg/llms/testing"
	"github.com/google/go-react/pkg/llms/vertex"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func init() {
	injection.Register[vertex.Params](
		func(ctx context.Context) vertex.Params {
			return vertex.Params{}
		},
	)
	injection.Register[llms.LLM[vertex.Params]](
		func(ctx context.Context) llms.LLM[vertex.Params] {
			return &llmstesting.Fake[vertex.Params]{}
		},
	)
}
