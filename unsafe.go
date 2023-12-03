package goget

import (
	"reflect"
	_ "unsafe"
)

//go:linkname valueInterface reflect.valueInterface
func valueInterface(v reflect.Value, safe bool) any
