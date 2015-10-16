package templates
import (
	"io"
	"reflect"
	"html/template"
)

type Renderer struct {
	Object interface{}
	Template *template.Template
	Title string
	manager *Manager
}

func (r Renderer) Render(writer io.Writer) error {
	println("RUNNING TEMPLATE RUNNER", reflect.TypeOf(r.Object).Elem().Name())
	r.Template.Execute(writer, map[string]interface{} {
		"instance": r.Object,
		"tools": Tools{},
		"Title": r.Title,
		"test": "test1",
	})

	return nil
}