package templates

import "bytes"

type Tools struct {
	renderer *Renderer
}

func (t Tools) Render(object interface{}, view string, trySuper bool) string {
	logger.Debug("Tool render {}, {}", object, view)
	renderer, err:= t.renderer.manager.RenderObject(t.renderer.ThemeName, t.renderer.PackageName, object, view, trySuper)

	if err != nil {
		return err.Error()
	}

	b := new(bytes.Buffer)
	renderer.Render(b)
	return b.String()
}

func (t Tools) RenderStatic() string {
	return "HEJ"
}