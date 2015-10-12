package templates

import (
	"log"
	"io/ioutil"
	"strings"
)

type Theme struct {
	packages map[string]*Package
	name string
}

func (t Theme) String() string {
	return "Theme " + t.name
}

func NewTheme(path string) *Theme {
	theme := &Theme{make(map[string]*Package), strings.Trim(path, "/\\")}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			log.Printf("DIR IS %v", f.Name())
			p := NewPackage(path + f.Name())
			p.Print()
			theme.packages[f.Name()] = p
		}
	}

	return theme
}