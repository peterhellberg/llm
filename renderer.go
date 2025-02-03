package llm

// Renderer is used to convert a template string and a map of values into a final string.
type Renderer interface {
	RenderTemplate(tmpl string, values map[string]any) (string, error)
}

// RendererFunc is a convenience function type that implements the Renderer interface by calling itself.
type RendererFunc func(string, map[string]any) (string, error)

func (render RendererFunc) RenderTemplate(tmpl string, values map[string]any) (string, error) {
	return render(tmpl, values)
}
