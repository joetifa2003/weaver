package builtin

import (
	"time"

	"github.com/joetifa2003/weaver/vm"
)

func registerTimeModule(builder *vm.RegistryBuilder) {
	m := map[string]vm.Value{
		// --- Time Constants ---
		"ansic":       vm.NewString(time.ANSIC),
		"unixDate":    vm.NewString(time.UnixDate),
		"rubyDate":    vm.NewString(time.RubyDate),
		"rfc822":      vm.NewString(time.RFC822),
		"rfc822z":     vm.NewString(time.RFC822Z),
		"rfc850":      vm.NewString(time.RFC850),
		"rfc1123":     vm.NewString(time.RFC1123),
		"rfc1123z":    vm.NewString(time.RFC1123Z),
		"rfc3339":     vm.NewString(time.RFC3339),
		"rfc3339Nano": vm.NewString(time.RFC3339Nano),
		"kitchen":     vm.NewString(time.Kitchen),
		"stamp":       vm.NewString(time.Stamp),
		"stampMilli":  vm.NewString(time.StampMilli),
		"stampMicro":  vm.NewString(time.StampMicro),
		"stampNano":   vm.NewString(time.StampNano),
		"dateTime":    vm.NewString(time.DateTime),
		"dateOnly":    vm.NewString(time.DateOnly),
		"timeOnly":    vm.NewString(time.TimeOnly),

		// --- Duration Constants ---
		"nanosecond":  vm.NewNumber(float64(time.Nanosecond)),
		"microsecond": vm.NewNumber(float64(time.Microsecond)),
		"millisecond": vm.NewNumber(float64(time.Millisecond)),
		"second":      vm.NewNumber(float64(time.Second)),
		"minute":      vm.NewNumber(float64(time.Minute)),
		"hour":        vm.NewNumber(float64(time.Hour)),

		// --- Time Functions ---
		"now": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			return vm.NewTime(time.Now())
		}),

		"unix": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			secArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return secArg
			}
			nsecArg := vm.NewNumber(0) // Default nsec to 0
			if len(args) > 1 {
				nsecArg, ok = args.Get(1, vm.ValueTypeNumber)
				if !ok {
					return nsecArg
				}
			}
			return vm.NewTime(time.Unix(int64(secArg.GetNumber()), int64(nsecArg.GetNumber())))
		}),

		"unixMilli": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			milliArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return milliArg
			}
			return vm.NewTime(time.UnixMilli(int64(milliArg.GetNumber())))
		}),

		"unixMicro": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			microArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return microArg
			}
			return vm.NewTime(time.UnixMicro(int64(microArg.GetNumber())))
		}),

		"parse": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			layoutArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return layoutArg
			}
			valueArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return valueArg
			}

			t, err := time.Parse(layoutArg.GetString(), valueArg.GetString())
			if err != nil {
				return vm.NewErrFromErr(err)
			}
			return vm.NewTime(t)
		}),

		"since": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(time.Since(timeArg.GetTime())))
		}),

		"until": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(time.Until(timeArg.GetTime())))
		}),

		// --- Duration Functions ---
		"parseDuration": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationStrArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return durationStrArg
			}
			d, err := time.ParseDuration(durationStrArg.GetString())
			if err != nil {
				return vm.NewErrFromErr(err)
			}
			return vm.NewNumber(float64(d))
		}),

		// --- Time Methods (as functions taking time as first arg) ---
		"add": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			durationArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			t := timeArg.GetTime()
			d := time.Duration(durationArg.GetNumber())
			return vm.NewTime(t.Add(d))
		}),

		"sub": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			time1Arg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return time1Arg
			}
			time2Arg, ok := args.Get(1, vm.ValueTypeTime)
			if !ok {
				return time2Arg
			}
			t1 := time1Arg.GetTime()
			t2 := time2Arg.GetTime()
			return vm.NewNumber(float64(t1.Sub(t2)))
		}),

		"addDate": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			yearsArg, ok := args.Get(1, vm.ValueTypeNumber)
			if !ok {
				return yearsArg
			}
			monthsArg, ok := args.Get(2, vm.ValueTypeNumber)
			if !ok {
				return monthsArg
			}
			daysArg, ok := args.Get(3, vm.ValueTypeNumber)
			if !ok {
				return daysArg
			}
			t := timeArg.GetTime()
			return vm.NewTime(t.AddDate(int(yearsArg.GetNumber()), int(monthsArg.GetNumber()), int(daysArg.GetNumber())))
		}),

		"after": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			time1Arg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return time1Arg
			}
			time2Arg, ok := args.Get(1, vm.ValueTypeTime)
			if !ok {
				return time2Arg
			}
			t1 := time1Arg.GetTime()
			t2 := time2Arg.GetTime()
			return vm.NewBool(t1.After(t2))
		}),

		"before": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			time1Arg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return time1Arg
			}
			time2Arg, ok := args.Get(1, vm.ValueTypeTime)
			if !ok {
				return time2Arg
			}
			t1 := time1Arg.GetTime()
			t2 := time2Arg.GetTime()
			return vm.NewBool(t1.Before(t2))
		}),

		"equal": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			time1Arg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return time1Arg
			}
			time2Arg, ok := args.Get(1, vm.ValueTypeTime)
			if !ok {
				return time2Arg
			}
			t1 := time1Arg.GetTime()
			t2 := time2Arg.GetTime()
			return vm.NewBool(t1.Equal(t2))
		}),

		"format": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			layoutArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return layoutArg
			}
			t := timeArg.GetTime()
			return vm.NewString(t.Format(layoutArg.GetString()))
		}),

		"isZero": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewBool(timeArg.GetTime().IsZero())
		}),

		"date": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			t := timeArg.GetTime()
			year, month, day := t.Date()
			return vm.NewArray([]vm.Value{
				vm.NewNumber(float64(year)),
				vm.NewNumber(float64(month)),
				vm.NewNumber(float64(day)),
			})
		}),

		"getYear": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Year()))
		}),

		"getMonth": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Month()))
		}),

		"getDay": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Day()))
		}),

		"weekday": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Weekday()))
		}),

		"clock": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			t := timeArg.GetTime()
			hour, min, sec := t.Clock()

			return vm.NewObject(map[string]vm.Value{
				"hour":   vm.NewNumber(float64(hour)),
				"minute": vm.NewNumber(float64(min)),
				"second": vm.NewNumber(float64(sec)),
			})
		}),

		"getHour": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Hour()))
		}),

		"getMinute": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Minute()))
		}),

		"getSecond": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Second()))
		}),

		"getNanosecond": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Nanosecond()))
		}),

		"getUnixTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().Unix()))
		}),

		"getUnixMilliTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().UnixMilli()))
		}),

		"getUnixMicroTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().UnixMicro()))
		}),

		"getUnixNanoTime": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewNumber(float64(timeArg.GetTime().UnixNano()))
		}),

		"utc": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewTime(timeArg.GetTime().UTC())
		}),

		"local": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			return vm.NewTime(timeArg.GetTime().Local())
		}),

		// --- Duration Methods (as functions taking duration number as first arg) ---
		"getHours": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(d.Hours())
		}),

		"getMinutes": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(d.Minutes())
		}),

		"getSeconds": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(d.Seconds())
		}),

		"getMilliseconds": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(float64(d.Milliseconds()))
		}),

		"getMicroseconds": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(float64(d.Microseconds()))
		}),

		"getNanoseconds": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewNumber(float64(d.Nanoseconds()))
		}),

		"getDurationString": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			durationArg, ok := args.Get(0, vm.ValueTypeNumber)
			if !ok {
				return durationArg
			}
			d := time.Duration(durationArg.GetNumber())
			return vm.NewString(d.String())
		}),

		// --- Timezone Functions ---
		"parseInLocation": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			layoutArg, ok := args.Get(0, vm.ValueTypeString)
			if !ok {
				return layoutArg
			}
			valueArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return valueArg
			}
			locNameArg, ok := args.Get(2, vm.ValueTypeString)
			if !ok {
				return locNameArg
			}

			loc, err := time.LoadLocation(locNameArg.GetString())
			if err != nil {
				return vm.NewErrFromErr(err)
			}

			t, err := time.ParseInLocation(layoutArg.GetString(), valueArg.GetString(), loc)
			if err != nil {
				return vm.NewErrFromErr(err)
			}
			return vm.NewTime(t)
		}),

		"inLocation": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			locNameArg, ok := args.Get(1, vm.ValueTypeString)
			if !ok {
				return locNameArg
			}

			loc, err := time.LoadLocation(locNameArg.GetString())
			if err != nil {
				return vm.NewErrFromErr(err)
			}

			t := timeArg.GetTime()
			return vm.NewTime(t.In(loc))
		}),

		"getZone": vm.NewNativeFunction(func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
			timeArg, ok := args.Get(0, vm.ValueTypeTime)
			if !ok {
				return timeArg
			}
			t := timeArg.GetTime()
			name, offset := t.Zone()
			return vm.NewObject(map[string]vm.Value{
				"name":   vm.NewString(name),
				"offset": vm.NewNumber(float64(offset)), // Offset in seconds
			})
		}),
	}

	builder.RegisterModule("time", vm.NewObject(m))
}
