package templates

import (
	"testing"
	"os"
	"fmt"
	"bytes"
)

type object struct {

}

func TestManagerCreation(t *testing.T) {
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)
}

func TestManagerGet(t *testing.T) {
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	o1 := object{}

	m.GetObjectTemplate("theme1", "scenario1", o1, "index")
}

func TestManagerPrint(t *testing.T) {
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	o1 := object{}

	tpl, err := m.GetObjectTemplate("theme1", "package1", o1, "index")
	if err != nil {
		fmt.Errorf("FAILED, %v", err)
	}

	tpl.Execute(os.Stdout, o1)
}

func TestManagerDefault(t *testing.T) {
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	theme2 := NewTheme("testdata/theme2/")
	m.AddTheme(theme2)

	m.defaultThemeName = "theme1"

	o1 := object{}

	tpl, err := m.GetObjectTemplate("theme2", "package1", o1, "config")
	if err != nil {
		fmt.Errorf("FAILED, %v", err)
	}

	var buffer bytes.Buffer
	tpl.Execute(&buffer, o1)

	if buffer.String() != "CONFIG" {
		t.Errorf("Failed")
	}
}