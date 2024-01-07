package goget_test

import (
	"fmt"
	"github.com/richardliao/goget"
	"time"
)

type Person struct {
	Name    string
	Join    time.Time
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

func Example() {
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

	// Query simple objects.
	fmt.Println(goget.Any("string1"))      // string1
	fmt.Println(goget.Any("string1", "1")) // <nil>

	// Query compound objects.
	fmt.Println(goget.MustAny(person, goget.S, "Name"))                      // Vin Mars
	fmt.Println(goget.MayAny(person, goget.S, "name"))                       // Vin Mars
	fmt.Println(goget.AnyResult(person, goget.C, "name"))                    // <nil> QueryError[1]: [struct] value not found by field name
	fmt.Println(goget.MustString(person, goget.T, "name"))                   // Vin Mars
	fmt.Println(goget.IntResult(person, goget.C|goget.T, "name"))            // 0 QueryError[1]: [struct] value not found by field name
	fmt.Println(goget.IntResult(person, goget.S, "Name"))                    // 0 <nil>
	fmt.Println(goget.IntResult(person, goget.T, "Name"))                    // 0 QueryError[2]: result type not match: need int got string
	fmt.Println(goget.MustAny(person, goget.S, "Join"))                      // 2003-10-06 08:46:40 +0000 UTC
	fmt.Println(goget.MustAny(person, goget.N, "Join,ext"))                  // 63201026800
	fmt.Println(goget.AnyResult(person, goget.S, "Join,ext"))                // <nil> QueryError[1]: [struct] field ext not exported
	fmt.Println(goget.MustInt(person, goget.S, "Age"))                       // 30
	fmt.Println(goget.MustString(person, goget.S, "Age"))                    // 30
	fmt.Println(goget.MustBool(person, goget.N, "isStaff"))                  // true
	fmt.Println(goget.MustInt(person, goget.N, "isStaff"))                   // 1
	fmt.Println(goget.MustString(person, goget.N, "isStaff"))                // true
	fmt.Println(goget.MustString(person, goget.N, "address, Country"))       // Malawi
	fmt.Println(goget.MustString(person, goget.N, "address, street"))        // 123 Main St
	fmt.Println(goget.MustString(person, goget.N, "address, owner", "Name")) // Vin Mars
	fmt.Println(goget.MustString(person, goget.N, "tags"))                   // [tag1 tag2 map[d:d] {Malawi Mesa 123 Main St <nil>}]
	fmt.Println(goget.StringResult(person, goget.N, "tags,1"))               // tag2 <nil>
	fmt.Println(goget.StringResult(person, goget.N, "tags,-1"))              // {Malawi Mesa 123 Main St <nil>} <nil>
	fmt.Println(goget.AnyResult(person, goget.N, "tags,-5"))                 // <nil> QueryError[1]: [slice] invalid key: -5
	fmt.Println(goget.Any(person, "tags,-5"))                                // <nil>
	fmt.Println(goget.StringResult(person, goget.N, "tags,d=d"))             // map[d:d] <nil>
	fmt.Println(goget.StringResult(person, goget.N, "tags,City=Mesa"))       // {Malawi Mesa 123 Main St <nil>} <nil>
	fmt.Println(goget.StringResult(person, goget.N, "tags,=tag2"))           // tag2 <nil>
	fmt.Println(goget.Any(person, "meta,addr,Country"))                      // Malawi
	fmt.Println(goget.AnyResult(person, goget.N, "meta,addr"))               // {Malawi Mesa 123 Main St <nil>} <nil>
	fmt.Println(goget.AnyResult(person, goget.N, "meta,e\\,f"))              // e,f <nil>

	// Get generic targets.
	fmt.Println(goget.MustSlice[any](person, goget.N, "tags"))              // [tag1 tag2 map[d:d] {Malawi Mesa 123 Main St <nil>}]
	fmt.Println(goget.MapResult[string, any](person, goget.N, "meta,addr")) // map[City:Mesa Country:Malawi] <nil>
	fmt.Println(goget.MapResult[int, any](person, goget.N, "meta,addr"))    // map[] QueryError[2]: cannot convert result {Malawi Mesa 123 Main St <nil>} to map[int]interface {}

	// Output:
	// string1
	// <nil>
	// Vin Mars
	// Vin Mars
	// <nil> QueryError[1]: [struct] value not found by field name
	// Vin Mars
	// 0 QueryError[1]: [struct] value not found by field name
	// 0 <nil>
	// 0 QueryError[2]: result type not match: need int got string
	// 2003-10-06 08:46:40 +0000 UTC
	// 63201026800
	// <nil> QueryError[1]: [struct] field ext not exported
	// 30
	// 30
	// true
	// 1
	// true
	// Malawi
	// 123 Main St
	// Vin Mars
	// [tag1 tag2 map[d:d] {Malawi Mesa 123 Main St <nil>}]
	// tag2 <nil>
	// {Malawi Mesa 123 Main St <nil>} <nil>
	// <nil> QueryError[1]: [slice] invalid key: -5
	// <nil>
	// map[d:d] <nil>
	// {Malawi Mesa 123 Main St <nil>} <nil>
	// tag2 <nil>
	// Malawi
	// {Malawi Mesa 123 Main St <nil>} <nil>
	// e,f <nil>
	// [tag1 tag2 map[d:d] {Malawi Mesa 123 Main St <nil>}]
	// map[City:Mesa Country:Malawi] <nil>
	// map[] QueryError[2]: cannot convert result {Malawi Mesa 123 Main St <nil>} to map[int]interface {}
}
