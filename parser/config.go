package parser

import (
	"github.com/jpicht/polyjson/parser/att"
	"golang.org/x/tools/go/packages"
)

type Config struct {
	Interfaces []att.Marker
	Markers    []att.Marker
	Packages   packages.Config
}

var DefaultConfig = Config{
	Markers: []att.Marker{
		MarkerCommon,
		MarkerIs,
	},
	Interfaces: []att.Marker{
		EasyJSONUnmarshaler,
		EncodingJSONUnmarshaler,
		AdditionalFields,
		JWriterWritable,
	},
	Packages: packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedDeps,
	},
}