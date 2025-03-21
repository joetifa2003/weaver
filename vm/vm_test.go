package vm_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joetifa2003/weaver/builtin"
	"github.com/joetifa2003/weaver/compiler"
	"github.com/joetifa2003/weaver/ir"
	"github.com/joetifa2003/weaver/parser"
	"github.com/joetifa2003/weaver/vm"
)

func TestVM(t *testing.T) {
	t.Parallel()

	tests := []string{
		0: `
			x := 1;
			x == 1 |> assert();
		`,
		1: `
			x := 1;
			x = 2;
			x == 2 |> assert();
		`,
		2: `
			x := 1;
			cond := true;

			if cond {
				x = 2;
			} else {
				x = 3;
			}

			x == 2 |> assert();
		`,
		3: `
			x := 1;
			cond := false;

			if cond {
				x = 2;
			} else {
				x = 3;
			}

			x == 3 |> assert();
		`,
		4: `
			even := 0;
			odd := 0;
			for i := 0; i < 10; i = i + 1 {
				if i % 2 == 0 {
					even = even + 1;
				}	else {
					odd = odd + 1;
				}
			}	

			even == 5 |> assert();
			odd == 5  |> assert();
		`,
		5: `
			even := 0;
			odd := 0;
			isEven := |x| x % 2 == 0;

			for i := 0; i < 10; i = i + 1 {
				if isEven(i) {
					even = even + 1;
				}	else {
					odd = odd + 1;
				}
			}	

			even == 5 |> assert();
			odd == 5  |> assert();
		`,
		6: `
			even := 0;
			odd := 0;
			isEven := |x| x % 2 == 0;
			
			i := 0;
			while i < 10 {
				if isEven(i) {
					even = even + 1;
				}	else {
					odd = odd + 1;
				}

				i = i + 1;
			}	

			even == 5 |> assert();
			odd == 5  |> assert();
		`,
		7: `
			adder := |x| |y| x + y;
			addFive := adder(5);
			addFive(5) == 10 |> assert();
		`,
		8: `
			x := [1, 2, 3];
			l := x 
				|> map(|x| x + 1) 
				|> filter(|x| x % 2 == 0) 
				|> len() ;
			l == 2 |> assert();
		`,
		9: `
			a := 1;
			b := 2;

			x := [1, 2, 3];
			l := x 
				|> map(|x| x + a) 
				|> filter(|x| x % b == 0) 
				|> len();

			l == 2 |> assert();
		`,
		10: `
			x := [1, 2, 3];
			x[0] = 2;
			x[0] == 2 |> assert();
		`,
		11: `
			x := [[1], [2], [3]];
			x[0][0] = 2;
			x[0][0] == 2 |> assert();
		`,
		12: `
			x := { a: 1 };
			x["a"] = 2;
			x.a == 2 |> assert();
		`,
		13: `
			x := { a: 1 };
			x.a = 2;
			x.a == 2 |> assert();
		`,
		14: `
			x := [{a: [9]}];
			x[0].a[0] == 9 |> assert();
		`,
		15: `
			x := [{a: [9]}];
			x[0].a[0] = 41;
			x[0].a[0] == 41 |> assert();
		`,
		16: `
			x := [1, 2, 3];
			x |> push(4);
			assert(len(x) == 4);
			x[0] == 1 |> assert();
			x[1] == 2 |> assert();
			x[2] == 3 |> assert();
			x[3] == 4 |> assert();
		`,
		17: `
		x := 0;
		y := 5;
		match x {
			0 => {
				match y {
					0 => {
						false |> assert();
					},
					1 => {
						false |> assert();
					},
					2 => {
						false |> assert();
					},
					4 => {
						false |> assert();
					},
					5 => {
						true |> assert();
					},
					else => {
						false |> assert();
					}
				}
			},
			1 => {
				false |> assert();
			},
			2 => {
				false |> assert();
			},
			3 => {
				false |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		18: `
		x := 0.5;
		y := 0.7;
		match x {
			0.5 => {
				match y {
					0 => {
						false |> assert();
					},
					1 => {
						false |> assert();
					},
					2 => {
						false |> assert();
					},
					4 => {
						false |> assert();
					},
					0.7 => {
						true |> assert();
					},
					else => {
						false |> assert();
					}
				}
			},
			1 => {
				false |> assert();
			},
			2 => {
				false |> assert();
			},
			3 => {
				false |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		19: `
		match "foo" {
			"bar" => {
				false |> assert();
			},
			"baz" => {
				false |> assert();
			},
			"foo" => {
				true |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		20: `
		match [[0], 1, 2] {
			[[0], 1, 2] => {
				true |> assert();
			},
			[2, 3, 4] => {
				false |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		21: `
			a := || false |> assert();
			true && true && false && a();
		`,
		22: `
			a := || false |> assert();
			false || true || a();
		`,
		23: `
			aCalled := {value: false};
			a := || {
				aCalled.value = true;
			};
			true && true && true && a();
			aCalled.value |> assert();
		`,
		24: `
			aCalled := {value: false};
			a := || {
				aCalled.value = true;
			};
			false || false || false || a();
			aCalled.value |> assert();
		`,
		25: `
		c := [0, {name: "hello"}];
		match c {
			[0, {name: "hello"}] => {
				true |> assert();
			},
			[2, 3, 4] => {
				false |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		26: `
		students := [
			{name: "joe", age: 30},
			{name: "foo", age: 20},
			{name: "bar", age: 10},
		];

		res := "";
		res2 := "";

		for i := 0; i < len(students); i = i + 1 {
			match students[i] {
				{name: n, age: a} if a >= 10 && a <= 20 => {
					res = res + n;
				},
				else => {
					res2 = res2 + students[i].name;
				}
			}
		}

		res == "foobar" |> assert();
		res2 == "joe" |> assert();
		`,
		27: `
			even := 0;
			odd := 0;
			for i := 0; i < 10; i++ {
				if i % 2 == 0 {
					even = even + 1;
				}	else {
					odd = odd + 1;
				}
			}	

			even == 5 |> assert();
			odd == 5  |> assert();
		`,
		28: `
			x := 1;
			x--;
			x == 0 |> assert();
		`,
		29: `
			even := 0;
			odd := 0;
			for i := 10; i > 0; i-- {
				if i % 2 == 0 {
					even++;
				} else {
					odd++;
				}
			}

			even == 5 |> assert();
			odd == 5  |> assert();
		`,
		30: `
			fib := |n| {
				if n <= 1 {
					return n;
				}

				return fib(n - 1) + fib(n - 2);
			};

			fib(10) == 55 |> assert();
		`,
		31: `
		true  ? assert(true)  | assert(false);
		false ? assert(false) | assert(true);
		`,
		32: `
		e := error("test error", {name: "test"});
		e.data.name = "hi";

		match e {
			error(_, {name: n}) => {
				n == "hi" |> assert();
			},
			else => {
				false |> assert();
			}
		}
		`,
		33: `
		counter := || {
			x := 0;

			return {
				increment: || x = x + 1,
				decrement: || x = x - 1,
				value: || x,
			};
		};

		c := counter();

		c.increment();
		c.value() == 1 |> assert();
		c.increment();
		c.value() == 2 |> assert();
		c.decrement();
		c.value() == 1 |> assert();
		`,
		34: `
		counter := || {
			x := 0;

			return {
				increment: || x++,
				decrement: || x--,
				value: || x,
			};
		};

		c := counter();

		c.increment();
		c.value() == 1 |> assert();
		c.increment();
		c.value() == 2 |> assert();
		c.decrement();
		c.value() == 1 |> assert();
		`,
		35: `
		x := "foo";
		match x {
			string(s) => {
				s == "foo" |> assert();
				return;
			},
			number() => {
				false |> assert();
			},
			else => {
				false |> assert();
			}
		}

		false |> assert();
		`,
		36: `
		x := 15;
		match x {
			string() => {
				false |> assert();
			},
			number(n) => {
				n == 15|> assert();
				return;
			},
			else => {
				false |> assert();
			}
		}

		false |> assert();
		`,
		37: `
			tempD := tempDir();
			tempD |> io:isDir() |> assert();
			
			path := io:join(tempD, "test.txt");
			io:dirname(path) == tempD |> assert();
			io:basename(path) == "test.txt" |> assert();
			io:extname(path) == ".txt" |> assert();
			
			io:writeFile(path, "hello world");
			io:exists(path) |> assert();
			io:readFile(path) == "hello world" |> assert();
			
			subDir := io:join(tempD, "subdir");
			io:mkdir(subDir);
			io:isDir(subDir) |> assert();
			
			io:size(path) == 11 |> assert();
			!io:isDir(path) |> assert();
			
			newPath := io:join(tempD, "renamed.txt");
			io:rename(path, newPath) ;
			io:exists(newPath) |> assert();
			!io:exists(path) |> assert();
			
			io:remove(newPath);
			!io:exists(newPath) |> assert();
			io:remove(subDir);
			!io:exists(subDir) |> assert();
		`,
		38: `
		match 19 {
			0..19 => {
				return;
			}
		}
		return error("19 should match 0..19");
		`,
		39: `
		match len([0, 1, 2]) {
			2..3 => {
				return;
			}
		}
		return error("3 should match 2..3");
		`,
		40: `
		match {a: {b: [1, 2, "3", 4.5]}} |> json:stringify() |> json:parse() {
			{a: {b: [1, 2, "3", 4.5]}} => {
				return;
			}
		}

		false |> assert();
		`,
		41: `
			match { status: 200 } {
				{ status } => {
					status == 200 |> assert();
					return;
				}
			}
			false |> assert();
		`,
		42: `
			i := 0;

			f := || i;

			i = i + 1;

			f() == 1 |> assert();
		`,
	}

	for i, tc := range tests {
		for _, opt := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d opt=%t", i, opt), func(t *testing.T) {
				t.Parallel()
				assert := require.New(t)

				p, err := parser.Parse(tc)
				assert.NoError(err)

				irc := ir.NewCompiler()
				ircr, err := irc.Compile(p)
				assert.NoError(err)

				reg := builtin.NewRegBuilderFrom(builtin.StdReg).
					RegisterFunc("tempDir", func(v *vm.VM, args vm.NativeFunctionArgs) vm.Value {
						return vm.NewString(t.TempDir())
					}).
					Build()
				c := compiler.New(reg, compiler.WithOptimization(opt))
				instructions, vars, constants, err := c.Compile(ircr)
				assert.NoError(err)

				executor := vm.NewExecutor(constants)
				val := executor.Run(
					&vm.Frame{
						Instructions: instructions,
						NumVars:      vars,
						HaltAfter:    true,
					},
					0,
				)
				if val.IsError() {
					t.Error(val.GetError())
				}
			})
		}
	}
}
