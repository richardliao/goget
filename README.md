# Overview

Goget is a Go package that provides an easy way to get any sub-element from any Go object.

```go
// Get the city where the person lives.
goget.String(person, "address, city")
```

## Performance Implication

goget is reflection-heavy code, the performance will be much slower than the native one.

See benchmark for details.

## Requirement

Go v1.18 or above is required. 

# Getting Started

## Installation

```shell
go get github.com/richardliao/goget
```

## Option

* (N)one: Try best to get sub-element.
* (C)ase: Match keys case-sensitive, otherwise match keys case-insensitive.
* (S)afe: Skip unexported fields of a struct, otherwise search for unexported fields of a struct.
* (T)ype: Strict match target type, otherwise try to convert result to target type.

### Note About Option Safe

The ability to get unexported fields of a struct is a sharp knife. 

Always use option Safe unless necessary.

When the sub-element is an unexported field (private field), there are two situations:

1. The struct is defined in your package: your package can access it, but goget cannot.
2. The struct is defined outside your package: neither your package nor goget can access it.

When goget searches for a field, it has no way of knowing whether the struct is within a user package.

In case one, if you need to access the unexported fields of the package, do not use the option Safe.

## Usage 

Check example_test.go for more usage.

Example: Get a sub-element from a Go object.

```go
package main

import (
	"fmt"
	"github.com/richardliao/goget"
	"time"
)

type Person struct {
	Name    string
	Join   time.Time
	Age     uint
	isStaff bool
	address *Address
	tags    []any
	meta    map[string]any
}

type Address struct {
	Country string
	City    string
	street  string
	owner   *Person
}

func main() {
	addr := Address{
		Country: "Malawi",
		City:    "Mesa",
		street:  "123 Main St",
	}

	person := &Person{
		Name:    "Vin Mars",
		Join:    time.Unix(1065430000, 0).In(time.UTC),
		Age:     30,
		isStaff: true,
		address: &addr,
		tags:    []any{"tag1", "tag2", map[string]any{"d": "d"}, addr},
		meta:    map[string]any{"a": "a", "b": 1, "c": nil, "addr": addr, "e,f": "e,f"},
	}

	addr.owner = person

	// Best-effort search for Name, matching keys are case-insensitive, matching unexported fields of the structure
	// and converting to the target type on a best-effort basis..
	fmt.Println(goget.String(person, "Name")) // Vin Mars

	// Must will panic when error.
	fmt.Println(goget.MustString(person, goget.Safe, "Name")) // Vin Mars
	// May ignore any errors, always return target object.
	fmt.Println(goget.MayInt(person, goget.Safe, "Join,ext"))              // 0
	fmt.Println(goget.MustInt(person, goget.None, "Join,ext"))             // 63201026800
	fmt.Println(goget.MustInt(person, goget.None, "address, owner", "age")) // 30
	fmt.Println(goget.MustBool(person, goget.Case|goget.Type, "isStaff"))   // true
	// Support generic slice and map.
	fmt.Println(goget.MustSlice[any](person, goget.None, "tags"))            // [tag1 tag2 map[d:d] {Malawi Mesa 123 Main St <nil>}]
	fmt.Println(goget.MustMap[string, any](person, goget.None, "meta,addr")) // map[City:Mesa Country:Malawi]

	// Try best to get sub-element.
	street := goget.Any(person, "tags,City=Mesa,street")
	if s, ok := street.(string); ok {
		fmt.Println("street:", s) // 123 Main St
	} else {
		fmt.Println("get street type not match")
	}

	// Skip unexported fields of the struct.
	ext, err := goget.IntResult(person, goget.Safe, "Join,ext")
	if err != nil {
		fmt.Println("get ext error:", err) // QueryError[1]: [struct] field ext not exported
	} else {
		fmt.Println("ext:", ext)
	}

	// Strict match target type.
	isStaff, err := goget.IntResult(person, goget.Type, "isStaff")
	if err != nil {
		fmt.Println("get isStaff error:", err) // QueryError[2]: result type not match: need int got bool
	} else {
		fmt.Println("isStaff:", isStaff)
	}
}

```

## Error

QueryError code:

* ErrNotFound: Not found by paths.
* ErrTypeMatch: Target type not match.

## Why Need This

```go
package main

import (
	"fmt"
	"github.com/richardliao/goget"
)

type Person struct {
	Address *Address
}

type Address struct {
	Meta any
}

func main() {
	addr := Address{
		Meta: map[string]any{"tags": []any{"tag1", "tag2"}},
	}

	person := &Person{
		Address: &addr,
	}

	var tag2 string

	// Naive and danger way.
	tag2 = person.Address.Meta.(map[string]any)["tags"].([]any)[1].(string)
	fmt.Println("tag2:", tag2)

	// Safe but annoying way.
	if person.Address != nil {
		if person.Address.Meta != nil {
			if meta, ok := person.Address.Meta.(map[string]any); ok {
				if tags, ok := meta["tags"].([]any); ok {
					if len(tags) > 1 {
						if _tag, ok := tags[1].(string); ok {
							tag2 = _tag
							fmt.Println("tag2:", tag2)
						}
					}
				}
			}
		}
	}

	// goget way.
	tag2 = goget.MayString(person, goget.S, "Address, Meta, tags, 1")
	fmt.Println("tag2:", tag2)
}
```

## Benchmark

```
BenchmarkMap
BenchmarkMap-8            	 1808322	       623.2 ns/op
BenchmarkStruct
BenchmarkStruct-8         	 1000000	      1045 ns/op
BenchmarkSlice
BenchmarkSlice-8          	  469048	      2425 ns/op
BenchmarkMapNative
BenchmarkMapNative-8      	100000000	        11.59 ns/op
BenchmarkStructNative
BenchmarkStructNative-8   	1000000000	         0.6238 ns/op
BenchmarkSliceNative
BenchmarkSliceNative-8    	1000000000	         0.7844 ns/op
```
