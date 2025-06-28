#import "@preview/zebraw:0.4.3": *

#let logo_width = 3cm

// Header with logos
#align(center)[
  #grid(
    columns: (1fr, 1fr),
    column-gutter: 2cm,
    align: center,
    [
      #image("logo1.jpg", width: logo_width)
    ],
    [
      #image("logo2.jpg", width: logo_width)
    ]
  )
]

#v(2cm)

// Project Title and Description
#align(center)[
  #text(size: 24pt, weight: "bold")[
    Weaver
  ]
  
  
  #text(size: 16pt)[
    A Scripting Language
  ]
]

#v(3cm)

// Team Members Section
#align(center)[
  #text(size: 18pt, weight: "bold")[
    Submitted By
  ]

  #grid(
    columns: 1,
    row-gutter: 0.8cm,
    [#text(size: 14pt)[Youssef Ahmed Nour El-Dien - 2127137]]
  )
]

#v(1fr)

// Professor Supervisor at bottom
#align(center)[
  #line(length: 60%, stroke: 0.5pt)
  #v(0.5cm)
  #text(size: 16pt, weight: "bold")[
    SUPERVISED BY
  ]
  #v(0.3cm)
  #text(size: 14pt)[
    Dr. Hossam Hassan
  ]
  #v(0.2cm)
  #text(size: 12pt, weight: "bold")[
    Cairo University
  ]
  #v(0.2cm)
  #text(size: 12pt, weight: "bold")[
    Department of Computer Science
  ]
]

#v(1cm)

// Date at the very bottom
#align(center)[
  #text(size: 12pt)[
    #datetime.today().display("[month repr:long] [year]")
  ]
]

#pagebreak()

= Abstract

The proliferation of scripting languages like Python and JavaScript has simplified many programming tasks. However, they often introduce complexities such as implicit type coercion and verbose asynchronous code, which can lead to subtle bugs and reduced readability. This paper introduces Weaver, a novel scripting language designed to address these shortcomings. Weaver's core philosophy is to provide a simple, expressive, and predictable programming experience, combining the convenience of dynamic typing with the safety of a strong type system.

Key features of Weaver include the absence of type coercion, first-class functions, a powerful pipe operator for elegant data flow, and comprehensive pattern matching capabilities. These features encourage a functional programming style that is both concise and easy to reason about.

Furthermore, Weaver incorporates a modern concurrency model based on lightweight green threads, called Fibers. This allows for writing non-blocking, concurrent code in a simple, sequential style, eliminating the need for `async/await` syntax. Benchmarks demonstrate that this approach not only simplifies code but also yields significant performance and memory efficiency improvements compared to established runtimes like Node.js.

By integrating these features, Weaver aims to fill a gap in the current landscape of scripting languages, offering a tool that is both a joy to use and powerful enough for demanding applications.

#pagebreak()

= Project Goals

The development of Weaver is guided by a set of core principles aimed at addressing common pain points in modern scripting languages. The primary goals for the project are as follows:

- *Simplicity and Readability:*
  To create a language with a clean, intuitive syntax that is easy to learn and read. Weaver is designed to be expressive without being cryptic, enabling developers to write maintainable code with less cognitive overhead. Features like the pipe operator (`|>`), pattern matching, and first-class functions are central to this goal.

- *Predictability and Safety:*
  A core objective is to eliminate the unpredictable behavior caused by implicit type coercion, a common source of bugs in languages like JavaScript. By enforcing a strong type system without coercion, Weaver ensures that operations are explicit and that what you see is what you get.

- *Improved Error Handling:*
  Weaver aims to provide a robust and ergonomic error handling mechanism. The goal is to offer a system that is less verbose than Go's explicit error returns but more predictable than traditional exceptions. By treating errors as values and integrating them with pattern matching, Weaver makes error handling a natural part of the development workflow.

- *Modern Concurrency with Fibers:*
  To simplify concurrent programming, Weaver introduces Fibers, a model based on green threads. This approach abstracts away the complexity of asynchronous code, eliminating the need for `async/await` syntax. The goal is to allow developers to write non-blocking, concurrent applications in a straightforward, sequential style, while achieving high performance and efficiency.

