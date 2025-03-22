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

echo(evenNumbers) // [3, 5]
```

```weaver
arr := [1, 2, 3, 4]
evenNumbers := filter(arr, |n| n % 2 == 0)
echo(evenNumbers) // [3, 5]
```

```weaver
arr := [1, 2, 3, 4]
echo(filter(arr, |n| n % 2 == 0)) // [3, 5]
```

```weaver
echo(filter([1, 2, 3, 4], |n| n % 2 == 0)) // [3, 5]
```

```weaver
filter([1, 2, 3, 4], |n| n % 2 == 0) |> echo() // [3, 5]
```

```weaver
[1, 2, 3, 4] |> filter(|n| n % 2 == 0) |> echo() // [3, 5]
```

```weaver
[1, 2, 3, 4] 
    |> filter(|n| n % 2 == 0) 
    |> echo() // [3, 5]
```


````

</v-click>
