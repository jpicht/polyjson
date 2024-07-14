package polyjson

import (
	"github.com/jpicht/polyjson/markers"
)

type Interface[IF any] interface{ markers.Interface }
type Implements[IF any] struct {
	markers.Implements `json:"-"`
	Parent             *IF `json:"-"`
}
type Common[IF any] struct{ markers.Common }
