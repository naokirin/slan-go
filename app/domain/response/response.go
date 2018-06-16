package response

import (
	"bytes"
	"log"
	"text/template"
)

var _ TemplateInterface = (*Template)(nil)

// TemplateInterface is interface for response message templates
type TemplateInterface interface {
	AddTemplates(templates map[string]string)
	GetText(key string, templateArgs map[string]string) string
}

// Template for response message templates
type Template struct {
	tpls map[string]*template.Template
}

// AddTemplates add template map
func (t *Template) AddTemplates(templates map[string]string) {
	if t.tpls == nil {
		t.tpls = make(map[string]*template.Template)
	}
	for key := range templates {
		tp, err := template.New(key).Parse(templates[key])
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		t.tpls[key] = tp
	}
}

// GetText returns corresponding text by template
func (t *Template) GetText(key string, templateArgs map[string]string) string {
	if templateArgs == nil {
		templateArgs = make(map[string]string)
	}
	var text bytes.Buffer
	t.tpls[key].Execute(&text, templateArgs)
	return text.String()
}
