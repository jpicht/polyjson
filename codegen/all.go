package codegen

var All = []CodeGen{
	VisitorInterfaceGen{},
	PolyStructGen{},
	SliceTypeGen{},
	AcceptFuncGen{},
	MarshalFuncGen{},
	UnmarshalFuncGen{},
	DefaultVisitorGen{},
	VisitorFuncGen{},
}
