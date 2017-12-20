---
id: conventions
title: Conventions
sidebar_label: Conventions
---

A basic Jenny project has 3 packages;

1. `cmd` for keeping the main executable.
2. `transport` package where we put service definitions in their respective
   versions as their package names.
3. And a final `foo` as package that implements the `FooService` interface
   defined in the transport version.

```text
 .
 ├── cmd
 │   └── userservice
 │       └── main.go
 ├── transport
 │   └── v1
 │       ├── jenny.go
 │       └── swagger.yaml
 └── user
     └── user.go
```

Putting aside everything else the transport package with the version numbers is
propbably the best convention that we stumbled upon. The jenny cli app expects
this project layout however you can change this layout by passing flags in the
command line.