- *Expressive Power through Pattern Matching:*
  To provide developers with powerful tools for handling complex data structures and control flow, Weaver includes a comprehensive pattern matching system. This feature allows for elegant destructuring, type matching, and conditional logic in a single, unified construct.

- *High Performance:*
  A key goal is to deliver a high-performance runtime that can compete with and, in some cases, surpass established platforms like Node.js. By building a multi-threaded runtime and using an efficient memory model, Weaver is designed to handle high-concurrency workloads with lower resource consumption.

- *Comprehensive Standard Library:*
  To make Weaver a practical tool for real-world applications, the project aims to provide a rich and useful standard library. This includes built-in modules for common tasks such as handling HTTP requests, file I/O, and JSON parsing, enabling developers to be productive out of the box.

#pagebreak()

= Table of Contents
#outline()

#pagebreak()

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

Everything is plain objects and functions, most of the code is just concatenating functions together to solve problems.

#figure(
```weaver
// before
arr := [1, 2, 3, 4]
echo(len(filter(arr, |n| n % 2 == 0))) // [2, 4]

// after
arr := [1, 2, 3, 4]
arr
    |> filter(|n| n % 2 == 0)
    |> len()
    |> echo() // [2, 4]
```,
  caption: [Piping functions in Weaver. A core concept for writing readable code.],
) <fig:piping-functions>


= Basics

Weaver is a dynamically typed language, which means that variables can hold values of different types. The type of a value is determined at runtime, based on the value itself.

== Values

Weaver has the typical types of expressions as shown in @fig:expressions, What stands out more is the fact that functions are first class citizens in Weaver, and can be passed around as values anywhere.

#figure(
```weaver
"Hello World!"          // string
123                     // number
1.23                    // number
true|false              // bool
[1, "2", [3]]           // array
{"a": 1, b: 2}          // dict (b is shorthand for "b")
|a, b| a + b            // lambda
|a, b| { return a + b } // function
nil                     // nil
```,
  caption: [Weaver expressions.],
) <fig:expressions>

There is only one type in Weaver that indicates the absence of a value, it is `nil`.

== Binary Operations

Generally speaking, Weaver has many of the same operators as other C-like languages. The highlight being the pipe operator (`|>`), and lazy evaluation of binary `&&` and `||` operators that works for booleans and other values also based on if they are "truthy" or not.

#figure(
```weaver
1 + 2             // 3     (number)
1.0 + 2           // 3.0   (number)
2.3 + 3.4         // 5.7   (number)
1 - 2             // -1    (number)
1 * 2             // 2     (number)
1 / 2             // 0.5   (number)
8 % 2             // 0     (number)

true && false     // false (boolean)
true || false     // true  (boolean)

nil || 1          // 1     (number)
error() || "foo"  // "foo" (string)

"hello " + "world" // "hello world" (string)

add := |a, b| a + b
add(1, 2)         // 3     (number)
1 |> add(2)       // 3     (number)
```,
  caption: [Binary operations in Weaver.],
) <fig:binary-ops>

== Truthy and Falsy Values

Weaver boolean operators work with boolean expressions `true` and `false` as well as any other value in the language. Values that are considered "falsy" are `nil`, `error`, and `false`. Everything else is "truthy".

The `||` operator returns the first "truthy" value, and `&&` returns the first "falsy" value (or the last value if all are truthy).

#figure(
```weaver
// Examples of truthy and falsy values
if true  { echo("true is truthy!") }  // Prints
if 0     { echo("0 is truthy!") }     // Prints
if ""    { echo("'' is truthy!") }    // Prints
if nil   { echo("nil is truthy!") }   // Does not print
if false { echo("false is truthy!") } // Does not print

// Or operator (||) returns the first truthy value
nil || "foo"  // "foo" (string)
true || "foo" // true (boolean)

// And operator (&&) returns the first falsy value
false && "foo" // false (boolean)
nil && true    // nil
true && nil    // nil
true && "foo"  // "foo" (string)
```,
  caption: [Truthy and Falsy values in Weaver.],
) <fig:truthy-falsy>

