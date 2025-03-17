#import "@preview/zebraw:0.4.3": *
#import "lib.typ": ieee


#show: ieee.with(
  title: [Weaver - Scripting language],
  // abstract: [],
  authors: (
    (
      name: "Youssef Ahmed",
      department: [Department of Computer Science],
      organization: [Cairo University],
      location: [Egypt, Cairo],
      email: "joetifa2003@gmail.com"
    ),
  ),
  figure-supplement: [Fig.],
)

#show: zebraw
#show raw: set text(
  font: "0xProto Nerd Font Mono", 
  spacing: 100% + 0pt, 
  tracking: 0pt,
  ligatures: true,
)
#set raw(syntaxes: ("./weaver.syntax.yml"))
#show heading: set text(10pt, weight: 600)
#set heading(numbering: "1.")

= Introduction

Weaver is a scripting language that allows users to write and execute code in a simple and intuitive way. It is designed to be easy to learn and use, with a focus on simplicity and readability.

It's designed to fill in the gap in the current dynamic scripting languages, like Python, Javascript, Bash etc..

= Basics

Weaver is a dynamically typed language, which means that variables can hold values of different types. The type of a value is determined at runtime, based on the value itself.

== Types Of Expressions

Weaver has the typical types of expressions as shown in @fig:expressions, What stands out more is the fact that functions are first class citizens in Weaver, and can be passed around as values anywhere.

#figure(
```weaver
"Hello World!"          // string
123                     // int
1.23                    // float
true|false              // bool
[1, "2", [3]]           // array
{"a": 1, b: 2}          // dict
|a, b| a + b            // lambda
|a, b| { return a + b } // function
nil                     // nil 
```,
  caption: [Weaver expressions.],
) <fig:expressions>

== Operators

Generally speaking, Weaver has many of the same operators as other C-like languages.

#figure(
table(
  columns: (auto, auto, auto, auto),
  table.header(
    "operand", "operator", "operand", "result"
  ),
  "number",    `+`,  "number",        "number",
  "number",    `-`,  "number",        "number",
  "number",    `*`,  "number",        "number",
  "number",    `/`,  "number",        "number",
  "string",    `+`,  "string",        "string",
  "any A",     `|>`, "any B()",       "B(A)",
),
  caption: [
  Binary operators in Weaver.
  ],
) <fig:binary-operators-table>

For any other combination of types not in @fig:binary-operators-table is illegal, and Weaver will throw an error.

For equality operator (`==`|`!=`), two types are considered equal if they are the same type AND have the same value, as shown in @fig:equality-operators.

#figure(
table(
  columns: (1fr, 1fr, 1fr, auto),
  "A", `==`, "B", $"true" "iff" "type"(A) = "type"(B) and A = B$,
  "A", `!=`, "B", $"true" "iff" "type"(A) != "type"(B) or A != B$,
),
  caption: [Equality operators in Weaver.],
) <fig:equality-operators>

== Type Coercion

In Weaver, there is no type coercion, I.E. values are not automatically converted to a different type based on the operators, famous examples of this in javascript shown in @fig:js-type-coercion.

#figure(
```javascript
true + false    == 1
12 / "6"        == 2
"foo" + 15 + 3  == "foo153"
{} + []         == 0
[] == ![]       == true
0 == "0"        == true
0 == []         == true
"0" == []       == false // !!!
```,
  caption: [Javascript type coercion.],
) <fig:js-type-coercion>

Instead, Weaver has a stricter type system, which means there are no implicit type conversions, instead you should convert types manually, as shown in @fig:weaver-type-conversion.

#figure(
```weaver
number(true) + number(false)  == 1 
12 / number("6")           == 12
"foo" + string(15 + 3)  == "foo18" 
```,
  caption: [Weaver type conversion.],
) <fig:weaver-type-conversion>

This makes reading the code much easier, Since what you see is what you get.

If operators are used with incorrect types, Weaver will error at runtime, and if equality operators (==|!=) are used with diffenrent types, it always returns false, as shown in @fig:weaver-equality-different-types.

