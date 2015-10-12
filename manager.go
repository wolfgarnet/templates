package templates
import (
	"reflect"
	"html/template"
	"path/filepath"
	"log"
	"fmt"
)

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

/*
func (m *Manager) Parse(path string) {

}
*/

// GetObjectTemplate retrieves a template given an anonymous object
func (m *Manager) GetObjectTemplate(themeName, packageName string, object interface{}, method string) (*template.Template, error) {

	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	log.Printf("Getting template for %v of type %v", object, t)

	return m.GetTypeTemplate(themeName, packageName, t, object, method)
}

// GetObjectTemplate retrieves a template given a certain type
func (m *Manager) GetTypeTemplate(themeName, packageName string, t reflect.Type, object interface{}, method string) (*template.Template, error) {
	name := filepath.Join(t.PkgPath(), t.Name())
	println("OBJ, theme:", themeName, "platform:", packageName, "Name:", name, "method:", method, "tname:", t.Name(), "pkg.", t.PkgPath())
	tpl, err := m.GetTemplate(themeName, packageName, name, method + m.extension)

	if err != nil {
		superType := reflect.TypeOf((*Super)(nil)).Elem()
		if t.Implements(superType) {
			log.Printf("Supertype: %v", superType)
			sm, err2 := object.(Super)
			log.Printf("SM: %v -- %v", sm, err2)
			// Recursive call the GetObjectTemplate method
			return m.GetObjectTemplate(themeName, packageName, sm.Super(), method)
		}

		return nil, err
	}

	return tpl, nil
}

// Already has extension!?
func (m *Manager) GetTemplate(themeName, packageName, object, name string) (*template.Template, error) {
	println(themeName, packageName, object, name)

	theme, err := m.getTemplate(themeName, packageName, object)
	if err != nil {
		return nil, err
	}

	if theme == nil {
		return nil, fmt.Errorf("%v: %v was not found", themeName, theme)
	}

	template := theme.Lookup(name)

	if template == nil {
		if themeName != m.defaultThemeName {
			// If not the default theme, try that!
			return m.GetTemplate(m.defaultThemeName, packageName, object, name)
		} else {
			// If the default theme, fail!
			return nil, fmt.Errorf("Template, " + name + ", does not exist")
		}
	} else {
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
		return nil, fmt.Errorf("Platform " + packageName + " does not exist")
	}
	println("Platform exists:", packageName, ok2)

	obj, ok3 := pack.objects[object]
	if obj == nil {
		return nil, fmt.Errorf("Object " + object + " does not exist")
	}
	println("Object exists:", object, ok3)

	return obj, nil
}