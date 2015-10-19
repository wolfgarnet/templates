package templates

import (
	"bytes"
	"reflect"
	"html/template"
)

type Tools struct {
	renderer *Renderer
}

func (t Tools) Render(object interface{}, view string, trySuper bool) template.HTML {
	logger.Debug("Tool render {}, {}", object, view)
	renderer, err:= t.renderer.manager.RenderObject(t.renderer.ThemeName, t.renderer.PackageName, object, view, trySuper)

	if err != nil {
		return template.HTML(err.Error())
	}

	b := new(bytes.Buffer)
	renderer.Render(b)
	return template.HTML(b.String())
}

func (t Tools) RenderStatic(tt reflect.Type, view string) template.HTML {
	renderer, err := t.renderer.manager.RenderType(t.renderer.ThemeName, t.renderer.PackageName, tt, view)

	if err != nil {
		return template.HTML(err.Error())
	}

	b := new(bytes.Buffer)
	renderer.Render(b)
	return template.HTML(b.String())
}