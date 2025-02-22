package builtin

var StdReg *Registry

func init() {
	builder := NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerIOModule(builder)
	registerStringModule(builder)
	registerJSONModule(builder)
	registerModuleRL(builder)

	StdReg = builder.Build()
}
