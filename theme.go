package templates

import (
	"log"
	"io/ioutil"
	"path/filepath"
)

type Theme struct {
	packages map[string]*Package
	name string
}

func (t Theme) String() string {
	return t.name
}

func NewTheme(path string) *Theme {
	name := filepath.Base(path)

	theme := &Theme{make(map[string]*Package), name}
	log.Printf("New theme: %v", theme)

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