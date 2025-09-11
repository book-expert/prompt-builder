// ./internal/promptbuilder/builder.go
package promptbuilder

import (
	"bytes"
	"fmt"
	"text/template"
)

// Builder holds the state for a single prompt-building operation.
type Builder struct {
	template *template.Template
	vars     map[string]any
}

// New creates and initializes a new prompt builder with the given template string.
func New(templateStr string) (*Builder, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Builder{
		template: tmpl,
		vars:     make(map[string]any),
	}, nil
}

// Set adds or updates a single variable for the template.
func (b *Builder) Set(key string, value any) {
	b.vars[key] = value
}

// SetMap adds or updates multiple variables from a map.
func (b *Builder) SetMap(vars map[string]any) {
	for key, value := range vars {
		b.vars[key] = value
	}
}

// Build executes the template with the stored variables and returns the result.
func (b *Builder) Build() (string, error) {
	var buf bytes.Buffer
	if err := b.template.Execute(&buf, b.vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}
