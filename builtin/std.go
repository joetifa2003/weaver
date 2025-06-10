package builtin

import (
	"github.com/joetifa2003/weaver/vm"
)

var StdReg *vm.Registry

func init() {
	builder := vm.NewRegBuilder()

	registerBuiltinFuncs(builder)
	registerBuiltinFuncsModules(builder)
	registerBuiltinFuncsArr(builder)

	registerIOModule(builder)
	registerStringModule(builder)
	registerJSONModule(builder)
	registerMathModule(builder)
	registerHTTPModule(builder)
	registerFiberModule(builder)
	registerTimeModule(builder)
	registerModuleRL(builder)
	registerHtmlModule(builder)

	StdReg = builder.Build()
}
