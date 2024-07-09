package codegen

import (
	"go/types"
	"io"
	"log"
	"strings"

	"github.com/jpicht/polyjson/parser/att"
)

type Context struct {
	Imports map[string]string
	Package *att.Package
	io.Writer
}

func NewContext(p *att.Package, w io.Writer) *Context {
	return &Context{
		Package: p,
		Imports: map[string]string{},
		Writer:  w,
	}
}

func (c *Context) New(t types.Type) string {
	name, k := c.typeString(t)
	log.Printf("New(%v)", t)
	switch {
	case k&ptr == ptr:
		return "new(" + name + ")"
	case k&_map == _map:
		return "make(" + name + ")"
	}
	return name
}

func (c *Context) TypeString(t types.Type) string {
	name, k := c.typeString(t)
	return ptrPrefix[k&ptr] + name
}

type kind int

const (
	invalid kind = -1
	basic   kind = 1 << iota
	named
	ptr
	_map
)

func (c *Context) typeString(t types.Type) (string, kind) {
	//	qual := types.RelativeTo(c.Package.Types)
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name(), basic

	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()

		name := obj.Name()
		_, k := c.typeString(tt.Underlying())

		if pkg.Path() != c.Package.PkgPath {
			importedAs, ok := c.Imports[pkg.Path()]
			if !ok {
				importedAs = pkg.Name()
				for imported := range c.Imports {
					if imported == importedAs {
						importedAs = strings.ReplaceAll(pkg.Path(), "/", "_")
					}
				}
				c.Imports[pkg.Path()] = importedAs
			}
			log.Printf("import %s %q // for %s", importedAs, pkg.Path(), name)
			name = importedAs + "." + name
		}

		ta := tt.TypeArgs()
		if ta != nil && ta.Len() > 0 {
			name += "["
			first := true
			for i := 0; i < ta.Len(); i++ {
				arg := ta.At(i)
				e, k := c.typeString(arg)
				name += ptrPrefix[k&ptr] + e
				if first {
					first = false
					continue
				}
				name += ", "
			}
			name += "]"
		}

		return name, k | named

	case *types.Pointer:
		t, k := c.typeString(tt.Elem())
		return t, ptr | k

	case *types.Map:
		key, keyKind := c.typeString(tt.Key())
		val, valKind := c.typeString(tt.Elem())
		return "map[" + ptrPrefix[keyKind&ptr] + key + "]" + ptrPrefix[valKind&ptr] + val, _map

	case *types.Slice:
		t, _ := c.typeString(tt.Elem())
		return "[]" + t, basic
	default:
		log.Fatalf("unhandled type %T", t)
	}

	return "", invalid
}

var ptrPrefix = map[kind]string{
	ptr: "*",
}

func init() {
	log.Default().SetFlags(log.Flags() | log.Lshortfile)
}
