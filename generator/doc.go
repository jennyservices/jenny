// Package generator generates code from Service Definitions
//
// Jenny's approach to the code generation pipeline consists of 2 steps, decode
// and encode.
//
// IR representations exist for the purposee of converting a meta description to
// code. Service description defined here is a super-set of many meta languages;
// for instance consumes and produces are Swagger concepts that are highly tied
// to HTTP package, while parameter order is required for gRPC
//
//   graphQL schema ↘            ↗ js
//    swagger {2,3} → serviceIR  → go
//             gRPC ↗            ↘ swift
//
package generator