#figure(
```weaver
// illegal operands dict + array
{} + []

// illegal operands string + int
"foo" + 8 

// array == bool
   []    == ![]   == false 

// int == string
   0   == "0"     == false 

// int == array
   0   == []      == false 

// string == array
   "0"    == []   == false
```,
  caption: [Weaver operators with incorrect/different type operands.],
) <fig:weaver-equality-different-types>

== Control Flow

If statements are used to execute a block of code if a certain condition expression is true, as shown in @fig:if.

#figure(
```weaver
if 5 > 3 {
  echo("5 is greater than 3")
}
```,
  caption: [If statement in Weaver.],
) <fig:if>


If statements can also have an else block to execute if the condition is false, as shown in @fig:if-else.

#figure(
```weaver
if 5 > 3 {
  "5 is greater than 3" 
    |> echo()
} else {
  "5 is not greater than 3" 
    |> echo()
}
```,
  caption: [If statement in Weaver.],
) <fig:if-else>


While loops are used to repeatedly execute a block of code as long as a certain condition expression is true, as shown in @fig:while.

#figure(
```weaver
i := 0
while i < 5 {
  echo(i)
  i += 1
}
```,
  caption: [While loop in Weaver, printing numbers from 0 to 4.],
) <fig:while>

For loops are used as an alternative to while loops, typically used when iterating over a collection of items like arrays, @fig:for is another way for writing @fig:while.

#figure(
```weaver
  for i := 0; i < 5; i += 1 {
    echo(i)
  }
```,
  caption: [For loop in Weaver, printing numbers from 0 to 4.],
) <fig:for>

Iterating over an array is a common use case for for loops, as shown in @fig:for-array.

#figure(
```weaver
  arr := [1, 2, 3, 4]
  for i := 0; i < len(arr); i++ {
    echo(arr[i])
  }
```,
  caption: [For loop over an array, printing numbers from 1 to 4.],
) <fig:for-array>

#set page(columns: 1)

== Pattern Matching <sec:pattern-matching>

Weaver provides powerful pattern matching capabilities that allow developers to write expressive and concise code. Pattern matching can be used to match values against specific patterns and execute corresponding code blocks. The syntax follows the form:

```weaver
match expression {
  pattern1 => {
    // code to execute if pattern1 matches
  },
  pattern2 => {
    // code to execute if pattern2 matches
  },
  _ => {
    // default case if no patterns match
  }
}
```

=== Basic Patterns

Weaver supports matching against literal values including integers, floats, strings, and booleans:

```weaver
x := 5
match x {
  0 => echo("zero"),
  5 => echo("five"),
  _ => echo("other")
}
```

=== Type Matching

Patterns can match based on value types using type predicates:

```weaver
match value {
  string(s) => echo("got string: " + s),
  number(n) => echo("got number: " + string(n)),
  _ => echo("other type")
}
```

=== Destructuring

Weaver supports destructuring of arrays and objects in patterns:

```weaver
// Array destructuring
match [1, 2, 3] {
  [a, b, c] => echo(a + b + c),
  _ => echo("no match")
}

// Object destructuring  
match {name: "Alice", age: 30} {
  {name: n, age: a} => echo(n + " is " + string(a)),
  _ => echo("no match")
}
```

=== Nested Patterns

Patterns can be nested to match complex data structures:

```weaver
students := [
  {name: "Alice", grades: [90, 85]},
  {name: "Bob", grades: [80, 75]}
]

match students[0] {
  {name: "Alice", grades: [math, _]} => echo("Alice's math grade: " + string(math)),
  _ => echo("no match")
}
```

=== Range Matching

Weaver supports range patterns for numeric values:

```weaver
match age {
  0..17 => echo("child"),
  18..64 => echo("adult"), 
  65.. => echo("senior"),
  else => echo("invalid age")
}
```

=== Guards

Additional conditions can be added to patterns using guards:

```weaver
match student {
  {name: n, age: a} if a >= 18 => echo(n + " is an adult"),
  {name: n, age: a} => echo(n + " is a minor")
}
```

