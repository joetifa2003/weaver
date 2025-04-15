package builtin

import (
	"strings"

	"github.com/joetifa2003/weaver/internal/pkg/helpers"
	"github.com/joetifa2003/weaver/registry"
	"github.com/joetifa2003/weaver/vm"
)

func registerStringModule(builder *registry.RegistryBuilder) {
	builder.RegisterModule("strings", map[string]vm.Value{
		"concat": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			var res string
			for _, arg := range args {
				if arg.IsError() {
					return arg
				}
				res += arg.String()
			}

			return vm.NewString(res)
		}),
		"split": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}

			sepArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return sepArg
			}

			str := strArg.GetString()
			sep := sepArg.GetString()

			parts := helpers.SliceMap(strings.Split(str, sep), func(s string) vm.Value {
				return vm.NewString(s)
			})

			return vm.NewArray(parts)
		}),
		"lower": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			return vm.NewString(strings.ToLower(strArg.GetString()))
		}),
		"upper": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			return vm.NewString(strings.ToUpper(strArg.GetString()))
		}),
		"trim": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}

			str := strArg.GetString()
			return vm.NewString(strings.TrimSpace(str))
		}),
		"contains": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			substrArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return substrArg
			}
			return vm.NewBoolean(strings.Contains(strArg.GetString(), substrArg.GetString()))
		}),
		"startsWith": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			prefixArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return prefixArg
			}
			return vm.NewBoolean(strings.HasPrefix(strArg.GetString(), prefixArg.GetString()))
		}),
		"endsWith": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			suffixArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return suffixArg
			}
			return vm.NewBoolean(strings.HasSuffix(strArg.GetString(), suffixArg.GetString()))
		}),
		"replace": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			oldArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return oldArg
			}
			newArg, ok := args.Get(2, vm.ValueTypeString)
			if !ok {
				return newArg
			}

			count := -1 // Replace all by default
			if len(args) > 3 {
				countArg, ok := args.Get(3, vm.ValueTypeNumber)
				if !ok {
					return countArg // Error if type is wrong
				}
				count = int(countArg.GetNumber())
			}

			return vm.NewString(strings.Replace(strArg.GetString(), oldArg.GetString(), newArg.GetString(), count))
		}),
		"substring": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			startArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return startArg
			}

			str := strArg.GetString()
			start := int(startArg.GetNumber())
			end := len(str) // Default to end of string

			if len(args) > 2 {
				endArg, ok := args.Get(2, vm.ValueTypeNumber)
				if !ok {
					return endArg // Error if type is wrong
				}
				end = int(endArg.GetNumber())
			}

			// Basic bounds checking
			if start < 0 {
				start = 0
			}
			if end > len(str) {
				end = len(str)
			}
			if start > end || start >= len(str) {
				return vm.NewString("") // Return empty string for invalid ranges
			}

			return vm.NewString(str[start:end])
		}),
		"indexOf": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			substrArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return substrArg
			}

			return vm.NewNumber(float64(strings.Index(strArg.GetString(), substrArg.GetString())))
		}),
		"lastIndexOf": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			substrArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return substrArg
			}

			return vm.NewNumber(float64(strings.LastIndex(strArg.GetString(), substrArg.GetString())))
		}),
		"padStart": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			lengthArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return lengthArg
			}
			padStrArg := vm.NewString(" ") // Default pad string
			if len(args) > 2 {
				padStrArg, ok = args.Get(2, vm.ValueTypeString)
				if !ok {
					return padStrArg
				}
			}

			str := strArg.GetString()
			targetLength := int(lengthArg.GetNumber())
			padStr := padStrArg.GetString()

			return vm.NewString(padString(str, targetLength, padStr, true))
		}),
		"padEnd": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			strArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return strArg
			}
			lengthArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return lengthArg
			}
			padStrArg := vm.NewString(" ") // Default pad string
			if len(args) > 2 {
				padStrArg, ok = args.Get(2, vm.ValueTypeString)
				if !ok {
					return padStrArg
				}
			}

			str := strArg.GetString()
			targetLength := int(lengthArg.GetNumber())
			padStr := padStrArg.GetString()

			return vm.NewString(padString(str, targetLength, padStr, false))
		}),
	})
}

// Helper function for padding
func padString(str string, targetLength int, padString string, padStart bool) string {
	if len(str) >= targetLength {
		return str
	}

	if padString == "" {
		padString = " " // Default pad string is space
	}

	padLen := targetLength - len(str)
	repeatCount := padLen / len(padString)
	remainingPad := padLen % len(padString)

	padding := strings.Repeat(padString, repeatCount) + padString[:remainingPad]

	if padStart {
		return padding + str
	}
	return str + padding
}
