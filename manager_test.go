package templates

import (
	"testing"
	"os"
	"fmt"
)

type object struct {

}

func TestManagerCreation(t *testing.T) {
	println("TEST")
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	println("TEST2")
	/*
	o1 := object{}

	m.GetObjectTemplate("sce")
	*/
}

func TestManagerGet(t *testing.T) {
	println("TEST")
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	println("TEST2")

	o1 := object{}

	m.GetObjectTemplate("theme1", "scenario1", o1, "index")
}


func TestManagerPrint(t *testing.T) {
	println("TEST")
	m := NewManager()
	theme := NewTheme("testdata/theme1/")
	m.AddTheme(theme)

	println("TEST2")

	o1 := object{}

	tpl, err := m.GetObjectTemplate("theme1", "package1", o1, "index")
	if err != nil {
		fmt.Errorf("FAILED, %v", err)
	}
	println("awdawdad", tpl)
	tpl.Execute(os.Stdout, o1)
}
