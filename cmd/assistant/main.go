package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-react/pkg/llms"
	"github.com/google/go-react/pkg/llms/vertex"
	"github.com/poy/assistant/pkg/tools/tasks"
	"github.com/poy/assistant/pkg/tools/userinput"
	"github.com/poy/go-dependency-injection/pkg/injection"
)

var model = flag.String("model", "text-bison@001", "The model to use for the prompt")
var apiEndpoint = flag.String("api-endpoint", "us-central1-aiplatform.googleapis.com", "The API endpoint to use")
var projectID = flag.String("project-id", os.Getenv("GCP_PROJECT_ID"), "The project ID to use")
var maxTokens = flag.Int("max-tokens", 1024, "The maximum number of tokens to generate")
var temperature = flag.Float64("temperature", 0.2, "The temperature to use for the prompt")
var topK = flag.Int("top-k", 40, "The top-k value to use for the prompt")
var topP = flag.Float64("top-p", 0.9, "The top-p value to use for the prompt")

func main() {
	log.SetFlags(0)
	flag.Parse()

	ctx := context.Background()
	registerLLM(ctx)
	setupStoreFilePath()

	ctx = injection.WithInjection(ctx)

	taskAgent := injection.Resolve[tasks.TaskAgent](ctx).Agent

	for {
		fmt.Println("AI: What is the goal?")
		fmt.Print("You: ")
		goal, err := userinput.ReadLine()
		if err != nil {
			log.Fatalf("failed to read line: %v", err)
		}

		finalAnswer, err := taskAgent.Run(ctx, goal)
		if err != nil {
			log.Fatalf("agent.Run failed: %v", err)
		}
		fmt.Println(finalAnswer)
	}
}

func registerLLM(ctx context.Context) {
	injection.Register[vertex.Params](
		func(ctx context.Context) vertex.Params {
			return vertex.Params{
				Model:       *model,
				MaxTokens:   *maxTokens,
				Temperature: *temperature,
				TopK:        *topK,
				TopP:        *topP,
			}
		},
	)
	injection.Register[llms.LLM[vertex.Params]](
		func(ctx context.Context) llms.LLM[vertex.Params] {
			return getLLM(ctx)
		},
	)
}

func getLLM(ctx context.Context) llms.LLM[vertex.Params] {
	if *projectID == "" {
		log.Fatalf("you must set the project-id flag or GCP_PROJECT_ID environment variable")
	}

	llm, err := vertex.New(ctx, *apiEndpoint, *projectID)
	if err != nil {
		log.Fatalf("failed to create LLM: %v", err)
	}
	return llm
}

func setupStoreFilePath() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home dir: %v", err)
	}

	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(home+"/.assistant", 0755); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	tasks.ProvideStorePath(home + "/.assistant/tasks.json")
}
