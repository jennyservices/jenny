---
id: options
title: Options
sidebar_label: Options
---

## Design

Options package provides some opinionated middlewares for your `Go` service,
where possible it follows the
[`middleware`](https://godoc.org/github.com/go-kit/kit/endpoint#Middleware)
patterns established in `go-kit`.

There is a certain order which the options package would chain your middlewares

```text
Current ordering of the middlewares goes as follows

Request
	↓
Requests-ID (enabled by default)
	↓
Tracing (enabled noop by default)
	↓
Error reporting (enabled noop by default)
	↓
JWT parser (disabled by default, enable by passing WithJWTParser)
	↓
User parser (disabled by default, enable by passing WithUserParser); (this is useful for ratelimiting by user)
	↓
Scopes parser (disabled by default, enable by passing WithScopesParser)
```

Please see errors package for error reporting.
