package templates

import (
	"io"
	"html/template"
)

type Renderer struct {
	ThemeName, PackageName string
	Object interface{}
	Template *template.Template
	Title string
	manager *Manager
	data map[string]interface{}
}

func (r *Renderer) AddData(field string, data interface{}) {
	r.data[field] = data
}

func (r *Renderer) Render(writer io.Writer) error {
	//println("RUNNING TEMPLATE RUNNER", reflect.TypeOf(r.Object).Elem().Name())
	r.Template.Execute(writer, map[string]interface{} {
		"instance": r.Object,
		"tools": Tools{r},
		"Title": r.Title,
		"test": "test1",
		"data": r.data,
	})

	return nil
}