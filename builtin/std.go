package builtin

var StdReg *Registry

func init() {
	builder := NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerIOModule(builder)
	registerStringModule(builder)

	StdReg = builder.Build()
}
