package templates
import "testing"

func TestDefault(t *testing.T) {
	path := "testdata/scenario1"

	p := NewPackage(path)
	p.Print()
}
