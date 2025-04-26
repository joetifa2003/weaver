package builtin

import (
	"os"
	"path/filepath"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/internal/pkg/ds"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncsModules(builder *vm.RegistryBuilder) {
	moduleCache := ds.NewConcMap[string, vm.Value]()

	builder.RegisterFunc("import", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		pathArg, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return pathArg
		}

		builtinModule, ok := v.Executor.Reg.ResolveModule(pathArg.GetString())
		if ok {
			return builtinModule
		}

		return initModule(v, pathArg.GetString(), moduleCache)
	})
}

func initModule(v *vm.VM, path string, moduleCache *ds.ConcMap[string, vm.Value]) vm.Value {
	absPath := filepath.Join(filepath.Dir(v.CurrentFrame().Path), path)
	if v, ok := moduleCache.Get(absPath); ok {
		return v
	}

	srcData, err := os.ReadFile(absPath)
	if err != nil {
		return vm.NewErrFromErr(err)
	}
	src := string(srcData)

	p, err := parser.Parse(src)
	if err != nil {
		return vm.NewErrFromErr(err)
	}

	irc := ir.NewCompiler()
	ircr, err := irc.Compile(absPath, p)
	if err != nil {
		return vm.NewErrFromErr(err)
	}

	c := compiler.New(StdReg)
	instructions, vars, constants, err := c.Compile(ircr)
	if err != nil {
		return vm.NewErrFromErr(err)
	}

	val := v.Run(vm.Frame{
		NumVars:      vars,
		Instructions: instructions,
		Constants:    constants,
		Path:         absPath,
		HaltAfter:    true,
	}, 0)

	moduleCache.Set(absPath, val)

	return val
}
