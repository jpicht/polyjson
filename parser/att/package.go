package att

import (
	"golang.org/x/tools/go/packages"
)

type Package struct {
	*packages.Package
	NamedStructs []*NamedStruct
	TypeMarkers  []TypeMarker
}
