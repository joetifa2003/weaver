package builtin

var StdReg *Registry

func init() {
	builder := NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerIOModule(builder)

	StdReg = builder.Build()
}
