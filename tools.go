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
	logger.Debug("Tool render %v, %v - %v", object, view, trySuper)
	renderer, err:= t.renderer.manager.RenderObject(t.renderer.ThemeName, t.renderer.PackageName, object, view, trySuper)

	logger.Debug("-------> %v(%v)", renderer, err)

	if err != nil {
		return template.HTML(err.Error())
	}

	b := new(bytes.Buffer)
	renderer.Render(b)
	logger.Debug("I GOT: %v", b.String())
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