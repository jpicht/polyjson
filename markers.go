package polyjson

type Interface[IF any] interface{ isInterface() }
type Implements[IF any] struct{}
type Common[IF any] struct{}

func (Implements[IF]) isImplementation() {}
func (Common[IF]) isCommon()             {}
