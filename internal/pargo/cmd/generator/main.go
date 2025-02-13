package main

import (
	"fmt"
	"os"
	"strings"
)

// func Sequence[T any, O any](mapper func(T) O, psT Parser[T]) Parser[O] {
// return func(state State) (O, State, error) {
// 	resT, newState, err := psT(state)
// 	if err != nil {
// 		return zero[O](), state, err
// 	}
//
// 	res := mapper(resT)
//
// 	return res, newState, nil
// }
// }

func main() {
	var res strings.Builder
	res.WriteString("package main\n\n")

	for i := range 16 {
		genericParameters := ""
		for j := range i + 1 {
			genericParameters += fmt.Sprintf("T%d", j)
			if j != i {
				genericParameters += ", "
			}
		}

		parameters := ""
		for j := range i + 1 {
			parameters += fmt.Sprintf("psT%d Parser[T%d]", j, j)
			if j != i {
				parameters += ", "
			}
		}

		parsers := ""
		for j := range i + 1 {
			parsers += fmt.Sprintf(`
    resT%d, newState, err := psT%d(newState)
    if err != nil {
      return zero[O](), state, err
    }
`, j, j)
		}

		results := ""
		for j := range i + 1 {
			results += fmt.Sprintf("resT%d", j)
			if j != i {
				results += ", "
			}
		}

		res.WriteString(fmt.Sprintf(`func Sequence%d[%s, O any](%s, mapper func(%s) O) Parser[O] {
  return func(state State) (O, State, error) {
    newState := state

    %s 

    res := mapper(%s)

    return res, newState, nil
  }
}

`, i+1, genericParameters, parameters, genericParameters, parsers, results))
	}

	f, err := os.Create("sequence.generated.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(res.String())
	if err != nil {
		panic(err)
	}
}
