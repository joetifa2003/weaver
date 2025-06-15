package builtin

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func registerBuiltinFuncsModules(builder *vm.RegistryBuilder) {
	moduleCache := map[string]vm.Value{}
	lock := sync.Mutex{}

	builder.RegisterFunc("import", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
		lock.Lock()
		defer lock.Unlock()

		pathArg, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return pathArg, false
		}

		pathStr := pathArg.GetString()

		if mod, ok := moduleCache[pathStr]; ok {
			return mod, true
		}

		modInit, ok := v.Executor.Reg.ResolveModule(pathStr)
		if ok {
			mod := modInit()
			moduleCache[pathStr] = mod
			return mod, true
		}

		userMod, ok := initModule(v, pathStr)
		if !ok {
			return userMod, false
		}

		moduleCache[pathStr] = userMod
		return userMod, true
	})
}

func initModule(v *vm.VM, path string) (vm.Value, bool) {
	absPath := filepath.Join(filepath.Dir(v.CurrentFrame().Path), path)

	srcData, err := os.ReadFile(absPath)
	if err != nil {
		return vm.NewErrFromErr(err), false
	}
	src := string(srcData)

	p, err := parser.Parse(src)
	if err != nil {
		return vm.NewErrFromErr(err), false
	}

	irc := ir.NewCompiler()
	ircr, err := irc.Compile(absPath, p)
	if err != nil {
		return vm.NewErrFromErr(err), false
	}

	c := compiler.New(StdReg)
	instructions, vars, constants, err := c.Compile(ircr)
	if err != nil {
		return vm.NewErrFromErr(err), false
	}

	return v.RunFunction(vm.NewFunction(
		vm.FunctionValue{
			NumVars:      vars,
			Instructions: instructions,
			Constants:    constants,
			Path:         absPath,
		},
	))
}
