package codegen_test

import (
	"bytes"
	"testing"

	"github.com/jpicht/polyjson/codegen"
	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
	"github.com/stretchr/testify/require"
)

func TestParserToPolyStruct(t *testing.T) {
	c := parser.DefaultConfig
	r, err := c.Parse("../testdata")
	require.NoError(t, err)
	require.NotNil(t, r)

	buf := new(bytes.Buffer)

	for _, pkg := range r.Packages {
		poly, err := generator.PolyStructs(pkg)
		ctx := codegen.DefaultConfig.NewContext(pkg, buf)
		require.NoError(t, err)
		for _, s := range poly {
			t.Logf("%s: %v", s.Name, s.Impls)
			for _, g := range codegen.All {
				g.GeneratePolyStruct(ctx, s)
			}
		}
	}

	t.Logf("output:\n%s", buf.String())
}