This is often used to provide default values:
```weaver
greet := |name| echo("Hello " + (name || "unknown"));
greet("John") // Hello John
greet()       // Hello unknown
```

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
  caption: [Javascript type coercion. AVOID!],
) <fig:js-type-coercion>

Instead, Weaver has a stricter type system, which means there are no implicit type conversions, instead you should convert types manually, as shown in @fig:weaver-type-conversion.

#figure(
```weaver
number(true) + number(false)  == 1
12 / number("6")           == 2
"foo" + string(15 + 3)  == "foo18"
```,
  caption: [Weaver type conversion.],
) <fig:weaver-type-conversion>

This makes reading the code much easier, Since what you see is what you get.

If operators are used with incorrect types, Weaver will error at runtime.

#figure(
```weaver
// illegal operands bool + bool
true + false

// illegal operands int / string
12 / "6"

// illegal operands string + int
"foo" + 15 + 3

// illegal operands dict + array
{} + []

// illegal operands array == bool
[] == ![]

// illegal operands int == string
0 == "0"

// illegal operands int == array
0 == []

// illegal operands string == array
"0" == []
```,
  caption: [Weaver operators with incorrect type operands will error.],
) <fig:weaver-errors>

== Functions

Functions are first class, they can be passed around and used as values.

There are no special syntax for functions, you just assign a function value to a variable and call it.

#figure(
```weaver
// long form
add := |a, b| {
    return a + b;
}

add(1, 2)  // 3

// short form (auto-return)
add := |a, b| a + b

add(1, 2)  // 3
```,
  caption: [Defining and calling functions in Weaver.],
) <fig:functions>

This allows for expressive and concise code.

#figure(
```weaver
// imperative style
arr := [1, 2, 3, 4]
evenNumbers := []
for i := 0; i < len(arr); i++ {
    if arr[i] % 2 == 0 {
        evenNumbers |> push(arr[i])
    }
}
echo(evenNumbers) // [2, 4]

// functional style
arr := [1, 2, 3, 4]
evenNumbers := filter(arr, |n| n % 2 == 0)
echo(evenNumbers) // [2, 4]

// functional and piped
[1, 2, 3, 4]
    |> filter(|n| n % 2 == 0)
    |> echo() // [2, 4]
```,
  caption: [Different ways to filter a list in Weaver.],
) <fig:filtering-list>

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
  caption: [If-else statement in Weaver.],
) <fig:if-else>

Weaver also supports a ternary operator for concise conditional expressions:

#figure(
```weaver
n := 1
what := n % 2 == 0 ? "even" | "odd"
echo(what) // "odd"
```,
  caption: [Ternary operator in Weaver.],
) <fig:ternary>


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

Weaver also supports a `for-in` loop for iterating over ranges and collections:

#figure(
```weaver
// prints 0 to 4
for i in 0..4 {
    echo(i);
}

// prints each element of the array
arr := ["a", "b", "c"]
for item in arr {
    echo(item)
}
```,
  caption: [For-in loop in Weaver.],
) <fig:for-in>


#set page(columns: 1)

== Pattern Matching <sec:pattern-matching>

Weaver provides powerful pattern matching capabilities that allow developers to write expressive and concise code. Pattern matching can be used to match values against specific patterns and execute corresponding code blocks. The syntax follows the form:

```weaver
match expression {
  pattern1 => { /* code to execute if pattern1 matches */ },
  pattern2 => { /* code to execute if pattern2 matches */ },
  _ => { /* default case if no patterns match */ }
}
```

Match cases are evaluated in order, from top to bottom, until a match is found.

=== Basic Patterns

Weaver supports matching against literal values including integers, floats, strings, and booleans:

