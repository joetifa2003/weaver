---
# You can also start simply with 'default'
theme: ./theme
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
# some information about your slides (markdown enabled)
title: Weaver - A Simple Scripting Language
author: Joe Tifa
# apply unocss classes to the current slide
# https://sli.dev/features/drawing
drawings:
  persist: false
# slide transition: https://sli.dev/guide/animations.html#slide-transitions
transition: fade
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
fonts:
  sans: 'Roboto'
  serif: 'DM Serif Text'
lineNumbers: true
remoteAssets: true
colorSchema: dark
addons:
  - slidev-component-spotlight
  - slidev-addon-excalidraw
---

# Weaver

Simple scripting language for the joy of coding.

```weaver
"Hello World!"
    |> echo()
```

---
layout: cover
---

# Syntax

Basic syntax and feal of the language.

In this section we will cover the basics of the language, and some of the idioms that make it unique.

---
layout: cover
---

## Values

```weaver
"Hello World!"          // string
123                     // number
1.23                    // number
true                    // boolean
false                   // boolean
[1, "2", [3]]           // array
{"a": 1, "b": 2}        // object
{a: 1, b: 2}            // object (shorthand)
|a, b| { return a + b } // lambda
nil                     // nil (null)
```

There is only one type in Weaver that indicates the absence of a value, it is `nil`.

*Note: `{a: 1, b: 2}` is shorthand for `{"a": 1, "b": 2}`.*

---

## Truthy Values

Weaver boolean operators work with boolean expressions `true` and `false` as well as any other value in the language.

Values that are considered "falsey" are `nil`, `error`, `false`.

```weaver
// Examples of truthy and falsy values
if (true)  { echo("true is truthy!") }  // Prints
if (0)     { echo("0 is truthy!") }     // Prints
if ("")    { echo("'' is truthy!") }    // Prints
if (nil)   { echo("nil is truthy!") }   // Does not print
if (false) { echo("false is truthy!") } // Does not print

// Or operator (||) returns the first truthy value
nil || "foo"  // "foo" (string)
true || "foo" // true (boolean)

// And operator (&&) returns the first falsy value
false && "foo" // false (boolean)
true && "foo"  // "foo" (string)
```

```weaver
greet := |name| {
    echo("Hello " + (name || "unknown"))
};
greet("John") // Hello John
greet() // Hello unknown
```

---

## Type Coercion

Weaver is dynamically typed, like the other scripting languages, but it's *strongly typed*.

There is no type coercion, so you must be explicit about the conversion of types, This is a deliberate design decision to avoid mistakes of other languages, like the enfamous JavaScript examples below.

````md magic-move
```js
// Javascript
true + false   == 1
12 / "6"       == 2
"foo" + 15 + 3 == "foo153"
{} + []        == 0
[] == ![]      == true
0 == "0"       == true
0 == []        == true
"0" == []      == false // !!!
```

```weaver
// Weaver
true + false   // ERROR! illegal operands bool + bool
12 / "6"       // ERROR! illegal operands int / string
"foo" + 15 + 3 // ERROR! illegal operands string + int
{} + []        // ERROR! illegal operands dict + array
[] == ![]      // ERROR! illegal operands array == bool
0 == "0"       // ERROR! illegal operands int == string
0 == []        // ERROR! illegal operands int == array
"0" == []      // ERROR! illegal operands string == array
```

```weaver
// Weaver
int(true) + int(false)  == 1
12 / int("6")           == 12
"foo" + string(15 + 3)  == "foo18"
```

````

<v-click>

In other words: What you see is what you get.

</v-click>

---

## Functions

Functions are first class, they can be passed around and used as values.

There are no special syntax for functions, you just assign a function value to a variable and call it.

````md magic-move {lines: true}
```weaver
add := |a, b| {
    return a + b;
}

add(1, 2)   // 3
1 |> add(2) // 3 
1 |> add(2) |> add(3) // 6
```

```weaver
add := |a, b| a + b

add(1, 2)  // 3
1 |> add(2) // 3 
1 |> add(2) |> add(3) // 6
```

````

<v-click>

This allows for expressive and concise code.

````md magic-move {lines: true}
```weaver
arr := [1, 2, 3, 4]

evenNumbers := []
for i := 0; i < len(arr); i++ {
    if (arr[i] % 2 == 0) {
        evenNumbers |> push(arr[i])
    }
}

echo(evenNumbers) // [2, 4]
```

```weaver
[1, 2, 3, 4] |> filter(|n| n % 2 == 0) |> echo() // [2, 4]
```
````

</v-click>

---

## Control Flow

### If Statement

```weaver
if (true) {
    echo("true is truthy!")
}

if (false) {
    echo("this will not run")
}

arr := [1, 2, 3, 4]
if (len(arr) >= 4) {
    arr[0] + arr[3] |> echo() // 5
}
```

### Ternary Operator

```weaver
n := 1
what := n % 2 == 0 ? "even" | "odd"
echo(what) // "even"
```

---

### Loops

```weaver
// prints 0 to 9
for i := 0; i < 10; i++ {
    echo(i);
}

// prints 0 to 9
for i in 0..9 {
    echo(i);
}

// prints 0 to 9
i := 0;
while (i < 10) {
    echo(i);
    i++;
}
```

---

### Match Statement

Pattern matching is a very powerful feature of Weaver, it allows you to write conditional logic based on the "shape" of the value.

Match cases are evaluated in order, from top to bottom, until a match is found.

```weaver
x := "foo"
match x {
    "bar" => echo("bar is matched"),
    "nor" => echo("nor is matched"),
    "foo" => echo("finally foo is matched"),
    _ => echo("if nothing else matches"),
}
```

