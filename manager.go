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

	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	logger.("Getting template for %v of type %v", object, t)

	return m.GetTypeTemplate(themeName, packageName, t, object, view, trySuper)
}

// GetObjectTemplate retrieves a template given a certain type
func (m *Manager) GetTypeTemplate(themeName, packageName string, t reflect.Type, object interface{}, view string, trySuper bool) (*template.Template, error) {
	logger.Debug("NAME=%v+%v", t.PkgPath(), t.Name())
	name := filepath.Join(t.PkgPath(), t.Name())
	log.Printf("Theme: %v, Package: %v, View: %v, method: %v, tname: %v, pkg: ", themeName, packageName, name, view, t.Name(), t.PkgPath())
	tpl, err := m.GetTemplate(themeName, packageName, name, view + m.extension)

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

// Already has extension!?
func (m *Manager) GetTemplate(themeName, packageName, object, name string) (*template.Template, error) {
	log.Printf("%v - %v - %v - %v", themeName, packageName, object, name)

	theme, err := m.getTemplate(themeName, packageName, object)
	if err != nil {
		return nil, err
	}

	if theme == nil {
		return nil, fmt.Errorf("%v: %v was not found", themeName, theme)
	}

	template := theme.Lookup(name)

	if template == nil {
		log.Printf("NAY")
		if themeName != m.defaultThemeName {
			// If not the default theme, try that!
			return m.GetTemplate(m.defaultThemeName, packageName, object, name)
		} else {
			// If the default theme, fail!
			return nil, fmt.Errorf("Template, " + name + ", does not exist")
		}
	} else {
		log.Printf("WWATATAT")
		return template, nil
	}


}

// getTemplate retrieves the template given the object name
func (m *Manager) getTemplate(themeName, packageName, object string) (*template.Template, error) {
	theme, ok1 := m.themes[themeName]
	if theme == nil {
		return nil, fmt.Errorf("Theme " + themeName + " does not exist")
	}
	println("Theme exists:", themeName, ok1)

	pack, ok2 := theme.packages[packageName]
	if pack == nil {
		return nil, fmt.Errorf("Package " + packageName + " does not exist")
	}
	println("Package exists:", packageName, ok2)

	obj, ok3 := pack.objects[object]
	if obj == nil {
		return nil, fmt.Errorf("Object " + object + " does not exist")
	}
	println("Object exists:", object, ok3)

	return obj, nil
}

func (m *Manager) RenderObject(themeName, packageName string, object interface{}, view string, trySuper bool) (*Renderer, error) {
	logger.Debug("Rendering object: {} with {}", object, view)

	t, err := m.GetObjectTemplate(themeName, packageName, object, view, trySuper)
	if err != nil {
		return nil, err
	}

	return &Renderer{object, t, "TEST"}, nil

}