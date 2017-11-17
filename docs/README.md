# Jenny

_Jenny the Generator_ is a toolkit developed by [Typeform](http://typeform.com)
for rapid development of production ready services.

## Design

_Jenny_ is designed with one thing in mind, **production**!

It achieves production readiness with two main principles;

1. Spec 1st
2. Prod for dev.

**Spec 1st**: When we talk about Spec 1st, we don't talk about a particular
technology like `swagger` or `gRPC`, we talk about an interface that people
involved in development can work with. For **Typeform**, this is `swagger`, for
Google it's `gRPC` (amongst others), for your startup it could be a XML Service
spec you had lying around from `SOAP` days, only requirement for a service spec
is that it's a machine parse-able data format. We use this parse-able format to
generate everything from documentation to server and client code. This ensures
that everything we do is consistent and well documented.

**Prod for dev**: Jenny's follows modern "cloud development" conventions and
provides services that you are likely going to find in modern production
environments.

For example, `debug` mode does hot-reloading for your services, it uses your
services `/_health` endpoint to switch between the services during hot-reload.
Jenny also provides `tracing` capabilities in the debug console.

Design decisions and more are discussed in individual components sections

* [Generator](generator.md)
* [Debug](debug.md) [WIP, Help Wanted!]
* [Auth](auth.md)
* [Options](options.md)
* [Healthy](healthy.md) [WIP]
* [Conventions and defaults](conventions.md)

## Roadmap

1. Open Source Jenny
2. Discuss Plug-ins Architecture
3. Refactor dashboard code with generated JS

## [Tutorials](tutorials)

* [User Service](tutorials/userservice)
