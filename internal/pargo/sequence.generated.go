package pargo

func Sequence1[T0, O any](psT0 Parser[T0], mapper func(T0) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0)

		return res, newState, nil
	}
}

func Sequence2[T0, T1, O any](psT0 Parser[T0], psT1 Parser[T1], mapper func(T0, T1) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1)

		return res, newState, nil
	}
}

func Sequence3[T0, T1, T2, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], mapper func(T0, T1, T2) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2)

		return res, newState, nil
	}
}

func Sequence4[T0, T1, T2, T3, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], mapper func(T0, T1, T2, T3) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3)

		return res, newState, nil
	}
}

func Sequence5[T0, T1, T2, T3, T4, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], mapper func(T0, T1, T2, T3, T4) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4)

		return res, newState, nil
	}
}

func Sequence6[T0, T1, T2, T3, T4, T5, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], mapper func(T0, T1, T2, T3, T4, T5) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5)

		return res, newState, nil
	}
}

func Sequence7[T0, T1, T2, T3, T4, T5, T6, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], mapper func(T0, T1, T2, T3, T4, T5, T6) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6)

		return res, newState, nil
	}
}

func Sequence8[T0, T1, T2, T3, T4, T5, T6, T7, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], mapper func(T0, T1, T2, T3, T4, T5, T6, T7) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7)

		return res, newState, nil
	}
}

func Sequence9[T0, T1, T2, T3, T4, T5, T6, T7, T8, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8)

		return res, newState, nil
	}
}

func Sequence10[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9)

		return res, newState, nil
	}
}

func Sequence11[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10)

		return res, newState, nil
	}
}

func Sequence12[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], psT11 Parser[T11], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT11, newState, err := psT11(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10, resT11)

		return res, newState, nil
	}
}

func Sequence13[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], psT11 Parser[T11], psT12 Parser[T12], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT11, newState, err := psT11(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT12, newState, err := psT12(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10, resT11, resT12)

		return res, newState, nil
	}
}

func Sequence14[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], psT11 Parser[T11], psT12 Parser[T12], psT13 Parser[T13], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT11, newState, err := psT11(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT12, newState, err := psT12(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT13, newState, err := psT13(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10, resT11, resT12, resT13)

		return res, newState, nil
	}
}

func Sequence15[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], psT11 Parser[T11], psT12 Parser[T12], psT13 Parser[T13], psT14 Parser[T14], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT11, newState, err := psT11(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT12, newState, err := psT12(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT13, newState, err := psT13(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT14, newState, err := psT14(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10, resT11, resT12, resT13, resT14)

		return res, newState, nil
	}
}

func Sequence16[T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, O any](psT0 Parser[T0], psT1 Parser[T1], psT2 Parser[T2], psT3 Parser[T3], psT4 Parser[T4], psT5 Parser[T5], psT6 Parser[T6], psT7 Parser[T7], psT8 Parser[T8], psT9 Parser[T9], psT10 Parser[T10], psT11 Parser[T11], psT12 Parser[T12], psT13 Parser[T13], psT14 Parser[T14], psT15 Parser[T15], mapper func(T0, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) O) Parser[O] {
	return func(state State) (O, State, error) {
		newState := state

		resT0, newState, err := psT0(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT1, newState, err := psT1(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT2, newState, err := psT2(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT3, newState, err := psT3(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT4, newState, err := psT4(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT5, newState, err := psT5(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT6, newState, err := psT6(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT7, newState, err := psT7(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT8, newState, err := psT8(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT9, newState, err := psT9(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT10, newState, err := psT10(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT11, newState, err := psT11(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT12, newState, err := psT12(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT13, newState, err := psT13(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT14, newState, err := psT14(newState)
		if err != nil {
			return zero[O](), state, err
		}

		resT15, newState, err := psT15(newState)
		if err != nil {
			return zero[O](), state, err
		}

		res := mapper(resT0, resT1, resT2, resT3, resT4, resT5, resT6, resT7, resT8, resT9, resT10, resT11, resT12, resT13, resT14, resT15)

		return res, newState, nil
	}
}