```weaver
x := "foo"
match x {
    "bar" => echo("bar is matched"),
    "nor" => echo("nor is matched"),
    "foo" => echo("finally foo is matched"),
    _ => echo("if nothing else matches"),
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
  [a, b, c] => echo(a + b + c), // 6
  _ => echo("no match")
}

// Object destructuring
match {name: "Alice", age: 30} {
  {name: n, age: a} => echo(n + " is " + string(a)), // Alice is 30
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

=== Overview of Patterns

#figure(
```weaver
// Match Statement Overview
match x {
    // matches string "foo"
    "foo" => {},
    // matches number 123
    123 => {},
    // matches number 1.4
    1.4 => {},
    // matches any number between 0 and 10
    0..10 => {},
    // matches any number less or equal to 10
    ..10 => {},
    // matches any number greater or equal to 5
    5.. => {},
    // matches array with at least two elements, where each element matches the pattern
    [<pattern>, <pattern>] => {},
    // matches object with "key" matching the pattern and "other" matching the pattern
    { key: <pattern>, other: <pattern> } => {},
    // matches any string and puts it in the variable s
    string(s) => {},
    // matches any number and puts it in the variable n
    number(n) => {},
    // matches error with the first pattern for the error message and the second pattern for the error details
    error(<pattern>, <pattern>) => {},
    // matches any value and puts it in the variable foo
    foo => {},
}
```,
  caption: [Overview of available patterns in Weaver.],
) <fig:patterns-overview>

Pattern matching in Weaver provides a powerful way to handle complex conditional logic in a readable and maintainable way. The combination of literal matching, type matching, destructuring, and guards makes it suitable for a wide range of use cases.


== Error Handling

Error handling is a critical aspect of any programming language, and Weaver provides a robust system for creating, propagating, and handling errors.

=== Creating and Propagating Errors

In Weaver, errors are values, just like numbers or strings. When a function raises an error, it's *automatically propagated* up the call stack unless explicitly handled. This is different from languages like JavaScript that use exceptions and `try/catch` blocks.

Errors are created and propagated using the `raise` keyword with the `error` function.

#figure(
```weaver
// Example: Automatic error propagation
divide := |a, b| {
    if b == 0 {
      raise error("Division by zero", {divisor: b})
    }

    return a / b
}

result := divide(10, 0)
echo(result) // This line will NOT execute
```,
  caption: [Creating and propagating errors in Weaver.],
) <fig:creating-propagating-errors>

=== Handling Errors

You can opt-out of automatic propagation using the `try` keyword. By adding `try` before an expression (like a function call), we tell Weaver that we want to handle the potential error ourselves. The result of the expression will be the error value if one was raised, or the successful result otherwise.

#figure(
```weaver
// Example: Opting out of automatic propagation
divide := |a, b| {
    if b == 0 {
      raise error("Division by zero", {divisor: b})
    }
    return a / b
}

result := try divide(10, 0)
echo("This line WILL execute")
echo(result) // Prints the error value
```,
  caption: [Handling errors with `try` in Weaver.],
) <fig:handling-errors-try>

We can then use pattern matching to handle the error.

#figure(
```weaver
result := try divide(10, 0)
match result {
    error(msg, data) => {
        echo("Error: " + msg)                    // Prints "Division by zero"
        echo("Divisor: " + string(data.divisor)) // Prints "Divisor: 0"
    },
    n => echo("Result: " + string(n)), // This won't execute in this case
}
```,
  caption: [Pattern matching on errors in Weaver.],
) <fig:error-pattern-matching>

Here's a more realistic example, fetching data from a URL:

#figure(
```weaver
response := try http.get("https://example.com/api/data")
match response {
    error(msg, { statusCode }) => {
        echo("HTTP request failed: " + msg)
        echo("Status code: " + string(statusCode))
    },
    res => {
        echo("Response body:")
        echo(res.body)
    }
}
```,
  caption: [Handling HTTP errors in Weaver.],
) <fig:http-error-handling>

This approach makes error handling explicit and integrates seamlessly with Weaver's pattern matching.

You can also use "truthyness" to handle errors, since errors are "falsy" values.

For example using `try/catch` in JavaScript:

```javascript
let user = null;
try {
    user = await fetchUser();
} catch (error) {
    user = { name: "Unknown" };
}
```

You can write the same thing like this in Weaver:

```weaver
user := try fetchUser() || { name: "Unknown" }
```

This is a very common pattern in Weaver to provide default values in case of an error.

=== Fibers

Weaver is a multi-threaded language, with support for non-blocking I/O operations.

You may have noticed that in the example, there is no `await` or `async` keyword, that's because Weaver is built on the notion of "Green Threads", which are lightweight threads that are managed by the runtime, and handles non-blocking IO without explicit yield points with `await`.

#figure(
```javascript
// JavaScript
const buyItem = async (itemID, userID, discountID) => {
    const item = await getItem(itemID);
    const itemWithDiscount = await applyDiscount(item, discount);
    const user = await getUser(userID);
    await pay(user, itemWithDiscount);
    return itemWithDiscount;
}
```,
  caption: [Asynchronous code in JavaScript with `async/await`.],
) <fig:async-js>

#figure(
```weaver
// Weaver
buyItem := |itemID, userID, discountID| {
    item := getItem(itemID) |> applyDiscount(discountID)
    user := getUser(userID)
    pay(user, item)
    return item
}
```,
  caption: [Concurrent code in Weaver without `async/await`.],
) <fig:concurrent-weaver>

Fibers consists of: Instructions (bytecode), Data (variables and constants), and the instruction pointer (IP).

Every code in the program is running inside a fiber, code in the top level is running on the main fiber, and you can start new fibers that run concurrently with the main fiber.

#figure(
```weaver
fiber := import("fiber")
io := import("io")

