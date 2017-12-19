---
id: generator
title: Generator
sidebar_label: Generator
---

## Design

Generator takes a service definition and generates code for it. For instance
given a swagger definition the code generation pipeline takes the swagger spec
and converts it in to a Intermediary Representation (`IR` for short). ServiceIR
is then consumed by a text template that is specific for a language like `Go` or
`Javascript`.

### Why IR?

The [IR](https://godoc.org/github.com/Typeform/jenny/generator/internal/ir/)
package has two higher level goals;

* Encapsulate multiple different service definitions.
* Ensure stability of code generation.

```
   graphQL Schema ⬂            ⬀ js
    swagger {2,3} ⇨ serviceIR  ⇨ go
             gRPC ⬀            ⬂ swift
```

### How does it work?

A pipeline consists of 3 steps;

1. Decode: A decoder such as swagger decoder, takes a service definition written
   in that language and converts it to `ir.Service`
2. Generator takes the `ir.Service` and does a couple of passes to make sure it
   has everything it needs.
3. An encoder such as the `Go` or `Javascript` encoders takes the `ir.Service`
   and uses it in it's templates to generate the code.

## Encoders and Decoders

Currently Jenny can read `swagger-2` definitions and spit out `Go` code that
uses `go-kit`. There are plans for supporting more languages and service
definitions but the initial set is limited to `Go` and `swagger`.

### Encoders

* [Go](https://github.com/Typeform/jenny/tree/master/generator/golang)
* [Javascript](https://github.com/Typeform/jenny/tree/master/generator/js) [WIP]

### Decoders

* [Swagger](https://github.com/Typeform/jenny/tree/master/generator/swagger)
