package vm_test

import (
	"fmt"
	"testing"

	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func TestVM(t *testing.T) {
	tests := []string{
		0: `
			x := 1
			x == 1 | assert()
		`,
		1: `
			x := 1
			x = 2
			x == 2 | assert()
		`,
		2: `
			x := 1
			cond := true

			if cond {
				x = 2
			} else {
				x = 3		
			}

			x == 2 | assert()
		`,
		3: `
			x := 1
			cond := false

			if cond {
				x = 2
			} else {
				x = 3		
			}

			x == 3 | assert()
		`,
		4: `
			even := 0
			odd := 0
			for i := 0; i < 10; i = i + 1 {
				if i % 2 == 0 {
					even = even + 1
				}	else {
					odd = odd + 1
				}
			}	

			even == 5 | assert()
			odd == 5  | assert()
		`,
		5: `
			even := 0
			odd := 0
			isEven := |x| x % 2 == 0
			for i := 0; i < 10; i = i + 1 {
				if isEven(i) {
					even = even + 1
				}	else {
					odd = odd + 1
				}
			}	

			even == 5 | assert()
			odd == 5  | assert()
		`,
		6: `
			even := 0
			odd := 0
			isEven := |x| x % 2 == 0
			
			i := 0
			while i < 10 {
				if isEven(i) {
					even = even + 1
				}	else {
					odd = odd + 1
				}

				i = i + 1
			}	

			even == 5 | assert()
			odd == 5  | assert()
		`,
		7: `
			adder := |x| |y| x + y
			addFive := adder(5)
			addFive(5) == 10 | assert()
		`,
		8: `
			x := [1, 2, 3]
			l := x 
				| map(|x| x + 1) 
				| filter(|x| x % 2 == 0) 
				| len() 
			l == 2 | assert()
		`,
		9: `
			a := 1
			b := 2

			x := [1, 2, 3]
			l := x 
				| map(|x| x + a) 
				| filter(|x| x % b == 0) 
				| len() 
			l == 2 | assert()
		`,
		10: `
			x := [1, 2, 3]
			x[0] = 2
			x[0] == 2 | assert()
		`,
		11: `
			x := [[1], [2], [3]]
			x[0][0] = 2
			x[0][0] == 2 | assert()
		`,
		12: `
			x := { a: 1 }
			x["a"] = 2
			x.a == 2 | assert()
		`,
		13: `
			x := { a: 1 }
			x.a = 2
			x.a == 2 | assert()
		`,
		14: `
			x := [{a: [9]}]
			x[0].a[0] == 9 | assert()
		`,
		15: `
			x := [{a: [9]}]
			x[0].a[0] = 41
			x[0].a[0] == 41 | assert()
		`,
		16: `
			x := [1, 2, 3]
			x | push(4)
			assert(len(x) == 4)
			x[0] == 1 | assert()
			x[1] == 2 | assert()
			x[2] == 3 | assert()
			x[3] == 4 | assert()
		`,
		17: `
		x := 0
		y := 5
		match x {
			0 => {
				match y {
					0 => {
						false | assert()
					},
					1 => {
						false | assert()
					},
					2 => {
						false | assert()
					},
					4 => {
						true | assert()
					}
				}
			},
			1 => {
				false | assert()
			},
			2 => {
				false | assert()
			},
			3 => {
				false | assert()
			}
		}
		`,
	}

	for i, tc := range tests {
		for _, opt := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d opt=%t", i, opt), func(t *testing.T) {
				p, err := parser.Parse(tc)
				if err != nil {
					t.Fatal(fmt.Errorf("failed to parse: %w", err))
				}

				irc := ir.NewCompiler()

				ircr, err := irc.Compile(p)
				if err != nil {
					t.Fatal(fmt.Errorf("failed to compile ir: %w", err))
				}

				c := compiler.New(compiler.WithOptimization(opt))
				frame, constants, err := c.Compile(ircr)
				if err != nil {
					t.Fatal(fmt.Errorf("failed to compile: %w", err))
				}
				vm := vm.New(constants, frame.Instructions, len(frame.Vars))
				vm.Run()
			})
		}
	}
}
