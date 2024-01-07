package goget

import (
	"fmt"
	"github.com/richardliao/goget/internal/ggtest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func _toString(i any) string {
	value := reflect.ValueOf(i)
	if !value.IsValid() {
		return ""
	}
	value, queryErr := toConcreteElem(value, false, 0)
	if queryErr != nil {
		return ""
	}

	return valueToString(value)
}

func TestString(t *testing.T) {
	i := 1<<16 + 1
	s := strconv.Itoa(i)
	privMap := map[string]any{"A": i}
	privStruct := ggtest.NewTestPrivStruct(i)
	pubStruct := ggtest.NewTestPubStruct(i, "PubAny", "PubString", privStruct, "privAny", "privString", privStruct, privMap)

	m := map[string]*ggtest.TestPubStruct{
		"pubStruct": &pubStruct,
	}

	ss := []string{
		"string1",
		"string2",
	}

	tests := []struct {
		value  any
		opt    Option
		keys   []string
		expect string
		err    bool
	}{
		{nil, None, nil, "", false},
		{"string", None, nil, "string", false},

		{pubStruct, None, []string{"PubAny"}, _toString(pubStruct.PubAny), false},
		{pubStruct, None, []string{"PubAny2"}, "", true},
		{pubStruct, None, []string{"privStruct"}, _toString(pubStruct.GetPrivStruct()), false},
		{pubStruct, None, []string{"privStruct2", "privString1"}, s, false},
		{&pubStruct, None, []string{"privStruct2", "privString1"}, s, false},

		{pubStruct, Safe, []string{"privStruct"}, "", true},
		{pubStruct, Safe, []string{"privStruct", "privString"}, "", true},

		{pubStruct, None, []string{"privStruct", "privBool"}, _toString(pubStruct.GetPrivStruct().GetPrivInt() == i), false},
		{pubStruct, None, []string{"privStruct", "privInt"}, _toString(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, None, []string{"privStruct", "privInt8"}, _toString(int8(pubStruct.GetPrivStruct().GetPrivInt())), false},
		{pubStruct, None, []string{"privStruct", "privUint8"}, _toString(uint8(pubStruct.GetPrivStruct().GetPrivInt())), false},
		{pubStruct, None, []string{"privStruct", "privFloat32"}, _toString(float32(pubStruct.GetPrivStruct().GetPrivInt())), false},
		{pubStruct, None, []string{"privStruct", "privComplex64"}, _toString(complex64(complex(float64(i), float64(i)))), false},
		{pubStruct, None, []string{"privStruct", "privArray"}, _toString([3]int{i - 1, i, i + 1}), false},
		{pubStruct, None, []string{"privStruct", "privChan"}, _toString(pubStruct.GetPrivStruct().GetPrivChan()), false},
		{pubStruct, None, []string{"privStruct", "privInterface"}, _toString(pubStruct.GetPrivStruct().GetPrivInterface()), false},
		{pubStruct, None, []string{"privStruct", "privAny"}, _toString(pubStruct.GetPrivStruct().GetPrivAny()), false},
		{pubStruct, None, []string{"privStruct", "privMap"}, _toString(map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}), false},
		{pubStruct, None, []string{"privStruct", "privMap1"}, _toString(&map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}), false},
		{pubStruct, None, []string{"privStruct", "privPointer"}, _toString(i), false},
		{pubStruct, None, []string{"privStruct", "privSlice"}, _toString([]int{i - 1, i, i + 1}), false},
		{pubStruct, None, []string{"privStruct", "privString"}, s, false},
		{pubStruct, None, []string{"privStruct", "privString1"}, s, false},
		{pubStruct, None, []string{"privStruct", "privStruct"}, _toString(pubStruct.GetPrivStruct().GetPrivStruct()), false},
		{pubStruct, None, []string{"privStruct", "privTime"}, _toString(time.Unix(1065430000, 0)), false},

		{pubStruct, None, []string{"privstruct", "privint"}, _toString(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, Case, []string{"privstruct", "privInt"}, "", true},
		{pubStruct, Case, []string{"privStruct", "privint"}, "", true},
		{pubStruct, Safe, []string{"privStruct", "privint"}, "", true},

		{pubStruct, None, []string{"privStruct", "privSlice2", "1", strconv.Itoa(i)}, _toString(i), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i + 1)}, _toString(i + 1), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "-1", strconv.Itoa(i + 1)}, _toString(i + 1), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "last", strconv.Itoa(i + 1)}, _toString(i + 1), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, _toString(i), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, _toString(i + 1), false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i-1, i-1), strconv.Itoa(i - 1)}, _toString(i - 1), false},

		{pubStruct, None, []string{"privStruct", "privSlice3", fmt.Sprintf("中文%d=%d", i, i), "中文" + strconv.Itoa(i)}, _toString(int32(i)), false},
		{pubStruct, None, []string{"privStruct", "privSlice4", fmt.Sprintf("%d=中文 %d", i, i), strconv.Itoa(i)}, "中文 " + strconv.Itoa(i), false},

		{pubStruct, None, []string{"privStruct", "privSlice2", "-4", strconv.Itoa(i + 1)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", "a", strconv.Itoa(i + 1)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i+1), strconv.Itoa(i)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=a", i), strconv.Itoa(i)}, "", true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("a=%d", i), strconv.Itoa(i)}, "", true},

		{pubStruct, None, []string{"privStructs"}, _toString(pubStruct.GetPrivStructs()), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i)}, _toString(pubStruct.GetPrivStruct()), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, _toString(i), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privString=%d", i)}, _toString(pubStruct.GetPrivStruct()), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privString1=%d", i)}, _toString(pubStruct.GetPrivStruct()), false},
		{pubStruct, None, []string{"privStructs", "constInt=1", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, _toString(i + 1), false},

		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privInt=%d", i+2)}, "", true},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, "", true},

		{m, None, []string{"pubStruct", "privAny"}, _toString(pubStruct.GetPrivAny()), false},

		{ss, None, []string{"1"}, _toString(ss[1]), false},
		{ss, None, []string{"-1"}, _toString(ss[1]), false},
		{ss, None, []string{}, _toString(ss), false},
		{ss, None, []string{"2"}, "", true},
		{ss, None, []string{"-3"}, "", true},
		{ss, None, []string{""}, "", true},
		{ss, None, []string{"key=val"}, "", true},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got, err := StringResult(tt.value, tt.opt, tt.keys...)
		if tt.err {
			if err == nil {
				t.Fatalf("expect error got nil")
			}
			continue
		}

		if !assert.Equalf(tt.expect, got, "got: %+v", got) {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}

func TestAny(t *testing.T) {
	i := 1<<16 + 1
	s := strconv.Itoa(i)
	privMap := map[string]any{"A": i}
	privStruct := ggtest.NewTestPrivStruct(i)
	pubStruct := ggtest.NewTestPubStruct(i, "PubAny", "PubString", privStruct, "privAny", "privString", privStruct, privMap)

	m := map[string]*ggtest.TestPubStruct{
		"pubStruct": &pubStruct,
	}

	ss := []string{
		"string1",
		"string2",
	}

	tests := []struct {
		value  any
		opt    Option
		keys   []string
		expect any
		err    bool
	}{
		{nil, None, nil, nil, false},
		{"string", None, nil, "string", false},

		{pubStruct, None, []string{"PubAny"}, pubStruct.PubAny, false},
		{pubStruct, None, []string{"PubAny2"}, nil, true},
		{pubStruct, None, []string{"privStruct"}, pubStruct.GetPrivStruct(), false},
		{pubStruct, None, []string{"privStruct2", "privString1"}, &s, false},
		{&pubStruct, None, []string{"privStruct2", "privString1"}, &s, false},

		{pubStruct, Safe, []string{"privStruct"}, nil, true},
		{pubStruct, Safe, []string{"privStruct", "privString"}, nil, true},

		{pubStruct, None, []string{"privStruct", "privBool"}, pubStruct.GetPrivStruct().GetPrivInt() == i, false},
		{pubStruct, None, []string{"privStruct", "privInt"}, pubStruct.GetPrivStruct().GetPrivInt(), false},
		{pubStruct, None, []string{"privStruct", "privInt8"}, int8(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, None, []string{"privStruct", "privUint8"}, uint8(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, None, []string{"privStruct", "privFloat32"}, float32(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, None, []string{"privStruct", "privComplex64"}, complex64(complex(float64(i), float64(i))), false},
		{pubStruct, None, []string{"privStruct", "privArray"}, [3]int{i - 1, i, i + 1}, false},
		{pubStruct, None, []string{"privStruct", "privChan"}, pubStruct.GetPrivStruct().GetPrivChan(), false},
		{pubStruct, None, []string{"privStruct", "privInterface"}, pubStruct.GetPrivStruct().GetPrivInterface(), false},
		{pubStruct, None, []string{"privStruct", "privAny"}, pubStruct.GetPrivStruct().GetPrivAny(), false},
		{pubStruct, None, []string{"privStruct", "privMap"}, map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}, false},
		{pubStruct, None, []string{"privStruct", "privMap1"}, &map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}, false},
		{pubStruct, None, []string{"privStruct", "privPointer"}, &i, false},
		{pubStruct, None, []string{"privStruct", "privSlice"}, []int{i - 1, i, i + 1}, false},
		{pubStruct, None, []string{"privStruct", "privString"}, s, false},
		{pubStruct, None, []string{"privStruct", "privString1"}, &s, false},
		{pubStruct, None, []string{"privStruct", "privStruct"}, pubStruct.GetPrivStruct().GetPrivStruct(), false},
		{pubStruct, None, []string{"privStruct", "privTime"}, time.Unix(1065430000, 0), false},

		{pubStruct, None, []string{"privstruct", "privint"}, pubStruct.GetPrivStruct().GetPrivInt(), false},
		{pubStruct, Case, []string{"privstruct", "privInt"}, nil, true},
		{pubStruct, Case, []string{"privStruct", "privint"}, nil, true},
		{pubStruct, Safe, []string{"privStruct", "privint"}, nil, true},

		{pubStruct, None, []string{"privStruct", "privSlice2", "1", strconv.Itoa(i)}, i, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "-1", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", "last", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, i, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i-1, i-1), strconv.Itoa(i - 1)}, i - 1, false},

		{pubStruct, None, []string{"privStruct", "privSlice3", fmt.Sprintf("中文%d=%d", i, i), "中文" + strconv.Itoa(i)}, int32(i), false},
		{pubStruct, None, []string{"privStruct", "privSlice4", fmt.Sprintf("%d=中文 %d", i, i), strconv.Itoa(i)}, "中文 " + strconv.Itoa(i), false},

		{pubStruct, None, []string{"privStruct", "privSlice2", "-4", strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", "a", strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i+1), strconv.Itoa(i)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=a", i), strconv.Itoa(i)}, nil, true},
		{pubStruct, None, []string{"privStruct", "privSlice2", fmt.Sprintf("a=%d", i), strconv.Itoa(i)}, nil, true},

		{pubStruct, None, []string{"privStructs"}, pubStruct.GetPrivStructs(), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, i, false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privString=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privString1=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, None, []string{"privStructs", "constInt=1", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, i + 1, false},

		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privInt=%d", i+2)}, nil, true},
		{pubStruct, None, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, nil, true},

		{m, None, []string{"pubStruct", "privAny"}, pubStruct.GetPrivAny(), false},

		{ss, None, []string{"1"}, ss[1], false},
		{ss, None, []string{"-1"}, ss[1], false},
		{ss, None, []string{}, ss, false},
		{ss, None, []string{"2"}, nil, true},
		{ss, None, []string{"-3"}, nil, true},
		{ss, None, []string{""}, nil, true},
		{ss, None, []string{"key=val"}, nil, true},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got, err := AnyResult(tt.value, tt.opt, tt.keys...)
		if tt.err {
			if err == nil {
				t.Fatalf("expect error got nil")
			}
			continue
		}

		if !assert.Equalf(tt.expect, got, "got: %+v", got) {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}

func TestQuery(t *testing.T) {
	i := 1<<16 + 1
	s := strconv.Itoa(i)
	privMap := map[string]any{"A": i}
	privStruct := ggtest.NewTestPrivStruct(i)
	pubStruct := ggtest.NewTestPubStruct(i, "PubAny", "PubString", privStruct, "privAny", "privString", privStruct, privMap)

	m := map[string]*ggtest.TestPubStruct{
		"pubStruct": &pubStruct,
	}

	ss := []string{
		"string1",
		"string2",
	}

	tests := []struct {
		value         any
		caseSensitive bool
		safe          bool
		keys          []string
		expect        any
		err           bool
	}{
		{nil, false, false, nil, nil, false},
		{"string", false, false, nil, "string", false},

		{pubStruct, false, false, []string{"PubAny"}, pubStruct.PubAny, false},
		{pubStruct, false, false, []string{"PubAny2"}, nil, true},
		{pubStruct, false, false, []string{"privStruct"}, pubStruct.GetPrivStruct(), false},
		{pubStruct, false, false, []string{"privStruct2", "privString1"}, &s, false},
		{&pubStruct, false, false, []string{"privStruct2", "privString1"}, &s, false},

		{pubStruct, false, true, []string{"privStruct"}, nil, true},
		{pubStruct, false, true, []string{"privStruct", "privString"}, nil, true},

		{pubStruct, false, false, []string{"privStruct", "privBool"}, pubStruct.GetPrivStruct().GetPrivInt() == i, false},
		{pubStruct, false, false, []string{"privStruct", "privInt"}, pubStruct.GetPrivStruct().GetPrivInt(), false},
		{pubStruct, false, false, []string{"privStruct", "privInt8"}, int8(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, false, false, []string{"privStruct", "privUint8"}, uint8(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, false, false, []string{"privStruct", "privFloat32"}, float32(pubStruct.GetPrivStruct().GetPrivInt()), false},
		{pubStruct, false, false, []string{"privStruct", "privComplex64"}, complex64(complex(float64(i), float64(i))), false},
		{pubStruct, false, false, []string{"privStruct", "privArray"}, [3]int{i - 1, i, i + 1}, false},
		{pubStruct, false, false, []string{"privStruct", "privChan"}, pubStruct.GetPrivStruct().GetPrivChan(), false},
		{pubStruct, false, false, []string{"privStruct", "privInterface"}, pubStruct.GetPrivStruct().GetPrivInterface(), false},
		{pubStruct, false, false, []string{"privStruct", "privAny"}, pubStruct.GetPrivStruct().GetPrivAny(), false},
		{pubStruct, false, false, []string{"privStruct", "privMap"}, map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}, false},
		{pubStruct, false, false, []string{"privStruct", "privMap1"}, &map[int]any{i - 1: i - 1, i: i, i + 1: i + 1}, false},
		{pubStruct, false, false, []string{"privStruct", "privPointer"}, &i, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice"}, []int{i - 1, i, i + 1}, false},
		{pubStruct, false, false, []string{"privStruct", "privString"}, s, false},
		{pubStruct, false, false, []string{"privStruct", "privString1"}, &s, false},
		{pubStruct, false, false, []string{"privStruct", "privStruct"}, pubStruct.GetPrivStruct().GetPrivStruct(), false},
		{pubStruct, false, false, []string{"privStruct", "privTime"}, time.Unix(1065430000, 0), false},

		{pubStruct, false, false, []string{"privstruct", "privint"}, pubStruct.GetPrivStruct().GetPrivInt(), false},
		{pubStruct, true, false, []string{"privstruct", "privInt"}, nil, true},
		{pubStruct, true, false, []string{"privStruct", "privint"}, nil, true},
		{pubStruct, false, true, []string{"privStruct", "privint"}, nil, true},

		{pubStruct, false, false, []string{"privStruct", "privSlice2", "1", strconv.Itoa(i)}, i, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", "-1", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", "last", strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, i, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, i + 1, false},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i-1, i-1), strconv.Itoa(i - 1)}, i - 1, false},

		{pubStruct, false, false, []string{"privStruct", "privSlice3", fmt.Sprintf("中文%d=%d", i, i), "中文" + strconv.Itoa(i)}, int32(i), false},
		{pubStruct, false, false, []string{"privStruct", "privSlice4", fmt.Sprintf("%d=中文 %d", i, i), strconv.Itoa(i)}, "中文 " + strconv.Itoa(i), false},

		{pubStruct, false, false, []string{"privStruct", "privSlice2", "-4", strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", "a", strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", "2", strconv.Itoa(i)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i+1), strconv.Itoa(i)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("%d=a", i), strconv.Itoa(i)}, nil, true},
		{pubStruct, false, false, []string{"privStruct", "privSlice2", fmt.Sprintf("a=%d", i), strconv.Itoa(i)}, nil, true},

		{pubStruct, false, false, []string{"privStructs"}, pubStruct.GetPrivStructs(), false},
		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privint=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i)}, i, false},
		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privString=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privString1=%d", i)}, pubStruct.GetPrivStruct(), false},
		{pubStruct, false, false, []string{"privStructs", "constInt=1", "privSlice2", fmt.Sprintf("%d=%d", i+1, i+1), strconv.Itoa(i + 1)}, i + 1, false},

		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privInt=%d", i+2)}, nil, true},
		{pubStruct, false, false, []string{"privStructs", fmt.Sprintf("privint=%d", i), "privSlice2", fmt.Sprintf("%d=%d", i, i), strconv.Itoa(i + 1)}, nil, true},

		{m, false, false, []string{"pubStruct", "privAny"}, pubStruct.GetPrivAny(), false},

		{ss, false, false, []string{"1"}, ss[1], false},
		{ss, false, false, []string{"-1"}, ss[1], false},
		{ss, false, false, []string{}, ss, false},
		{ss, false, false, []string{"2"}, nil, true},
		{ss, false, false, []string{"-3"}, nil, true},
		{ss, false, false, []string{""}, nil, true},
		{ss, false, false, []string{"key=val"}, nil, true},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got, err := query(reflect.ValueOf(tt.value), tt.caseSensitive, tt.safe, tt.keys)
		if tt.err {
			if err == nil {
				t.Fatalf("expect error got nil")
			}
			continue
		}

		if !assert.Equalf(tt.expect, valueToAny(got), "got: %+v", got) {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}

func TestValueToAny(t *testing.T) {
	i := 1<<16 + 1
	privMap := map[string]any{"A": i}
	privStruct := ggtest.NewTestPrivStruct(i)
	pprivStruct := &privStruct
	pubStruct := ggtest.NewTestPubStruct(i, "PubAny", "PubString", privStruct, "privAny", "privString", privStruct, privMap)

	m := map[string]ggtest.TestPubStruct{
		"A": pubStruct,
	}

	tests := []struct {
		field  reflect.Value
		expect any
	}{
		{reflect.ValueOf(nil), nil},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("PubAny"), "PubAny"},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("PubString"), "PubString"},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("PrivStruct"), privStruct},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("privAny"), "privAny"},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("privString"), "privString"},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("privStruct"), privStruct},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("privStruct2"), &pprivStruct},
		{reflect.ValueOf(&pubStruct).Elem().FieldByName("privMap"), privMap},

		{reflect.ValueOf(pubStruct).FieldByName("PubString"), "PubString"},
		{reflect.ValueOf(pubStruct).FieldByName("PrivStruct"), privStruct},
		{reflect.ValueOf(pubStruct).FieldByName("PubAny"), "PubAny"},
		{reflect.ValueOf(pubStruct).FieldByName("privAny"), "privAny"},
		{reflect.ValueOf(pubStruct).FieldByName("privString"), "privString"},
		{reflect.ValueOf(pubStruct).FieldByName("privStruct2"), &pprivStruct},

		{reflect.ValueOf(m["A"]).FieldByName("privAny"), "privAny"},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := valueToAny(tt.field)
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestToConcreteElem(t *testing.T) {
	var a = 0
	var aa = &a
	var aaa = &aa
	var aaaa any = &aaa
	var b any = "b"
	tests := []struct {
		value  any
		safe   bool
		depth  int
		expect any
		err    bool
	}{
		{"A", false, 0, "A", false},
		{1, false, 0, 1, false},
		{complex(1, 1), false, 0, complex(1, 1), false},
		{2.0, false, 0, 2.0, false},
		{map[int]any{}, false, 0, map[int]any{}, false},
		{a, false, 0, a, false},
		{&a, false, 0, a, false},
		{&aa, false, 0, a, false},
		{&aaaa, false, 0, a, false},
		{b, false, 0, "b", false},
		// nil value
		{nil, false, 0, nil, true},
		// exceed depth
		{&aaaa, false, 7, nil, true},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got, err := toConcreteElem(reflect.ValueOf(tt.value), tt.safe, tt.depth)
		if tt.err {
			if err == nil {
				t.Fatalf("expect error got nil")
			}
			continue
		}

		if !assert.Equalf(tt.expect, got.Interface(), "got: %+v", got) {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	i := 1
	privMap := map[string]any{"A": i}
	privStruct := ggtest.NewTestPrivStruct(i)
	pubStruct := ggtest.NewTestPubStruct(i, "PubAny", "PubString", privStruct, "privAny", "privString", privStruct, privMap)

	reflect.ValueOf(pubStruct).Type().FieldByName("privStruct")

	got, _ := toConcreteElem(reflect.ValueOf(pubStruct).FieldByName("privAny"), false, 0)
	assert.Equalf("privAny", got.Interface(), "got: %+v", got)

	_, err := toConcreteElem(reflect.ValueOf(pubStruct).FieldByName("privAny"), true, 0)
	if err == nil {
		t.Fatalf("expect error got nil")
	}
}

func TestStringToMapKeyType(t *testing.T) {
	type A struct {
		A any
	}

	tests := []struct {
		currentKey string
		value      any
		expect     any
		panic      bool
	}{
		// nil value: panic
		{"a", nil, nil, true},
		{"a", new(map[string]any), nil, true},
		{"a", map[string]any{}, "a", false},
		{"a", map[int]any{}, 0, false},
		{"a", map[float32]any{}, float32(0), false},
		{"a", map[bool]any{}, false, false},
		{"a", map[uint8]any{}, uint8(0), false},
		{"a", map[complex128]any{}, complex(0, 0), false},
		{"a", map[A]any{}, A{}, false},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		if tt.panic {
			assert.Panics(func() { stringToMapKeyType(tt.currentKey, reflect.ValueOf(tt.value).Type().Key()) })
			continue
		}

		got := stringToMapKeyType(tt.currentKey, reflect.ValueOf(tt.value).Type().Key())
		assert.Equalf(reflect.ValueOf(tt.expect).Type(), got.Type(), "got: %+v", got)
		assert.Equalf(reflect.ValueOf(tt.expect).Interface(), got.Interface(), "got: %+v", got)
	}
}

func TestFindMapKeyValue(t *testing.T) {
	tests := []struct {
		value         any
		currentKey    string
		caseSensitive bool
		expect        any
		panic         bool
	}{
		// nil value: no panic
		{nil, "", false, "", false},
		// empty map
		{map[string]string{}, "", false, "", false},
		{map[string]string{}, "A", false, "A", false},
		// empty key
		{map[string]string{"A": ""}, "", false, "", false},
		// empty value
		{map[string]string{}, "A", false, "A", false},
		// match key case sensitive
		{map[string]string{"A": ""}, "A", true, "A", false},
		// match key case insensitive
		{map[string]string{"A": ""}, "a", false, "A", false},
		// match key but case not match
		{map[string]string{"A": ""}, "a", true, "a", false},
		// not found case sensitive
		{map[string]string{"A": ""}, "b", true, "b", false},
		// not found case insensitive
		{map[string]string{"A": ""}, "b", false, "b", false},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		if tt.panic {
			assert.Panics(func() { findMapKeyValue(reflect.ValueOf(tt.value), tt.currentKey, tt.caseSensitive) })
			continue
		}

		got := findMapKeyValue(reflect.ValueOf(tt.value), tt.currentKey, tt.caseSensitive)
		assert.Equalf(tt.expect, got.Interface(), "got: %+v", got)
	}
}

func TestFindStructFieldName(t *testing.T) {
	tests := []struct {
		value         any
		currentKey    string
		caseSensitive bool
		expect        string
		panic         bool
	}{
		// nil value: panic
		{nil, "", false, "", true},
		// empty struct
		{struct{}{}, "", false, "", false},
		{struct{}{}, "A", false, "A", false},
		// empty key
		{struct{ A string }{""}, "", false, "", false},
		// empty value
		{struct{ A string }{}, "A", false, "A", false},
		// match key case sensitive
		{struct{ A string }{""}, "A", true, "A", false},
		// match key case insensitive
		{struct{ A string }{""}, "a", false, "A", false},
		// match key but case not match
		{struct{ A string }{""}, "a", true, "a", false},
		// not found case sensitive
		{struct{ A string }{""}, "b", true, "b", false},
		// not found case insensitive
		{struct{ A string }{""}, "b", false, "b", false},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		if tt.panic {
			assert.Panics(func() { findStructFieldName(reflect.ValueOf(tt.value), tt.currentKey, tt.caseSensitive) })
			continue
		}

		got := findStructFieldName(reflect.ValueOf(tt.value), tt.currentKey, tt.caseSensitive)
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestSplitWithEscape(t *testing.T) {
	tests := []struct {
		s      string
		sep    string
		esc    string
		expect []string
	}{
		{"a,b=day hour,c", ",", "\\", []string{"a", "b=day hour", "c"}},
		{"a,b\\,c", ",", "\\", []string{"a", "b,c"}},
		{" a, b\\\\,c ", ",", "\\\\", []string{" a", " b,c "}},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := splitWithEscape(tt.s, tt.sep, tt.esc)
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestGetAny(t *testing.T) {
	personMap := map[string]any{
		"name": "Vin Mars",
		"age":  30,
		"address": map[string]any{
			"street":  "123 Main St",
			"City":    "Mesa",
			"country": "Malawi",
		},
		"tags": []string{"a", "b", "c"},
		"meta": map[string]string{"a": "a", "b": "b"},
	}

	assert := assert.New(t)

	got, err := AnyResult(personMap, None, "address", "City")
	assert.NoErrorf(err, "error: %v", err)
	if !assert.Equalf(got, got, "got: %+v", got) {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestAnyToString(t *testing.T) {
	tests := []struct {
		value  any
		expect any
	}{
		{nil, "<nil>"},
		{"string1", "string1"},
		{1, "1"},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := valueToString(reflect.ValueOf(tt.value))
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestAnyToInt(t *testing.T) {
	tests := []struct {
		value  any
		expect int
	}{
		{nil, 0},
		{"string1", 0},
		{1, 1},
		{"1", 1},
		{true, 1},
		{1.1, 1},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := valueToInt(reflect.ValueOf(tt.value))
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestAnyToSlice(t *testing.T) {
	tests := []struct {
		value  any
		expect []int
	}{
		{nil, nil},
		{"string1", nil},
		{[]int{1}, []int{1}},
		{[]string{"1"}, nil},
		{[]bool{true}, nil},
		{[]float64{1.1}, nil},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := valueToSlice[int](reflect.ValueOf(tt.value))
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}

func TestAnyToMap(t *testing.T) {
	tests := []struct {
		value  any
		expect map[int]int
	}{
		{nil, nil},
		{"string1", nil},
		{map[int]int{1: 1}, map[int]int{1: 1}},
		{map[string]string{"1": "1"}, nil},
		{map[bool]bool{true: true}, nil},
		{map[float64]float64{1.1: 1.1}, nil},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		got := valueToMap[int, int](reflect.ValueOf(tt.value))
		assert.Equalf(tt.expect, got, "got: %+v", got)
	}
}
