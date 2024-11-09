---
theme: ./dracula.json 
---

# Weaver

```rust
~~~xargs cat
../examples/01_helloworld.w
~~~
```

Weaver is a simple, fast, and secure scripting language.

## Why Weaver?

- Simple to read and write
- Suitable for quick scripting
- Embeddable to any application

---

# Lexer 

```
~~~graph-easy --as=boxart
digraph {
    rankdir=LR;

    input [label="Input: 'x := 1'"];
    lexer [label="Lexical Analysis"];
    stream [label="{type=IDENT, lit='x'}\n {type=ASSIGN, lit=':='}\n {type=INT, lit='1'}"];
        
    input -> lexer;
    lexer -> stream;
    stream -> parser;
}
~~~
```

---

# Quick Tour
