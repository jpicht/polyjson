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
		PolyStructNewGen{},
		SliceAppendGen{},
		AcceptFuncGen{},
		MarshalFuncGen{},
		UnmarshalFuncGen{},
		DefaultVisitorGen{},
		VisitorFuncGen{},
	},
}
