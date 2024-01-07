package goget_test

import (
	gg "github.com/richardliao/goget"
	"testing"
	"time"
)

var (
	benchPersonMap = map[string]any{
		"name": "Vin Mars",
		"age":  30,
		"address": map[string]any{
			"country": "Malawi",
			"city":    "Mesa",
			"street":  "123 Main St",
		},
		"tags": []any{"tag1", "tag2", map[string]any{"d": "d"}},
		"meta": map[string]any{"a": "a", "b": 1, "c": nil},
	}

	benchAddr = Address{
		Country: "Malawi",
		City:    "Mesa",
		street:  "123 Main St",
	}

	benchPerson = &Person{
		Name:    "Vin Mars",
		Join:    time.Unix(1065430000, 0).In(time.UTC),
		Age:     30,
		isStaff: true,
		address: &benchAddr,
		tags:    []any{"tag1", "tag2", map[string]any{"d": "d"}, benchAddr},
		meta:    map[string]any{"a": "a", "b": 1, "c": nil, "addr": benchAddr},
	}
)

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.StringResult(benchPersonMap, gg.N, "address,city")
	}
}

func BenchmarkStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.StringResult(benchPerson, gg.N, "address,city")
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gg.StringResult(benchPerson, gg.N, "tags,City=Mesa,street")
	}
}

func BenchmarkMapNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchPersonMap["address"].(map[string]any)["city"]
	}
}

func BenchmarkStructNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchPerson.address.Country
	}
}

func BenchmarkSliceNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = benchPerson.tags[3].(Address).street
	}
}
