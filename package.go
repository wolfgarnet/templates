package templates

import (
	"html/template"
	"os"
	"log"
	"path/filepath"
	"strings"
)

type Package struct {
	objects map[string]*template.Template
	path string
	libPath string
}

func (p *Package) String() string {
	return "Package for " + p.path
}

func (p *Package) add(path string) {
	logger.Debug("Adding %v", path)

	dir, file := filepath.Split(path)
	name := strings.Trim(dir[len(p.path)+1:], "/\\")

	tpl, ok := p.objects[name]
	if !ok {
		logger.Debug("Creating new template, %v, for %v", file, dir)
		tpl = template.New(name)
		p.addLib(tpl)
		p.objects[name] = tpl
	}

	_, err := tpl.ParseFiles(path)
	if err != nil {
		log.Fatalf("Unable to parse template, %v", err)
	}
}

func (p *Package) addLib(tpl *template.Template) {
	glob := p.path + string(os.PathSeparator) + p.libPath + "/*.html"
	log.Printf("Adding lib: %v", glob)
	tpl.ParseGlob(glob)
}

func (p *Package) packageWalker(path string, f os.FileInfo, err error) error {
	log.Printf("Path: %v(%v)", path, p.path)
	l := len(p.path)
	l2 := len(p.libPath)

	if len(path) < l+l2 {
		return nil
	}

	sub := path[l+1:l+1+len(p.libPath)]
	//log.Printf("l: %v, l2: %v ===== %v", l, l2, sub)
	if sub == p.libPath {
		return nil
	}
	if f.Mode().IsRegular() && path[len(path)-1:] != "~" {
		p.add(path)
	}
	return nil
}

func (p *Package) Get(objectName string) *template.Template {
	object, ok := p.objects[objectName]
	if !ok {
		return nil
	}

	return object
}

func (p *Package) Contains(objectName, method string) bool {
	object := p.Get(objectName)
	if object == nil {
		return false
	}

	tpl := object.Lookup(method)

	return tpl != nil
}

func (p *Package) Print() {
	log.Printf("Printing %v", p.path)
	//printTemplate(p.objects, 2, 0)
	for i, j := range p.objects {
		//log.Printf("%v: %v", i, getTemplate(j.Name()))
		log.Printf("%v: [%v]", i, getTemplate(j))
	}
}

func getTemplate(tpl *template.Template) (s string) {
	for _, j := range tpl.Templates() {
		s += j.Name() + ", "
	}

	return s
}

func printTemplate(tpl *template.Template, max, level int) {
	if level > max {
		return
	}
	log.Printf("Printing template: %v @ %v", tpl.Name(), level)
	for i, j := range tpl.Templates() {
		log.Printf("TEMPLATE: %v=%v", i, j.Name())
	}

	for _, j := range tpl.Templates() {
		printTemplate(j, max, level + 1)
	}
}

func NewPackage(path string) *Package {
	logger.Debug("New package from %v", path)

	p := &Package{}
	p.objects = make(map[string]*template.Template)
	p.path = path
	p.libPath = "lib"

	filepath.Walk(path, p.packageWalker)

	return p
}