paths := ["foo.txt", "bar.txt"]

f1 := fiber.run(|| io.readFile(paths[0]))
f2 := fiber.run(|| io.readFile(paths[1]))

echo("main fiber")

files := fiber.wait([f1, f2])

echo(files[0]) // foo.txt
echo(files[1]) // bar.txt

// output:
// main fiber
// contents of foo.txt
// contents of bar.txt
```,
  caption: [Running I/O operations concurrently with fibers.],
) <fig:fibers-io>

== Benchmarks

A simple HTTP server with one endpoint reading a JSON file and returning a user by ID.

#figure(
```weaver
// Weaver
http := import("http");
io := import("io");
json := import("json");

router := http.newRouter();

router.get("/user/{id}", |req| {
    id := req.getParam("id");
    users := io.readFile("./main.json") 
             |> json.parse();
    user := users |> find(|user| {
      return user.id == number(id)
    });

    return user;
});

echo("starting server on port 8080");
http.listenAndServe(":8080", router);
```,
  caption: [HTTP server in Weaver.],
) <fig:http-server-weaver>

#figure(
```javascript
// JavaScript (Node.js with Express)
import express from "express";
import fs from "fs/promises";

const app = express();

app.get("/user/:id", async (req, res) => {
  const { id } = req.params
  const usersFile = await fs.readFile("./main.json")
  const users = JSON.parse(usersFile.toString())
  const user = users.find((u) => {
    return u.id === Number(id)
  });
  res.json(user)
});

console.log("Server running on port 3001");
app.listen(3001);
```,
  caption: [HTTP server in JavaScript.],
) <fig:http-server-js>

The benchmarks were run on a Lenovo Legion 5 pro with Ryzen 5 5800H (16 cores) and 32GB RAM.

#figure(
  image("./assets/http-bench/summary_comparison.svg", width: 100%),
  caption: [Benchmark summary comparison.],
) <fig:summary-comparison>

#figure(
  image("./assets/http-bench/memory_comparison.svg", width: 100%),
  caption: [Memory usage comparison.],
) <fig:memory-comparison>

#figure(
  image("./assets/http-bench/mean_latency_comparison.svg", width: 100%),
  caption: [Mean latency comparison.],
) <fig:mean-latency-comparison>

#figure(
  image("./assets/http-bench/p95_latency_comparison.svg", width: 100%),
  caption: [P95 latency comparison.],
) <fig:p95-latency-comparison>

#figure(
  image("./assets/http-bench/fibers-vs-node.svg", width: 100%),
  caption: [Fibers vs Node.js process model.],
) <fig:fibers-vs-node>

Weaver is multi-threaded, and for each request it creates a new fiber, so it's using all the cores within a single process, and fibers share the same memory space, so the memory usage is much lower. On the other hand, to utilize all cores in Node.js, we use `pm2` to run the server in cluster mode, which creates a separate process for each core.
