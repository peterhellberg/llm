# LLM 🔗

## What is this?

This is a _very_ small subset of the [Langchain](https://github.com/langchain-ai/langchain) project in [Go](https://go.dev/).

> [!Important]
> It is quite heavily based on the more ambitious [LangChainGo](https://github.com/tmc/langchaingo), for now you should likely use that instead.

My main objective with this module is to have declarations of some core interfaces (Such as `llm.Provider` and `llm.Chain`)
and other types that I can use as building blocks for LLM based Go applications. I do **not** intend for this module to
contain many _(or even any)_ implementations of these interfaces.

## Examples

See [./examples](./examples) for example usage.

## License (MIT)

Since **LangChainGo** is [licensed under MIT](https://github.com/tmc/langchaingo/blob/main/LICENSE)
it goes without saying that this project should be so as well :clipboard:
