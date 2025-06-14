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

				"abs": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Abs(numArg.GetNumber()))
				}),
				"acos": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Acos(numArg.GetNumber()))
				}),
				"acosh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Acosh(numArg.GetNumber()))
				}),
				"asin": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Asin(numArg.GetNumber()))
				}),
				"asinh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Asinh(numArg.GetNumber()))
				}),
				"atan": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Atan(numArg.GetNumber()))
				}),
				"atanh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Atanh(numArg.GetNumber()))
				}),
				"cbrt": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Cbrt(numArg.GetNumber()))
				}),
				"ceil": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Ceil(numArg.GetNumber()))
				}),
				"cos": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Cos(numArg.GetNumber()))
				}),
				"cosh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Cosh(numArg.GetNumber()))
				}),
				"exp": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Exp(numArg.GetNumber()))
				}),
				"floor": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Floor(numArg.GetNumber()))
				}),
				"log": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Log(numArg.GetNumber()))
				}),
				"log10": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg
					}

					return vm.NewNumber(math.Log10(numArg.GetNumber()))
				}),
				"max": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2
					}

					return vm.NewNumber(math.Max(numArg1.GetNumber(), numArg2.GetNumber()))
				}),
				"min": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2
					}

					return vm.NewNumber(math.Min(numArg1.GetNumber(), numArg2.GetNumber()))
				}),
				"pow": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					numArg2, ok := args.Get(1, vm.ValueTypeNumber)
					if !ok {
						return numArg2
					}

					return vm.NewNumber(math.Pow(numArg1.GetNumber(), numArg2.GetNumber()))
				}),
				"round": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Round(numArg1.GetNumber()))
				}),
				"roundEven": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.RoundToEven(numArg1.GetNumber()))
				}),
				"sin": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Sin(numArg1.GetNumber()))
				}),
				"sinh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Sinh(numArg1.GetNumber()))
				}),
				"sqrt": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Sqrt(numArg1.GetNumber()))
				}),
				"tan": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Tan(numArg1.GetNumber()))
				}),
				"tanh": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
					numArg1, ok := args.Get(0, vm.ValueTypeNumber)
					if !ok {
						return numArg1
					}

					return vm.NewNumber(math.Tanh(numArg1.GetNumber()))
				}),
			},
		)
	})
}
