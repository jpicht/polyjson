package codegen

type Config struct {
	Generators        []CodeGen
	OutputFileOptions []OutputFileOption
}

var DefaultConfig = Config{
	Generators: []CodeGen{
		TypeFieldGen{},
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
