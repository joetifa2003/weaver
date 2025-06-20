package builtin

import (
	"math"

	"github.com/joetifa2003/weaver/vm"
)

func registerMathModule(builder *vm.RegistryBuilder) {
	builder.RegisterModule("math", func() vm.Value {
		return vm.NewObject(
			map[string]vm.Value{
				"e":       vm.NewNumber(math.E),
				"pi":      vm.NewNumber(math.Pi),
				"phi":     vm.NewNumber(math.Phi),
				"sqrt2":   vm.NewNumber(math.Sqrt2),
				"sqrte":   vm.NewNumber(math.SqrtE),
				"sqrtpi":  vm.NewNumber(math.SqrtPi),
				"sqrtphi": vm.NewNumber(math.SqrtPhi),
				"ln2":     vm.NewNumber(math.Ln2),
				"log2e":   vm.NewNumber(math.Log2E),
				"ln10":    vm.NewNumber(math.Ln10),
				"log10e":  vm.NewNumber(math.Log10E),
				"inf":     vm.NewNumber(math.Inf(1)),
				"ninf":    vm.NewNumber(math.Inf(-1)),
				"nan":     vm.NewNumber(math.NaN()),

				"abs": vm.NewNativeFunction("abs", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Abs(numArg.GetNumber())), true
				}),
				"acos": vm.NewNativeFunction("acos", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Acos(numArg.GetNumber())), true
				}),
				"acosh": vm.NewNativeFunction("acosh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Acosh(numArg.GetNumber())), true
				}),
				"asin": vm.NewNativeFunction("asin", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Asin(numArg.GetNumber())), true
				}),
				"asinh": vm.NewNativeFunction("asinh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Asinh(numArg.GetNumber())), true
				}),
				"atan": vm.NewNativeFunction("atan", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Atan(numArg.GetNumber())), true
				}),
				"atanh": vm.NewNativeFunction("atanh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Atanh(numArg.GetNumber())), true
				}),
				"cbrt": vm.NewNativeFunction("cbrt", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Cbrt(numArg.GetNumber())), true
				}),
				"ceil": vm.NewNativeFunction("ceil", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Ceil(numArg.GetNumber())), true
				}),
				"cos": vm.NewNativeFunction("cos", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Cos(numArg.GetNumber())), true
				}),
				"cosh": vm.NewNativeFunction("cosh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Cosh(numArg.GetNumber())), true
				}),
				"exp": vm.NewNativeFunction("exp", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Exp(numArg.GetNumber())), true
				}),
				"floor": vm.NewNativeFunction("floor", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Floor(numArg.GetNumber())), true
				}),
				"log": vm.NewNativeFunction("log", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Log(numArg.GetNumber())), true
				}),
				"log10": vm.NewNativeFunction("log10", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg, false
					}

					return vm.NewNumber(math.Log10(numArg.GetNumber())), true
				}),
				"max": vm.NewNativeFunction("max", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2, false
					}

					return vm.NewNumber(math.Max(numArg1.GetNumber(), numArg2.GetNumber())), true
				}),
				"min": vm.NewNativeFunction("min", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2, false
					}

					return vm.NewNumber(math.Min(numArg1.GetNumber(), numArg2.GetNumber())), true
				}),
				"pow": vm.NewNativeFunction("pow", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2, false
					}

					return vm.NewNumber(math.Pow(numArg1.GetNumber(), numArg2.GetNumber())), true
				}),
				"round": vm.NewNativeFunction("round", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Round(numArg1.GetNumber())), true
				}),
				"roundEven": vm.NewNativeFunction("roundEven", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.RoundToEven(numArg1.GetNumber())), true
				}),
				"sin": vm.NewNativeFunction("sin", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Sin(numArg1.GetNumber())), true
				}),
				"sinh": vm.NewNativeFunction("sinh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Sinh(numArg1.GetNumber())), true
				}),
				"sqrt": vm.NewNativeFunction("sqrt", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Sqrt(numArg1.GetNumber())), true
				}),
				"tan": vm.NewNativeFunction("tan", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Tan(numArg1.GetNumber())), true
				}),
				"tanh": vm.NewNativeFunction("tanh", func(v *vm.VM, args vm.NativeFunctionArgs) (vm.Value, bool) {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1, false
					}

					return vm.NewNumber(math.Tanh(numArg1.GetNumber())), true
				}),
			},
		)
	})
}
