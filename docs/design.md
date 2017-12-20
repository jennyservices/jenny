---
id: design
title: Design
sidebar_label: Overview
---

_Jenny_ is designed with one thing in mind, **production**!

It achieves production readiness with two main principles;

1. Spec 1st
2. Prod for dev.

**Spec 1st**: When we talk about Spec 1st, we don't talk about a particular
technology like `swagger` or `gRPC`, we talk about an interface that people
involved in development can work with. At **Typeform**, this is `swagger`, for
Google it's `gRPC` (amongst others), for your startup it could be a XML Service
spec you had lying around from `SOAP` days, only requirement for a service spec
is that it's a machine parse-able data format. We use this parse-able format to
generate everything from documentation to server and client code. This ensures
that everything we do is consistent and well documented.

**Prod for dev**: Jenny's follows modern "cloud" development conventions and
provides services that you are likely going to find in modern production
environments.

For example, `debug` mode does hot-reloading for your services, it uses your
services `/_health` endpoint to switch between the services during hot-reload.
Jenny also provides `tracing` capabilities in the debug console.

Design decisions and more are discussed in individual components sections
