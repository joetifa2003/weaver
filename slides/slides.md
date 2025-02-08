---
# You can also start simply with 'default'
theme: default
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
# some information about your slides (markdown enabled)
title: Welcome to Slidev
author: Joe Tifa
info: |
  ## Slidev Starter Template
  Presentation slides for developers.

  Learn more at [Sli.dev](https://sli.dev)
# apply unocss classes to the current slide
# https://sli.dev/features/drawing
drawings:
  persist: false
# slide transition: https://sli.dev/guide/animations.html#slide-transitions
transition: slide-left
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
fonts:
    sans: 'DM Serif Text'
lineNumbers: true
remoteAssets: true
colorSchema: dark
addons:
  - slidev-component-spotlight
---

# Weaver

Simple scripting language for the joy of coding.

```rust
"Hello World!"
    |> echo()
```

---

<Toc />

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
1.23                    // float
true                    // boolean
false                   // boolean
[1, "2", [3]]           // array
{"a": 1, "b": 2}        // dictionary
{a: 1, b: 2}            // dictionary
|a, b| a + b            // lambda
|a, b| { return a + b } // lambda
```

---

## Binary

```weaver
1 + 2             // 3 (int)
1.0 + 2           // 3.0 (float)
2.3 + 3.4         // 5.7 (float)
"hello" + "world" // "helloworld" (string)
```

There is a distinction between integers and floats in Weaver, unlike many other scripting languages.

Also there is no type coercion, so you must be explicit about the conversion of types.

This is a deliberate desing decision to avoid mistakes of other languages, like the enfamous javascript.

```js
true + false   == 1
12 / "6"       == 2
"foo" + 15 + 3 == "foo153"
{} + []        == 0
[] == ![]      == true
```
