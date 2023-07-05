package userinput

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/go-react/pkg/tools"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

func Register[TToolGroup any]() {
	injection.Register[injection.Group[TToolGroup]](func(ctx context.Context) injection.Group[TToolGroup] {
		var t TToolGroup

		val := reflect.ValueOf(&t).Elem()
		toolVal := val.FieldByName("Tool")

		// This panicking probably implies that the underlying type needs to embed
		// tools.Tool.
		u := reflect.ValueOf(UserInputTool())
		toolVal.Set(u)

		return injection.AddToGroup[TToolGroup](ctx, t)
	})
}

// UserInputTool asks the user for input.
func UserInputTool() tools.Tool {
	return tools.Tool{
		Name:        "user-input",
		Description: "Ask the user a question. The input is what is displayed to the user.",
		Examples:    []string{"some question to the user"},
		Run: func(ctx context.Context, input string) (string, error) {
			fmt.Printf("AI: %s\n", input)

			// Read a line from the user.
			fmt.Print("You: ")
			return ReadLine()
		},
	}
}
