package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/ollama"
)

func main() {
	ctx := context.Background()
	env := llm.NewEnv(os.Getenv)

	if err := run(ctx, env, os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, env llm.Env, args []string) error {
	in, err := parse(env, args)
	if err != nil {
		return err
	}

	provider, err := ollama.New(
		ollama.WithHost(in.host),
		ollama.WithModel(in.model),
	)
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf("Human: %s\nAssistant:", in.prompt)

	_, err = llm.Call(ctx, provider, prompt,
		llm.WithTemperature(in.llmTemperature),
		llm.WithStreamingFunc(stream),
	)

	fmt.Println()

	return err
}

func stream(ctx context.Context, chunk []byte) error {
	fmt.Print(string(chunk))

	return nil
}

type input struct {
	host           string
	model          string
	prompt         string
	llmTemperature float64
}

func parse(env llm.Env, args []string) (input, error) {
	in := input{}

	if len(args) == 0 {
		return in, fmt.Errorf("no args provided")
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&in.host, "ollama-host", env.String("OLLAMA_HOST", "localhost"), "Hostname where your Ollama server is running")
	flags.StringVar(&in.model, "ollama-model", env.String("OLLAMA_MODEL", "smollm2:135m"), "Model to use by Ollama")

	flags.Float64Var(&in.llmTemperature, "llm-temperature", 0.8, "LLM temperature to use")

	if err := flags.Parse(args[1:]); err != nil {
		return in, err
	}

	rest := flags.Args()

	if len(rest) > 0 {
		in.prompt = strings.Join(rest, " ")
	}

	return in, nil
}
