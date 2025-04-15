package builtin

import "github.com/joetifa2003/weaver/registry"

var StdReg *registry.Registry

func init() {
	builder := registry.NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerBuiltinFuncsModules(builder)
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