Pattern matching in Weaver provides a powerful way to handle complex conditional logic in a readable and maintainable way. The combination of literal matching, type matching, destructuring, and guards makes it suitable for a wide range of use cases.


== Error Handling

Error handling is a critical aspect of any programming language, and Weaver provides a robust system for creating, propagating, and handling errors.

=== Creating Errors

In Weaver, errors are first-class values created using the `error` function, which takes an error message and optional data.

#figure(
```weaver
// Create a simple error
e := error("File not found")

// Create an error with additional data
e := error("User not authorized", {userId: 123, role: "guest"})
```,
  caption: [Creating errors in Weaver.],
) <fig:creating-errors>

=== Error Propagation

Weaver uses a unique approach to error propagation. By default, when a function call returns an error, it is automatically propagated up the call stack. This behavior can be overridden by using the bang operator (`!`) after a function call.

#figure(
```weaver
// Automatic error propagation
result := riskyFunction()  // If riskyFunction returns an error, it propagates

// Explicit error handling with bang operator
result := riskyFunction()! // Error won't propagate, must be handled manually
```,
  caption: [Error propagation in Weaver.],
) <fig:error-propagation>

=== Pattern Matching on Errors

Errors can be matched and destructured using pattern matching.

#figure(
```weaver
e := riskyFunction()! // notice the bang operator, error won't auto propagate

match e {
  error(msg, {name: n}) => {
    echo("Error message: " + msg + ", name: " + n)
  },
  else => {
    echo("Not an error")
  }
}
```,
  caption: [Pattern matching on errors in Weaver.],
) <fig:error-pattern-matching>

The `error(msg, data)` pattern is used to match on errors message and optional data, and can be used with other patterns shown in @sec:pattern-matching.

```weaver
r := http:get("https://foo.com") 

match r {
  error(_, {status: 400..499}) => {
    echo("Error with status in range 400-499")
  },
  else => {
    echo("Not an error")
  }
}
```

And it can be also used without any extra matching

#zebraw(
  highlight-lines: range(10, 13),
```weaver
r := http:get("https://foo.com")!

match r {
  error(_, {status: 400..499}) => {
    echo("Error with status in range 400-499")
  },
  error(msg, {status: 500..599}) => {
    echo("Error with status in range 500-599 with message " + msg)
  },
  error() => { // here we are just matching the type 
    echo("Some other error")
  },
  _ => {
    echo("Not an error")
  }
}
```
)

=== Comparison with Other Languages

Weaver's error handling approach differs from many mainstream languages, offering a balance between explicit and implicit error handling.

==== Go

Go uses explicit error handling with return values:

#figure(
```go
// Go error handling
result, err := riskyFunction()
if err != nil {
    // Handle error
    return err
}
// Use result
```,
  caption: [Error handling in Go.],
) <fig:go-error-handling>

Go's approach is more explicit than Weaver's default behavior, requiring errors to be checked at each call site. Weaver's automatic propagation reduces boilerplate while still allowing explicit handling when needed.

==== Rust

Rust uses the `Result` type for error handling:

#figure(
```rust
// Rust error handling
match risky_function() {
    Ok(result) => {
        // Use result
    },
    Err(e) => {
        // Handle error
    }
}

// Or using the ? operator
let result = risky_function()?;
```,
  caption: [Error handling in Rust.],
) <fig:rust-error-handling>

Rust's `?` operator is similar to Weaver's automatic propagation, but Rust requires functions to explicitly declare that they can return errors through their return type.

==== JavaScript

JavaScript traditionally uses exceptions for error handling:

#figure(
```javascript
// JavaScript error handling
try {
    const result = riskyFunction();
    // Use result
} catch (e) {
    // Handle error
}
```,
  caption: [Error handling in JavaScript.],
) <fig:javascript-error-handling>

JavaScript's exception-based approach is more implicit than Weaver's, making it harder to track error flows. Weaver's approach provides better visibility into error paths while still being concise.
