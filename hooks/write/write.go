package write

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.Hooks = Hooks{}

// Hooks that write to the embedded Writer.
type Hooks struct {
	io.Writer
}

func (h Hooks) ProviderGenerateContentStart(_ context.Context, ms []llm.Message) {
	fmt.Fprintln(h, "> Entering LLM with messages:")

	for _, m := range ms {
		var buf strings.Builder

		for _, t := range m.Parts {
			if t, ok := t.(llm.TextContent); ok {
				buf.WriteString(t.Text)
			}
		}

		fmt.Fprintln(h, "Role:", m.Role)
		fmt.Fprintln(h, "Text:", buf.String())
	}
}

func (h Hooks) ProviderGenerateContentEnd(_ context.Context, res *llm.ContentResponse) {
	fmt.Fprintln(h, "< Exiting LLM with response:")

	for _, c := range res.Choices {
		if c.Content != "" {
			fmt.Fprintln(h, "Content:", c.Content)
		}

		if c.StopReason != "" {
			fmt.Fprintln(h, "StopReason:", c.StopReason)
		}

		if len(c.GenerationInfo) > 0 {
			fmt.Fprintln(h, "GenerationInfo:")

			for k, v := range c.GenerationInfo {
				fmt.Fprintf(h, "%20s: %v\n", k, v)
			}
		}

		if c.FuncCall != nil {
			fmt.Fprintln(h, "FuncCall: ", c.FuncCall.Name, c.FuncCall.Arguments)
		}
	}
}

func (h Hooks) StreamingFunc(_ context.Context, chunk []byte) {
	fmt.Fprintln(h, string(chunk))
}

func (h Hooks) Text(_ context.Context, text string) {
	fmt.Fprintln(h, text)
}

func (h Hooks) ProviderStart(_ context.Context, prompts []string) {
	fmt.Fprintln(h, " 🚪 Entering LLM with prompts:", prompts)
}

func (h Hooks) ProviderError(_ context.Context, err error) {
	fmt.Fprintln(h, " 👋 Exiting LLM with error:", err)
}

func (h Hooks) ChainStart(_ context.Context, inputs map[string]any) {
	fmt.Fprintln(h, " 🚪 Entering chain with inputs:", formatChainValues(inputs))
}

func (h Hooks) ChainEnd(_ context.Context, outputs map[string]any) {
	fmt.Fprintln(h, " 👋 Exiting chain with outputs:", formatChainValues(outputs))
}

func (h Hooks) ChainError(_ context.Context, err error) {
	fmt.Fprintln(h, " 👋 Exiting chain with error:", err)
}

func (h Hooks) ToolStart(_ context.Context, input string) {
	fmt.Fprintln(h, " 🚪 Entering tool with input:", removeNewLines(input))
}

func (h Hooks) ToolEnd(_ context.Context, output string) {
	fmt.Fprintln(h, " 👋 Exiting tool with output:", removeNewLines(output))
}

func (h Hooks) ToolError(_ context.Context, err error) {
	fmt.Fprintln(h, " 👋 Exiting tool with error:", err)
}

func (h Hooks) AgentAction(_ context.Context, action llm.AgentAction) {
	fmt.Fprintln(h, " 🕵️ Agent selected action:", formatAgentAction(action))
}

func (h Hooks) AgentFinish(_ context.Context, finish llm.AgentFinish) {
	fmt.Fprintf(h, " 🕵️ Agent finish: %v \n", finish)
}

func (h Hooks) RetrieverStart(_ context.Context, query string) {
	fmt.Fprintln(h, " 🚪 Entering retriever with query:", removeNewLines(query))
}

func (h Hooks) RetrieverEnd(_ context.Context, query string, documents []llm.Document) {
	fmt.Fprintln(h, " 👋 Exiting retriever with documents for query:", documents, query)
}

func formatChainValues(values map[string]any) string {
	output := ""

	for key, value := range values {
		output += fmt.Sprintf("\"%s\" : \"%s\", ",
			removeNewLines(key),
			removeNewLines(value),
		)
	}

	return output
}

func formatAgentAction(action llm.AgentAction) string {
	return fmt.Sprintf("\"%s\" with input \"%s\"",
		removeNewLines(action.Tool),
		removeNewLines(action.ToolInput),
	)
}

func removeNewLines(s any) string {
	return strings.ReplaceAll(fmt.Sprint(s), "\n", " ")
}
