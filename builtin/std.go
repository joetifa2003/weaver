package builtin

var StdReg *Registry

func init() {
	builder := NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerBuiltinFuncsArr(builder)

	registerIOModule(builder)
	registerStringModule(builder)
	registerJSONModule(builder)
	registerHTTPModule(builder)
	registerFiberModule(builder)
	registerTimeModule(builder)

	registerModuleRL(builder)

	StdReg = builder.Build()
}
