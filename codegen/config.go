package codegen

type Config struct {
	Generators        []CodeGen
	OutputFileOptions []OutputFileOption
}

var DefaultConfig = Config{
	Generators: []CodeGen{
		VisitorInterfaceGen{},
		PolyStructGen{},
		SliceTypeGen{},
		AcceptFuncGen{},
		MarshalFuncGen{},
		UnmarshalFuncGen{},
		DefaultVisitorGen{},
		VisitorFuncGen{},
	},
}
