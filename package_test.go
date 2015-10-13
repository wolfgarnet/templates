package templates

import "testing"

func TestDefault(t *testing.T) {
	path := "testdata/theme1/package1/"

	p := NewPackage(path)
	p.Print()
}
