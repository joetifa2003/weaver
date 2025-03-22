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

# Overview

Weaver is a simple scripting language that prioritizes readability and simplicity.

Everything is plain objects and functions, most of the code is just concatenating functions together to solve problems.

````md magic-move {lines: true}
```weaver
arr := [1, 2, 3, 4]

echo(len(filter(arr, |n| n % 2 == 0))) // 2
```

```weaver
arr := [1, 2, 3, 4]

arr
    |> filter(|n| n % 2 == 0)
    |> len()
    |> echo() // 2
```
````

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
{a: 1, b: 2}            // object
|a, b| a + b            // lambda
|a, b| { return a + b } // lambda
nil                     // nil (null)
```

There is only one type in Weaver that indicates the absence of a value, it is `nil`.

---
layout: cover
---

## Binary Operations

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
```

Binary operations are very familiar to other languages, with highlight be operators like pipe operator (`|>`), and lazy evaluation of binray and `&&` and `||` operators that works for booleans and other values also based on if they are "truthy" or not.

---

## Truthy Values

Weaver boolean operators work with boolean expressions `true` and `false` as well as any other value in the language.

Values that are considered "falsey" are `nil`, `error`, `false`.

<v-click>

````md magic-move {lines: true}
```weaver
if (true)  { echo("true is truthy!") }
if (false) { echo("this will not run") }
if (nil)   { echo("this will not run") }
if (0)     { echo("zero is truthy!") }
if ("")    { echo("any string is truthy!") }

// Or operator stops at the first truthy value
nil || true  // true (boolean)
nil || "foo" // "foo" (boolean)

// And operator stops at the first falsy value
false && "foo" // false (boolean)
nil   && true  // nil   (boolean)
true  && nil   // nil   (boolean)
```

```weaver
greet := |name| echo("Hello " + name);
```

```weaver
greet := |name| echo("Hello " + name);
greet("John") // Hello John
```

```weaver
greet := |name| echo("Hello " + name);
greet("John") // Hello John
greet() // Hello nil
```

```weaver
// What if we want to display default name?
// no one is named "nil" :)
greet := |name| echo("Hello " + name);
greet("John") // Hello John
greet() // Hello nil
```

```weaver
greet := |name| echo("Hello " + (name || "unknown"));
greet("John") // Hello John
greet() // Hello unknown
```

```weaver
// What if we want to return error if name is not provided?
greet := |name| echo("Hello " + (name || "unknown"));
greet("John") // Hello John
greet() // Hello unknown
```

```weaver
greet := |name| {
    if (!name) { return error("name is required"); }
    echo("Hello " + name);
};
greet("John") // Hello John
greet()       // error: name is required
```

```weaver
// Can we do better?
greet := |name| {
    if (!name) { return error("name is required"); }
    echo("Hello " + name);
};
greet("John") // Hello John
greet()       // error: name is required
```

```weaver
greet := |name| {
    name || return error("name is required");
    echo("Hello " + name);
};
greet("John") // Hello John
greet()       // error: name is required
```

```weaver
greet := |name| {
    // Yes, return can be used as an expression
    name || return error("name is required");
    echo("Hello " + name);
};
greet("John") // Hello John
greet()       // error: name is required
```

````

</v-click>

---

## Type Coercion

Also there is no type coercion, so you must be explicit about the conversion of types, This is a deliberate design decision to avoid mistakes of other languages, like the enfamous javascript examples below.

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

Functions are the core of the language, they are "first class", that is they can be passed around and used as values.

There are no special syntax for functions, you just assign a function value to a variable and call it.

````md magic-move {lines: true}

```weaver
add := |a, b| {
    return a + b;
}
```

```weaver
add := |a, b| a + b
```

```weaver
add := |a, b| a + b
add(1, 2)  // 3
```

```weaver
add := |a, b| a + b
add(1)  // error: illegal operands number + nil
```

````

<v-click>

This allows for expressive and concise code, which will be otherwise very verbose.

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
arr := [1, 2, 3, 4]
evenNumbers := filter(arr, |n| n % 2 == 0)
echo(evenNumbers) // [2, 4]
```

```weaver
arr := [1, 2, 3, 4]
echo(filter(arr, |n| n % 2 == 0)) // [2, 4]
```

```weaver
echo(filter([1, 2, 3, 4], |n| n % 2 == 0)) // [2, 4]
```

```weaver
filter([1, 2, 3, 4], |n| n % 2 == 0) |> echo() // [2, 4]
```

```weaver
[1, 2, 3, 4] |> filter(|n| n % 2 == 0) |> echo() // [2, 4]
```

```weaver
[1, 2, 3, 4] 
    |> filter(|n| n % 2 == 0) 
    |> echo() // [2, 4]
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
what := n % 2 == 0 ? "even" : "odd"
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

```weaver
match [1, 2, 3, 4] {
    // matches an array which has at least 4 elements, and binds the first and last elements to a and b
    [a, _, __, b] => echo(a + b) // a=1; b=4; 5,
}
```

<v-click>

```weaver
n := 15
match n {
    0..10 => echo("n is between 0 and 10"),
    11..20 => echo("n is between 11 and 20"),
    _ => echo("n is greater than 20"),
}
```

</v-click>

<v-click>

```weaver
p := { 
    name: "Youssef", 
    bank: { debt: 100, cash: 10} 
};

match p {
    { bank: { debt: d, cash: c} } 
    if d >= 1 && d - c >= 0 => {
        echo("You have enough money to pay your debt");
    },

    _ => echo("You are still in debt"),
}
```

</v-click>

---

```weaver
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
    // matches array with two elements, where each element matches the pattern
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
```

---

### Match Patterns

Patterns can be as nested as you want.

```weaver
match x {
    [1, { someArray: [a, b, c] }] if a > b && b > c => echo("MATCH!"),
    _ => echo("NO MATCH!"),
}
```

### Match Guards

Match guards are a way to add additional conditions to a match case.

```weaver
match x {
    [..10, ..10] => {},
    // same as above
    [a, b] if a <= 10 && b <= 10 => {},
}
```

---

````md magic-move {lines: true}

```weaver
goodPartner := |p| {
    match p {
        { name: n, age: 18..30, isMarried: false } => return true,
        _ => return false,
    }
}
```

```weaver
goodPartner := |p| {
    match p {
        { name: n, age: 18..30, isMarried: false } => return true,
        _ => return false,
    }
}

goodPartner({ 
    name: "John",
    age: 25,
    isMarried: true,
}) // false

```

```weaver
goodPartner := |p| {
    match p {
        { name: n, age: 18..30, isMarried: false } => return true,
        _ => return false,
    }
}

goodPartner({ 
    name: "Youssef",
    age: 21,
    isMarried: false,
}) // true

```

```weaver
goodPartner := |p| {
    match p {
        { name: n, age: 18..30, isMarried: false } => return true,
        _ => return false,
    }
}

goodPartner({ 
    name: "Mahmoud",
    age: 45,
    isMarried: true,
}) // false
```

```weaver
goodPartner := |p| {
    if (!p.name) { return false }
    if (!p.age || p.age < 18 || p.age > 30) { return false }
    if (!p.isMarried) { return false }
    return false
}

goodPartner({ 
    name: "Mahmoud",
    age: 45,
    isMarried: true,
}) // false
```

````
