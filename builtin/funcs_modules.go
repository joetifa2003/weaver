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

	builder.RegisterFunc("import", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
		lock.Lock()
		defer lock.Unlock()

		pathArg, ok := args.Get(0, vm.ValueTypeString)
		if !ok {
			return pathArg
		}

		pathStr := pathArg.GetString()

		if mod, ok := moduleCache[pathStr]; ok {
			return mod
		}

		builtinModule, ok := v.Executor.Reg.ResolveModule(pathStr)
		if ok {
			moduleCache[pathStr] = builtinModule
			return builtinModule
		}

		userMod := initModule(v, pathStr)
		moduleCache[pathStr] = userMod
		return userMod
	})
}

func initModule(v *vm.VM, path string) vm.Value {
	absPath := filepath.Join(filepath.Dir(v.CurrentFrame().Path), path)

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

	return val
}
