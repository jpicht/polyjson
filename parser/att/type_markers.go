package att

import "go/types"

type TypeMarker struct {
	Marker Marker
	Name   string
	Target string
	Type   types.Type
}

func (tm TypeMarker) String() string {
	return tm.Name
}
