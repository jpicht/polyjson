package codegen

type Config struct {
	Generators        []CodeGen
	OutputFileOptions []OutputFileOption
}

var DefaultConfig = Config{
	Generators: All,
}
