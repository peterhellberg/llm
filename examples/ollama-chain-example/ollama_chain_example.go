package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/peterhellberg/llm"
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

	{ // First example
		prompt := llm.GoTemplate(
			"What is a good name for a company that makes {{.product}}?",
			[]string{"product"},
		)

		chain := llm.NewChain(provider, prompt)

		out, err := llm.ChainRun(ctx, chain, "socks")
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	{ // Second example
		translatePrompt := llm.GoTemplate(
			"Translate the following text from {{.inputLanguage}} to {{.outputLanguage}}. {{.text}}",
			[]string{"inputLanguage", "outputLanguage", "text"},
		)

		chain := llm.NewChain(provider, translatePrompt)

		fmt.Fprintf(w, "\n-------\n\n")

		// Otherwise the call function must be used.
		outputValues, err := llm.ChainCall(ctx, chain, map[string]any{
			"inputLanguage":  "English",
			"outputLanguage": "German",
			"text":           "I love programming computers.",
		})
		if err != nil {
			return err
		}

		outputKeys := chain.OutputKeys()

		out, ok := outputValues[outputKeys[0]].(string)
		if !ok {
			return fmt.Errorf("invalid chain return")
		}

		fmt.Println(out)
	}

	return nil
}

type input struct {
	host   string
	model  string
	prompt string
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

	rest := flags.Args()

	if len(rest) > 0 {
		in.prompt = strings.Join(rest, " ")
	}

	return in, nil
}
