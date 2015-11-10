package templates

import (
	"reflect"
	"html/template"
	"path/filepath"
	"log"
	"fmt"
	"io/ioutil"

	"github.com/wolfgarnet/logging"
)

var logger logging.Logger

type Super interface {
	Super() interface{}
}

type Manager struct {
	themes map[string]*Theme
	extension string
	defaultThemeName string
}

func NewManager() *Manager {
	m := &Manager{
		themes: make(map[string]*Theme),
		extension: ".html",
		defaultThemeName: "default",
	}

	return m
}

func (m *Manager) AddTheme(theme *Theme) {
	m.themes[theme.name] = theme
}

func (m *Manager) Parse(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			log.Printf("Parsing %v from %v", f.Name(), path)
			t := NewTheme(path + f.Name())
			m.AddTheme(t)
		}
	}
}

// GetObjectTemplate retrieves a template given an anonymous object
func (m *Manager) GetObjectTemplate(themeName, packageName string, object interface{}, view string, trySuper bool) (*template.Template, error) {
	logger.Debug("Getting object template for {} with {}", object, view)

	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	logger.Debug("Getting template for %v of type %v", object, t)

	tpl, err := m.GetTypeTemplate(themeName, packageName, t, view)

	if err != nil {
		if trySuper {
			log.Printf("FAILED2, %v", err)
			superType := reflect.TypeOf((*Super)(nil)).Elem()
			if t.Implements(superType) {
				log.Printf("Supertype: %v", superType)
				sm, err2 := object.(Super)
				log.Printf("SM: %v -- %v", sm, err2)

				// Recursive call the GetObjectTemplate method
				return m.GetObjectTemplate(themeName, packageName, sm.Super(), view, trySuper)
			}
		}

		return nil, err
	}

	return tpl, nil
}

// GetObjectTemplate retrieves a template given a certain type
func (m *Manager) GetTypeTemplate(themeName, packageName string, t reflect.Type, view string) (*template.Template, error) {
	logger.Debug("Getting type template=%v+%v", t.PkgPath(), t.Name())

	name := filepath.Join(t.PkgPath(), t.Name())
	logger.Debug("Theme: %v, Package: %v, View: %v, method: %v, tname: %v, pkg: %v", themeName, packageName, name, view, t.Name(), t.PkgPath())
	return m.GetTemplate(themeName, packageName, name, view + m.extension)
}

// Already has extension!?
func (m *Manager) GetTemplate(themeName, packageName, object, name string) (*template.Template, error) {
	logger.Debug("Getting template: %v - %v - %v - %v", themeName, packageName, object, name)

	theme, err := m.getTemplate(themeName, packageName, object)
	if err != nil {
		return nil, err
	}

	if theme == nil {
		return nil, fmt.Errorf("%v: %v was not found", themeName, theme)
	}

	template := theme.Lookup(name)

	if template == nil {
		logger.Debug("No template found, trying default if not....")
		if themeName != m.defaultThemeName {
			// If not the default theme, try that!
			return m.GetTemplate(m.defaultThemeName, packageName, object, name)
		} else {
			// If the default theme, fail!
			return nil, fmt.Errorf("Template, " + name + ", does not exist")
		}
	} else {
		logger.Debug("Returning template: %v", template)
		return template, nil
	}


}

// getTemplate retrieves the template given the object name
func (m *Manager) getTemplate(themeName, packageName, object string) (*template.Template, error) {
	theme, ok1 := m.themes[themeName]
	if theme == nil {
		return nil, fmt.Errorf("Theme " + themeName + " does not exist")
	}
	logger.Debug("Theme exists:", themeName, ok1)

	pack, ok2 := theme.packages[packageName]
	if pack == nil {
		return nil, fmt.Errorf("Package " + packageName + " does not exist")
	}
	logger.Debug("Package exists:", packageName, ok2)

	obj, ok3 := pack.objects[object]
	if obj == nil {
		return nil, fmt.Errorf("Object " + object + " does not exist")
	}
	logger.Debug("Object exists:", object, ok3)

	return obj, nil
}

// RenderObject renders an object
func (m *Manager) RenderObject(themeName, packageName string, object interface{}, view string, trySuper bool) (*Renderer, error) {
	logger.Debug("Rendering object: {} with {}", object, view)

	t, err := m.GetObjectTemplate(themeName, packageName, object, view, trySuper)
	if err != nil {
		return nil, err
	}

	return &Renderer{themeName, packageName, object, t, "TEST", m, make(map[string]interface{})	}, nil
}

// RenderType renders given a type
func (m *Manager) RenderType(themeName, packageName string, t reflect.Type, view string) (*Renderer, error) {
	logger.Debug("Rendering type: {} with {}", t, view)

	tpl, err := m.GetTypeTemplate(themeName, packageName, t, view)
	if err != nil {
		return nil, err
	}

	return &Renderer{themeName, packageName, nil, tpl, "TEST", m, make(map[string]interface{})}, nil

}
