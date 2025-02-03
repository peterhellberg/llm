package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/chains/stuffdocumentschain"
	"github.com/peterhellberg/llm/providers/ollama"
)

func main() {
	ctx := context.Background()
	env := llm.NewEnv(os.Getenv)

	if err := run(ctx, env, os.Args, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, env llm.Env, args []string, w io.Writer) error {
	in, err := parse(args, env)
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

	inputValues := map[string]any{
		"input_documents": []llm.Document{
			{PageContent: "Harrison went to Harvard."},
			{PageContent: "Ankush went to Princeton."},
			{PageContent: "Caleb also went to Harvard."},
		},
		"question": "Who went to Harward?",
	}

	answer, err := llm.ChainCall(ctx, load(provider), inputValues)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, answer["text"])

	return nil
}

func load(provider llm.Provider) llm.Chain {
	content := `Use the following pieces of context to answer the question at the end.

If you don't know the answer, just say that you don't know, don't try to make up an answer.

{{.context}}

Question: {{.question}}
Helpful Answer:`

	return stuffdocumentschain.New(
		llm.NewChain(provider,
			llm.GoTemplate(content, []string{
				"context",
				"question",
			}),
		),
	)
}

type input struct {
	host  string
	model string
}

func parse(args []string, env llm.Env) (input, error) {
	in := input{}

	if len(args) == 0 {
		return in, fmt.Errorf("no args provided")
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&in.host, "ollama-host", env.String("OLLAMA_HOST", "localhost"), "Hostname where your Ollama server is running")
	flags.StringVar(&in.model, "ollama-model", env.String("OLLAMA_MODEL", "phi4"), "Model to use by Ollama")

	if err := flags.Parse(args[1:]); err != nil {
		return in, err
	}

	return in, nil
}
