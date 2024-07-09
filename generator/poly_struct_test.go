package generator_test

import (
	"testing"

	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
	"github.com/jpicht/polyjson/parser/att"
	"github.com/stretchr/testify/require"
)

func TestParserToPolyStruct(t *testing.T) {
	c := parser.Config{
		Markers: []att.Marker{
			parser.MarkerCommon,
			parser.MarkerIs,
		},
	}
	r, err := c.Parse("../testdata")
	require.NoError(t, err)
	require.NotNil(t, r)

	for _, pkg := range r.Packages {
		poly, err := generator.PolyStructs(pkg)
		require.NoError(t, err)
		for _, s := range poly {
			t.Logf("%s: %v", s.Name, s.Impls)
		}
	}
}
