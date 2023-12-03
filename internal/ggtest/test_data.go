package ggtest

import (
	"io"
	"strconv"
	"time"
)

type testPrivStruct struct {
	constInt       int
	privBool       bool
	privInt        int
	privInt8       int8
	privInt64      int64
	privUint       uint
	privUint8      uint8
	privUint64     uint64
	privUintptr    uintptr
	privFloat64    float64
	privFloat32    float32
	privComplex128 complex128
	privComplex64  complex64
	privArray      [3]int
	privChan       chan int
	privInterface  io.Reader
	privAny        any
	privMap        map[int]any
	privMap1       *map[int]any
	privPointer    *int
	privSlice      []int
	privSlice2     []map[uint32]any
	privSlice3     []map[string]int32
	privSlice4     []map[string]string
	privString     string
	privString1    *string
	privStruct     struct{ i int }
	privTime       time.Time
}

// NewTestPrivStruct create a new testPrivStruct object.
func NewTestPrivStruct(i int) testPrivStruct {
	s := strconv.Itoa(i)

	return testPrivStruct{
		constInt:       1,
		privBool:       true,
		privInt:        i,
		privInt8:       int8(i),
		privInt64:      int64(i),
		privUint:       uint(i),
		privUint8:      uint8(i),
		privUint64:     uint64(i),
		privUintptr:    uintptr(i),
		privFloat64:    float64(i),
		privFloat32:    float32(i),
		privComplex128: complex(float64(i), float64(i)),
		privComplex64:  complex64(complex(float64(i), float64(i))),
		privArray:      [3]int{i - 1, i, i + 1},
		privChan:       make(chan int),
		privInterface:  nil,
		privAny:        nil,
		privMap:        map[int]any{i - 1: i - 1, i: i, i + 1: i + 1},
		privMap1:       &map[int]any{i - 1: i - 1, i: i, i + 1: i + 1},
		privPointer:    &i,
		privSlice:      []int{i - 1, i, i + 1},
		privSlice2:     []map[uint32]any{{uint32(i - 1): i - 1}, {uint32(i): i}, {uint32(i + 1): i + 1}},
		privSlice3:     []map[string]int32{{"中文" + strconv.Itoa(i-1): int32(i - 1)}, {"中文" + strconv.Itoa(i): int32(i)}, {"中文" + strconv.Itoa(i+1): int32(i + 1)}},
		privSlice4:     []map[string]string{{strconv.Itoa(i - 1): "中文 " + strconv.Itoa(i-1)}, {strconv.Itoa(i): "中文 " + strconv.Itoa(i)}, {strconv.Itoa(i + 1): "中文 " + strconv.Itoa(i+1)}},
		privString:     s,
		privString1:    &s,
		privStruct:     struct{ i int }{i: i},
		privTime:       time.Unix(1065430000, 0),
	}
}

func (s testPrivStruct) GetPrivInt() int {
	return s.privInt
}
func (s testPrivStruct) GetPrivChan() chan int {
	return s.privChan
}
func (s testPrivStruct) GetPrivInterface() io.Reader {
	return s.privInterface
}
func (s testPrivStruct) GetPrivAny() any {
	return s.privAny
}
func (s testPrivStruct) GetPrivStruct() any {
	return s.privStruct
}

// TestPubStruct is for testing.
type TestPubStruct struct {
	PubAny      any
	PubString   string
	PubString2  **string
	PrivStruct  testPrivStruct
	PrivStruct2 **testPrivStruct
	privAny     any
	privString  string
	privString2 **string
	privStruct  testPrivStruct
	privStruct2 **testPrivStruct
	privMap     map[string]any
	privMap2    **map[string]any
	privStructs []testPrivStruct
}

// NewTestPubStruct create a new TestPubStruct object.
func NewTestPubStruct(i int, PubAny any, PubString string, PrivStruct testPrivStruct, privAny any, privString string,
	privStruct testPrivStruct, privMap map[string]any) TestPubStruct {

	pstring := &PubString
	pstruct := &PrivStruct
	pmap := &privMap

	return TestPubStruct{
		PubAny:      PubAny,
		PubString:   PubString,
		PubString2:  &pstring,
		PrivStruct:  PrivStruct,
		PrivStruct2: &pstruct,
		privAny:     privAny,
		privString:  privString,
		privString2: &pstring,
		privStruct:  privStruct,
		privStruct2: &pstruct,
		privMap:     privMap,
		privMap2:    &pmap,
		privStructs: []testPrivStruct{NewTestPrivStruct(i - 1), privStruct, NewTestPrivStruct(i + 1)},
	}
}

// GetPrivAny returns privAny.
func (s TestPubStruct) GetPrivAny() any {
	return s.privAny
}

// GetPrivStruct returns privStruct.
func (s TestPubStruct) GetPrivStruct() testPrivStruct {
	return s.privStruct
}

// GetPrivStructs returns privStructs.
func (s TestPubStruct) GetPrivStructs() []testPrivStruct {
	return s.privStructs
}
