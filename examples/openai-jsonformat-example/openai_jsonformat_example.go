package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/openai"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	provider, err := openai.New(
		openai.WithModel("gpt-4o-mini"),
		openai.WithResponseFormat(responseFormat()),
	)
	if err != nil {
		return err
	}

	content := []llm.Message{
		llm.TextParts(llm.ChatMessageTypeSystem, "You are an expert at structured data responses in JSON. You should should answer with the given structure."),
		llm.TextParts(llm.ChatMessageTypeHuman, "Who is the current regent of Sweden, and how old are they this year (which is 2025)?"),
	}

	completion, err := provider.GenerateContent(ctx, content, llm.WithJSONMode())
	if err != nil {
		return err
	}

	data := []byte(completion.Choices[0].Content)

	var p Person

	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Person: %#v\n", p)

	return nil
}

type Person struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Age  int    `json:"age"`
}

func responseFormat() *openai.ResponseFormat {
	return &openai.ResponseFormat{
		Type: "json_schema",
		JSONSchema: &openai.ResponseFormatJSONSchema{
			Name: "object",
			Schema: &openai.ResponseFormatJSONSchemaProperty{
				Type: "object",
				Properties: map[string]*openai.ResponseFormatJSONSchemaProperty{
					"name": {
						Type:        "string",
						Description: "The name of the person, only first and last names",
					},
					"role": {
						Type:        "string",
						Description: "The role of the person",
					},
					"age": {
						Type:        "integer",
						Description: "The age of the person",
					},
				},
				AdditionalProperties: false,
				Required:             []string{"name", "age", "role"},
			},
			Strict: true,
		},
	}
}
