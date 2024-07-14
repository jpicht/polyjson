package parser

import (
	_ "encoding/json"
	"go/types"
	"log"

	_ "github.com/mailru/easyjson"
	"golang.org/x/tools/go/packages"

	_ "github.com/jpicht/polyjson/markers"
	"github.com/jpicht/polyjson/parser/att"
)

var (
	MarkerImplements        = att.Mark("github.com/jpicht/polyjson/markers", "Implements")
	MarkerCommon            = att.Mark("github.com/jpicht/polyjson/markers", "Common")
	MarkerInterface         = att.Mark("github.com/jpicht/polyjson/markers", "Interface")
	MarkerTypeID            = att.Mark("github.com/jpicht/polyjson/markers", "TypeID")
	AdditionalFields        = att.Mark("github.com/jpicht/polyjson/markers", "AdditionalFields")
	JWriterWritable         = att.Mark("github.com/launchdarkly/go-jsonstream/v3", "Writable")
	EasyJSONUnmarshaler     = att.Mark("github.com/mailru/easyjson", "Unmarshaler")
	EncodingJSONUnmarshaler = att.Mark("encoding/json", "Unmarshaler")
	EncodingJSONMarshaler   = att.Mark("encoding/json", "Marshaler")
)

func resolveMarkers(c *Config, markers []att.Marker, packages []*packages.Package) (map[att.Marker]func(types.Type) bool, bool) {
	all := true

	out := make(map[att.Marker]func(types.Type) bool, len(markers))
	for _, marker := range markers {
		fn, ok := resolveMarker(c, marker, packages)
		if ok {
			out[marker] = fn
			continue
		}
		all = false
	}

	return out, all
}

func resolveMarker(c *Config, marker att.Marker, pkgs []*packages.Package) (func(types.Type) bool, bool) {
	for _, inPkg := range pkgs {
		pkg, ok := inPkg.Imports[marker.Package]
		if !ok {
			continue
		}
		markerInterface := pkg.Types.Scope().Lookup(marker.InterfaceName)
		if markerInterface != nil {
			return makeMarkerCheck(markerInterface.Type()), true
		}
	}

	var fn func(types.Type) bool
	found := false

	packages.Visit(pkgs, func(p *packages.Package) bool {
		if p.ID == marker.Package {
			if c.Verbose {
				log.Printf("found package for marker %s.%s", marker.Package, marker.InterfaceName)
			}
			iface := p.Types.Scope().Lookup(marker.InterfaceName)
			if iface != nil {
				if c.Verbose {
					log.Printf("  found marker interface %s.%s", marker.Package, marker.InterfaceName)
				}
				found = true
				fn = makeMarkerCheck(iface.Type())
			}
		}
		return !found
	}, nil)

	if !found {
		log.Printf("warning: marker %s.%s could not be found.", marker.Package, marker.InterfaceName)
	}

	return fn, found
}
