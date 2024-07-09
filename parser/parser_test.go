package parser

import (
	"testing"

	"github.com/jpicht/polyjson/parser/att"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestParser(t *testing.T) {
	c := Config{
		Markers: []att.Marker{
			MarkerCommon,
			MarkerIs,
		},
	}
	r, err := c.Parse("../testdata")
	require.NoError(t, err)
	require.NotNil(t, r)

	for _, pkg := range r.Packages {
		for _, s := range pkg.NamedStructs {
			t.Logf("%s: %s (%v)", s.File, s.Name, maps.Keys(s.Markers))
		}
	}
}