```weaver
arr := [1, 2, 3, 4]
match arr {
    [1, 2] => {
        echo("arr starts with [1, 2]");
    },
    [2, 3] => {
        echo("arr starts with [2, 3]");
    },
    _ => {
        echo("otherwise");
    }
}
```

---

### Match Statement

```weaver
n := 15
match n {
    0..10 => echo("n is between 0 and 10"),
    11..20 => echo("n is between 11 and 20"),
    _ => echo("n is greater than 20"),
}
```

<v-click>

```weaver
students := [
    { name: "Youssef", gpa: 3.5 },
    { name: "John", gpa: 1.5 },
    { name: "Mahmoud", gpa: 2.0 },
];

for student in students {
    match student {
        { name, gpa: 0..1.5} => echo(name + " is good"),
        { name, gpa: 2..5} => echo(name + " is really good"),
    }
}

// output:
// Youssef is really good
// John is good
// Mahmoud is really good
```

</v-click>

---
layout: cover
---

### Match Guards

Match guards are a way to add additional conditions to a match case.

```weaver
match x {
    [..10, ..10] => {},
    // same as above
    [a, b] if a <= 10 && b <= 10 => {},
}
```

### Match Patterns

Patterns can be as nested as you want.

```weaver
match x {
    [1, { someArray: [a, b, c] }] if a > b && b > c => echo("MATCH!"),
    _ => echo("NO MATCH!"),
}
```

---

## Error Handling

Errors are values, just like numbers or strings. When a function raises an error, it's *automatically propagated* up the call stack unless explicitly handled. This is different from languages like JavaScript that use exceptions and `try/catch` blocks.

```weaver
// Example: Automatic error propagation
divide := |a, b| {
    if (b == 0) {
      raise error("Division by zero", {divisor: b});
    }

    return a / b;
};

result := divide(10, 0)
echo(result) // This line will NOT execute
```

You can opt-out of automatic propagation using the `try` keyword:

```weaver
// Example: Opting out of automatic propagation
result := try divide(10, 0)
echo("This line WILL execute")
echo(result) // Prints the error value
```

By adding `try` before the expression (function call), we tell Weaver that we want to handle the potential error ourselves. `result` will now contain the error value.

---

### Error Handling, Pattern Matching

```weaver
result := try divide(10, 0)
match result {
    error(msg, data) => {
        echo("Error: " + msg);                    // Prints "Division by zero"
        echo("Divisor: " + string(data.divisor)); // Prints "Divisor: 0"
    },
    n => echo("Result: " + string(n)), // This won't execute in this case
}
```

### Error Handling, Truthyness

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

---

### Fibers

Weaver is a multi-threaded language, with support for non-blocking I/O operations.

You may have noticed that in the example, there is no `await` or `async` keyword, that's because Weaver is built on the notion of "Green Threads", which are lightweight threads that are managed by the runtime, and handles non-blocking IO without explicit yield points with `await`.

```javascript [JavaScript]
const buyItem = async (itemID, userID, discountID) => {
    const item = await getItem(itemID);
    const itemWithDiscount = await applyDiscount(item, discount);
    const user = await getUser(userID);
    await pay(user, itemWithDiscount);
    return itemWithDiscount;
}
```

```weaver [Weaver]
buyItem := |itemID, userID, discountID| {
    item := getItem(itemID) |> applyDiscount(discountID)
    user := getUser(userID)
    pay(user, item)
    return item
}
```

---

### Fibers

Fibers consists of: Instructions (bytecode), Data (variables and constants), and the instruction pointer (IP).

Every code in the program is running inside a fiber, code in the top level is running on the main fiber, and you can start new fibers that run concurrently with the main fiber.

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
```

---

### Fibers In Action, Benchmarks

A simple HTTP server with one endpoint reading a JSON file and returning a user by ID.

<div class="grid grid-cols-2 gap-2">
```weaver [Weaver]
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
```

```javascript [JavaScript]
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
```
</div>

---
layout: full
---

### Fibers In Action, Benchmarks

Ran on a Lenovo Legion 5 pro with Ryzen 5 5800H (16 cores) and 32GB RAM.

<img src="./assets/http-bench/summary_comparison.svg" class="h-[90%] mx-auto" />

---
layout: full
---

### Fibers In Action, Benchmarks

Ran on a Lenovo Legion 5 pro with Ryzen 5 5800H (16 cores) and 32GB RAM.

<img src="./assets/http-bench/memory_comparison.svg" class="h-[90%] mx-auto" />

---
layout: full
---

### Fibers In Action, Benchmarks

Ran on a Lenovo Legion 5 pro with Ryzen 5 5800H (16 cores) and 32GB RAM.

<img src="./assets/http-bench/mean_latency_comparison.svg" class="h-[90%] mx-auto" />

---
clicks: 18
---

<Fibers />

---

### Comparing HTTP Servers

<img src="./assets/http-bench/fibers-vs-node.svg" class="h-[90%] mx-auto" />

---

### Standard Library

One of weaver goals is to have a comprehensive standard library, with a focus on being easy to use and easy to learn.

| Module | Description |
| ------ | ----------- |
| http   | HTTP server and client |
| html   | Writing HTML inside Weaver |
| io     | File system and networking and more |
| json   | JSON parser |
| math   | Math library |
| time   | Time library |
| raylib | Raylib (Graphics library) bindings |

---
layout: cover
---

## Thanks

Weaver is an ambitious project, and will continue to grow and improve over time.

Thanks for listening to the presentation.
