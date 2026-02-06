package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// Renderer is an implementation of the usecase.MarkdownRenderer interface.
type Renderer struct {
	md goldmark.Markdown
}

// NewRenderer creates a new Markdown renderer based on goldmark.
func NewRenderer() *Renderer {
	return &Renderer{
		md: goldmark.New(),
	}
}

// RenderToHTML converts markdown into HTML.
func (r *Renderer) RenderToHTML(input string) (string, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(input), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

