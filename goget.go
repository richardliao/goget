// Package goget provides an easy and simple way to get any sub-element from any Go object.
package goget

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// QueryError error code
type ErrCode int

const (
	ErrNotFound  ErrCode = 1 // Not found by paths
	ErrTypeMatch ErrCode = 2 // Target type not match
)

const (
	None Option = 0         // Try best
	Case Option = 1 << iota // Match keys case-sensitive
	Safe                    // Skip unexported fields of a struct
	Type                    // Strict match target type

	N Option = None
	C Option = Case
	S Option = Safe
	T Option = Type

	CS  = C | S
	CT  = C | T
	ST  = S | T
	CST = C | S | T
)

// Option
type Option uint8

type Result struct {
	val reflect.Value
	err *QueryError
}

// Any best-effort search for elements of an object by path and returns the result element as any.
// Best-effort means matching keys are case-insensitive, matching unexported fields of the structure
// and converting to the target type on a best-effort basis.
func Any(obj any, paths ...string) any {
	return MayAny(obj, None, paths...)
}

// AnyDefault like [Any], but returns default on error.
func AnyDefault(obj any, defaultVal any, paths ...string) any {
	r, err := AnyResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// AnyResult like [Any], but with option and returns error.
// If any error is returned, we can check if it is ErrNotFound or ErrTypeMatch.
func AnyResult(obj any, opt Option, paths ...string) (_ any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[any](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayAny like [AnyResult], but ignores errors.
func MayAny(obj any, opt Option, paths ...string) any {
	r, err := AnyResult(obj, opt, paths...)
	if err != nil {
		return nil
	}

	return r
}

// MustAny like [AnyResult], but panics on error.
func MustAny(obj any, opt Option, paths ...string) any {
	r, err := AnyResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// String like [Any], but returns string.
func String(obj any, paths ...string) string {
	return MayString(obj, None, paths...)
}

// StringDefault like [String], but returns default on error.
func StringDefault(obj any, defaultVal string, paths ...string) string {
	r, err := StringResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// StringResult like [AnyResult], but returns string.
func StringResult(obj any, opt Option, paths ...string) (_ string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[string](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayString like [MayAny], but returns string.
func MayString(obj any, opt Option, paths ...string) string {
	r, err := StringResult(obj, opt, paths...)
	if err != nil {
		return ""
	}

	return r
}

// MustString like [MustAny], but returns string.
func MustString(obj any, opt Option, paths ...string) string {
	r, err := StringResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Int like [Any], but returns int.
func Int(obj any, paths ...string) int {
	return MayInt(obj, None, paths...)
}

// IntDefault like [Int], but returns default on error.
func IntDefault(obj any, defaultVal int, paths ...string) int {
	r, err := IntResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// IntResult like [AnyResult], but returns int.
func IntResult(obj any, opt Option, paths ...string) (_ int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[int](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayInt like [MayAny], but returns int.
func MayInt(obj any, opt Option, paths ...string) int {
	r, err := IntResult(obj, opt, paths...)
	if err != nil {
		return 0
	}

	return r
}

// MustInt like [MustAny], but returns int.
func MustInt(obj any, opt Option, paths ...string) int {
	r, err := IntResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Uint like [Any], but returns uint.
func Uint(obj any, paths ...string) uint {
	return MayUint(obj, None, paths...)
}

// UintDefault like [Uint], but returns default on error.
func UintDefault(obj any, defaultVal uint, paths ...string) uint {
	r, err := UintResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// UintResult like [AnyResult], but returns uint.
func UintResult(obj any, opt Option, paths ...string) (_ uint, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[uint](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayUint like [MayAny], but returns uint.
func MayUint(obj any, opt Option, paths ...string) uint {
	r, err := UintResult(obj, opt, paths...)
	if err != nil {
		return 0
	}

	return r
}

// MustUint like [MustAny], but returns uint.
func MustUint(obj any, opt Option, paths ...string) uint {
	r, err := UintResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Complex like [Any], but returns complex128.
func Complex(obj any, paths ...string) complex128 {
	return MayComplex(obj, None, paths...)
}

// ComplexDefault like [Complex], but returns default on error.
func ComplexDefault(obj any, defaultVal complex128, paths ...string) complex128 {
	r, err := ComplexResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// ComplexResult like [AnyResult], but returns complex128.
func ComplexResult(obj any, opt Option, paths ...string) (_ complex128, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[complex128](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayComplex like [MayAny], but returns complex128.
func MayComplex(obj any, opt Option, paths ...string) complex128 {
	r, err := ComplexResult(obj, opt, paths...)
	if err != nil {
		return 0
	}

	return r
}

// MustComplex like [MustAny], but returns complex128.
func MustComplex(obj any, opt Option, paths ...string) complex128 {
	r, err := ComplexResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Float like [Any], but returns float64.
func Float(obj any, paths ...string) float64 {
	return MayFloat(obj, None, paths...)
}

// FloatDefault like [Float], but returns default on error.
func FloatDefault(obj any, defaultVal float64, paths ...string) float64 {
	r, err := FloatResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// FloatResult like [AnyResult], but returns float64.
func FloatResult(obj any, opt Option, paths ...string) (_ float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[float64](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayFloat like [MayAny], but returns float64.
func MayFloat(obj any, opt Option, paths ...string) float64 {
	r, err := FloatResult(obj, opt, paths...)
	if err != nil {
		return 0
	}

	return r
}

// MustFloat like [MustAny], but returns float64.
func MustFloat(obj any, opt Option, paths ...string) float64 {
	r, err := FloatResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Bool like [Any], but returns bool.
func Bool(obj any, paths ...string) bool {
	return MayBool(obj, None, paths...)
}

// BoolDefault like [Bool], but returns default on error.
func BoolDefault(obj any, defaultVal bool, paths ...string) bool {
	r, err := BoolResult(obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// BoolResult like [AnyResult], but returns bool.
func BoolResult(obj any, opt Option, paths ...string) (_ bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToAny[bool](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayBool like [MayAny], but returns bool.
func MayBool(obj any, opt Option, paths ...string) bool {
	r, err := BoolResult(obj, opt, paths...)
	if err != nil {
		return false
	}

	return r
}

// MustBool like [MustAny], but returns bool.
func MustBool(obj any, opt Option, paths ...string) bool {
	r, err := BoolResult(obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Slice like [Any], but returns slice.
func Slice[E any](obj any, paths ...string) []E {
	return MaySlice[E](obj, None, paths...)
}

// SliceDefault like [Slice], but returns default on error.
func SliceDefault[E any](obj any, defaultVal []E, paths ...string) []E {
	r, err := SliceResult[E](obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// SliceResult like [AnyResult], but returns slice.
func SliceResult[E any](obj any, opt Option, paths ...string) (_ []E, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToSlice[E](queryResult(obj, opt, paths), opt&Type == Type)
}

// MaySlice like [MayAny], but returns slice.
func MaySlice[E any](obj any, opt Option, paths ...string) []E {
	r, err := SliceResult[E](obj, opt, paths...)
	if err != nil {
		return nil
	}

	return r
}

// MustSlice like [MustAny], but returns slice.
func MustSlice[E any](obj any, opt Option, paths ...string) []E {
	r, err := SliceResult[E](obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

// Map like [Any], but returns map.
func Map[K comparable, E any](obj any, paths ...string) map[K]E {
	return MayMap[K, E](obj, None, paths...)
}

// MapDefault like [Map], but returns default on error.
func MapDefault[K comparable, E any](obj any, defaultVal map[K]E, paths ...string) map[K]E {
	r, err := MapResult[K, E](obj, None, paths...)
	if err != nil {
		return defaultVal
	}
	return r
}

// MapResult like [AnyResult], but returns map.
func MapResult[K comparable, E any](obj any, opt Option, paths ...string) (_ map[K]E, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = newQueryError(nil, ErrNotFound, "%v", r)
			return
		}
	}()

	return resultToMap[K, E](queryResult(obj, opt, paths), opt&Type == Type)
}

// MayMap like [MayAny], but returns map.
func MayMap[K comparable, E any](obj any, opt Option, paths ...string) map[K]E {
	r, err := MapResult[K, E](obj, opt, paths...)
	if err != nil {
		return nil
	}

	return r
}

// MustMap like [MustAny], but returns map.
func MustMap[K comparable, E any](obj any, opt Option, paths ...string) map[K]E {
	r, err := MapResult[K, E](obj, opt, paths...)
	if err != nil {
		panic(err)
	}

	return r
}

type QueryError struct {
	Code   ErrCode // ErrNotFound or ErrTypeMatch
	Detail string
	cause  error
}

func (w *QueryError) Error() string {
	if w == nil {
		return ""
	}

	msg := ""
	if w.cause != nil {
		msg = strings.TrimPrefix(w.cause.Error(), "QueryError: ")
	}

	if w.Detail != "" {
		if msg == "" {
			msg = w.Detail
		} else {
			msg += " -> " + w.Detail
		}
	}

	return fmt.Sprintf("QueryError[%d]: %s", w.Code, msg)
}

func (w *QueryError) Cause() error {
	if w == nil {
		return nil
	}
	return w.cause
}

func (w *QueryError) Unwrap() error {
	if w == nil {
		return nil
	}
	return w.cause
}

func newQueryError(cause error, code ErrCode, format string, args ...any) *QueryError {
	if cause != nil {
		var err *QueryError
		if errors.As(cause, &err) {
			return err
		}
	}

	return &QueryError{
		Code:   code,
		Detail: fmt.Sprintf(format, args...),
		cause:  cause,
	}
}

// queryResult search an object's elements by paths and returns the result element.
func queryResult(obj any, opt Option, paths []string) Result {
	if len(paths) == 0 {
		return Result{val: reflect.ValueOf(obj)}
	}

	value, err := query(reflect.ValueOf(obj), opt&Case == Case, opt&Safe == Safe, pathsToKeys(paths))
	if err != nil {
		return Result{err: err}
	}

	return Result{val: value}
}

// query search a Value by keys and returns the result element.
func query(value reflect.Value, caseSensitive, safe bool, keys []string) (_ reflect.Value, err *QueryError) {
	if !value.IsValid() {
		return value, newQueryError(nil, ErrNotFound, "invalid value")
	}

	if len(keys) == 0 {
		return value, nil
	}
	currentKey := keys[0]
	remainKeys := keys[1:]

	// Convert to concret element
	_value, err := toConcreteElem(value, safe, 0)
	if err != nil {
		return value, newQueryError(err, ErrNotFound, "invalid key: %s", currentKey)
	}
	value = _value

	switch value.Kind() {
	case reflect.Map:
		// Get key value
		keyValue := findMapKeyValue(value, currentKey, caseSensitive)
		fieldValue := value.MapIndex(keyValue)
		if !fieldValue.IsValid() {
			return value, newQueryError(nil, ErrNotFound, "[map] value not found by key %s", currentKey)
		}

		value, err = query(fieldValue, caseSensitive, safe, remainKeys)
		if err != nil {
			return value, newQueryError(err, ErrNotFound, "[map] query keys: %s", remainKeys)
		}
		return value, nil

	case reflect.Struct:
		fieldName := findStructFieldName(value, currentKey, caseSensitive)
		if safe {
			// Check unexported field when safe
			fieldType, ok := value.Type().FieldByName(fieldName)
			if !ok {
				return value, newQueryError(nil, ErrNotFound, "[struct] field %s not exists", currentKey)
			}
			if !fieldType.IsExported() {
				return value, newQueryError(nil, ErrNotFound, "[struct] field %s not exported", currentKey)
			}
		}

		fieldValue := value.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			return value, newQueryError(nil, ErrNotFound, "[struct] value not found by field %s", fieldName)
		}

		value, err = query(fieldValue, caseSensitive, safe, remainKeys)
		if err != nil {
			return value, newQueryError(err, ErrNotFound, "[struct] error query keys: %s", remainKeys)
		}
		return value, nil

	case reflect.Slice, reflect.Array:
		switch {
		case strings.Contains(currentKey, "="):
			parts := strings.SplitN(currentKey, "=", 2)
			if len(parts) != 2 {
				return value, newQueryError(nil, ErrNotFound, "[slice filter] invalid key format: %s", currentKey)
			}

			k, queryAttr := parts[0], parts[1]

			for index := 0; index < value.Len(); index++ {
				indexValue := value.Index(index)

				// Query the attribute value corresponding to k
				attrKeys := make([]string, 0)
				if k != "" {
					attrKeys = append(attrKeys, k)
				}
				attrVal, err := query(indexValue, caseSensitive, safe, attrKeys)
				if err != nil {
					continue
				}
				if !attrVal.IsValid() {
					continue
				}

				attrVal, err = toConcreteElem(attrVal, safe, 0)
				if err != nil {
					continue
				}

				// Determine whether the attribute meets the filter conditions
				indexAttr := valueToString(attrVal)
				if indexAttr != queryAttr {
					continue
				}

				// Query the remaining paths
				_value, err := query(indexValue, caseSensitive, safe, remainKeys)
				if err != nil {
					continue
				}
				value = _value

				return value, nil
			}

			return value, newQueryError(nil, ErrNotFound, "[slice filter] no elem by key: %s", currentKey)

		default:
			if strings.ToLower(currentKey) == "first" {
				currentKey = "0"
			}
			if strings.ToLower(currentKey) == "last" {
				currentKey = strconv.Itoa(value.Len() - 1)
			}

			// ConvErt: prevents changing the type of error
			index, convErr := strconv.Atoi(currentKey)
			if convErr != nil || index >= value.Len() || index < -value.Len() {
				return value, newQueryError(convErr, ErrNotFound, "[slice] invalid key: %s", currentKey)
			}

			if index < 0 {
				index = value.Len() + index
			}

			value, err = query(value.Index(index), caseSensitive, safe, remainKeys)
			if err != nil {
				return value, newQueryError(err, ErrNotFound, "[slice] query keys: %s", remainKeys)
			}
			return value, nil
		}
	default:
		return value, newQueryError(nil, ErrNotFound, "invalid kind: %s", value.Kind())
	}
}

// resultToAny convert a Result to target type.
func resultToAny[E any](result Result, typeStrict bool) (target E, err error) {
	if result.err != nil {
		return target, result.err
	}

	_v := valueToAny(result.val)

	if target, ok := _v.(E); ok {
		return target, nil
	}

	if typeStrict {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value := reflect.ValueOf(_v)
	if !value.IsValid() {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value, queryErr := toConcreteElem(value, false, 0)
	if queryErr != nil {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	// Convert value to target type
	var _vv any
	switch reflect.ValueOf(target).Kind() {
	case reflect.String:
		_vv = valueToString(value)
		return _vv.(E), nil
	case reflect.Int:
		_vv = valueToInt(value)
		return _vv.(E), nil
	case reflect.Uint:
		_vv = valueToUint(value)
		return _vv.(E), nil
	case reflect.Complex128:
		_vv = valueToComplex(value)
		return _vv.(E), nil
	case reflect.Float64:
		_vv = valueToFloat(value)
		return _vv.(E), nil
	case reflect.Bool:
		_vv = valueToBool(value)
		return _vv.(E), nil
	case reflect.Slice:
		_vv = valueToSlice[E](value)
		return _vv.(E), nil
	case reflect.Map:
		_vv = valueToMap[string, E](value)
		return _vv.(E), nil
	case reflect.Struct, reflect.Interface:
		if target, ok := _v.(E); !ok {
			return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
		} else {
			return target, nil
		}
	}

	return target, newQueryError(nil, ErrTypeMatch, "cannot convert result %v to %T", _v, target)
}

// resultToSlice convert a Result to target slice type.
func resultToSlice[E any](result Result, typeStrict bool) (target []E, err error) {
	if result.err != nil {
		return target, result.err
	}

	_v := valueToAny(result.val)

	if target, ok := _v.([]E); ok {
		return target, nil
	}

	if typeStrict {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value := reflect.ValueOf(_v)
	if !value.IsValid() {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value, queryErr := toConcreteElem(value, false, 0)
	if queryErr != nil {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	// Convert value to target type
	var _vv any
	switch reflect.ValueOf(target).Kind() {
	case reflect.Slice:
		_vv = valueToSlice[E](value)
		return _vv.([]E), nil
	}

	return target, newQueryError(nil, ErrTypeMatch, "cannot convert result %v to %T", _v, target)
}

// resultToMap convert a Result to target map type.
func resultToMap[K comparable, E any](result Result, typeStrict bool) (target map[K]E, err error) {
	if result.err != nil {
		return target, result.err
	}

	_v := valueToAny(result.val)

	if target, ok := _v.(map[K]E); ok {
		return target, nil
	}

	if typeStrict {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value := reflect.ValueOf(_v)
	if !value.IsValid() {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	value, queryErr := toConcreteElem(value, false, 0)
	if queryErr != nil {
		return target, newQueryError(nil, ErrTypeMatch, "result type not match: need %T got %T", target, _v)
	}

	// Convert value to target type
	switch reflect.ValueOf(target).Kind() {
	case reflect.Map:
		j, err := json.Marshal(valueToAny(value))
		if err != nil {
			return target, newQueryError(nil, ErrTypeMatch, "cannot convert result %v to %T", _v, target)
		}
		err = json.Unmarshal(j, &target)
		if err != nil {
			return target, newQueryError(nil, ErrTypeMatch, "cannot convert result %v to %T", _v, target)
		}
		return target, nil
	}

	return target, newQueryError(nil, ErrTypeMatch, "cannot convert result %v to %T", _v, target)
}

// valueToString convert a Value to string.
func valueToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	default:
		return fmt.Sprintf("%v", valueToAny(value))
	}
}

// valueToInt convert a Value to int.
// If the value cannot be converted to int, return 0.
func valueToInt(value reflect.Value) int {
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int(value.Uint())
	case reflect.Complex64, reflect.Complex128:
		return int(real(value.Complex()))
	case reflect.Float32, reflect.Float64:
		return int(value.Float())
	default:
		v, _ := strconv.ParseInt(fmt.Sprintf("%v", valueToAny(value)), 10, 64)
		return int(v)
	}
}

// valueToUint convert a Value to uint.
// If the value cannot be converted to uint, return 0.
func valueToUint(value reflect.Value) uint {
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uint(value.Uint())
	case reflect.Complex64, reflect.Complex128:
		return uint(real(value.Complex()))
	case reflect.Float32, reflect.Float64:
		return uint(value.Float())
	default:
		v, _ := strconv.ParseUint(fmt.Sprintf("%v", valueToAny(value)), 10, 64)
		return uint(v)
	}
}

// valueToComplex convert a Value to complex128.
// If the value cannot be converted to uint, return zero complex128.
func valueToComplex(value reflect.Value) complex128 {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return complex(float64(value.Int()), 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return complex(float64(value.Uint()), 0)
	case reflect.Complex64, reflect.Complex128:
		return value.Complex()
	case reflect.Float32, reflect.Float64:
		return complex(value.Float(), 0)
	default:
		v, _ := strconv.ParseComplex(fmt.Sprintf("%v", valueToAny(value)), 64)
		return v
	}
}

// valueToFloat convert a Value to float64.
// If the value cannot be converted to uint, return 0.
func valueToFloat(value reflect.Value) float64 {
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(value.Uint())
	case reflect.Complex64, reflect.Complex128:
		return float64(real(value.Complex()))
	case reflect.Float32, reflect.Float64:
		return value.Float()
	default:
		v, _ := strconv.ParseFloat(fmt.Sprintf("%v", valueToAny(value)), 64)
		return v
	}
}

// valueToBool convert a Value to bool.
// If the value cannot be converted to uint, return false.
func valueToBool(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Bool:
		return value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() != 0
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == complex(0, 0)
	case reflect.Float32, reflect.Float64:
		v := value.Float()
		return !math.IsNaN(v) && v != 0
	default:
		v, _ := strconv.ParseBool(fmt.Sprintf("%v", valueToAny(value)))
		return v
	}
}

// valueToSlice convert a Value to slice.
// If the value is not a slice or array, return nil.
func valueToSlice[E any](value reflect.Value) []E {
	if value.Kind() == reflect.Slice {
		if ss, ok := valueToAny(value).([]E); ok {
			return ss
		}
	} else if value.Kind() == reflect.Array {
		m := make([]E, value.Len())
		for i := 0; i < value.Len(); i++ {
			v, ok := valueToAny(value.Index(i)).(E)
			if !ok {
				return nil
			}
			m[i] = v
		}
		return m
	}

	return nil
}

// valueToMap convert a Value to map.
// If the value is not a map, return nil.
func valueToMap[K comparable, E any](value reflect.Value) map[K]E {
	if value.Kind() == reflect.Map {
		if m, ok := valueToAny(value).(map[K]E); ok {
			return m
		}
	}

	return nil
}

// stringToMapKeyType convert a string key to map key's type.
func stringToMapKeyType(key string, kind reflect.Type) reflect.Value {
	if !kind.Comparable() {
		return reflect.ValueOf(key)
	}

	value := reflect.ValueOf(key)

	val := reflect.New(kind)

	switch kind.Kind() {
	case reflect.Bool:
		val.Elem().SetBool(valueToBool(value))
	case reflect.String:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.Elem().SetInt(int64(valueToInt(value)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val.Elem().SetUint(uint64(valueToUint(value)))
	case reflect.Complex64, reflect.Complex128:
		val.Elem().SetComplex(valueToComplex(value))
	case reflect.Float32, reflect.Float64:
		val.Elem().SetFloat(valueToFloat(value))
	}

	return val.Elem()
}

// findStructFieldName try to find a map key by name(first match exactly, then try match caseinsensitive if specified).
// If not found, return key as is.
func findMapKeyValue(value reflect.Value, currentKey string, caseSensitive bool) (keyValue reflect.Value) {
	keyValue = reflect.ValueOf(currentKey)

	// Invalid, return as is
	if !value.IsValid() || value.IsNil() {
		return keyValue
	}

	// Map currentKey is not string, convert currentKey to map currentKey's type.
	if value.Type().Key().Kind() != reflect.String {
		return stringToMapKeyType(currentKey, value.Type().Key())
	}

	// Not caseinsensitive, return as is
	if caseSensitive {
		return keyValue
	}

	// First find currentKey exactly
	if value.MapIndex(keyValue).IsValid() {
		return keyValue
	}

	// Then search lower case
	for _, mapKeyValue := range value.MapKeys() {
		if strings.TrimSpace(strings.ToLower(mapKeyValue.String())) == strings.TrimSpace(strings.ToLower(currentKey)) {
			return mapKeyValue
		}
	}

	// Not found, return as is
	return keyValue
}

// findStructFieldName try to find a struct field by name(first match exactly, then try match caseinsensitive if specified).
// If not found, return key as is.
func findStructFieldName(value reflect.Value, currentKey string, caseSensitive bool) string {
	// Not caseinsensitive, return as is
	if caseSensitive {
		return currentKey
	}

	// First find currentKey exactly
	if value.FieldByName(currentKey).IsValid() {
		return currentKey
	}

	// Then search lower case
	for i := 0; i < value.NumField(); i++ {
		value.Type().Field(i).IsExported()
		fieldName := value.Type().Field(i).Name
		if strings.TrimSpace(strings.ToLower(fieldName)) == strings.TrimSpace(strings.ToLower(currentKey)) {
			return fieldName
		}
	}

	// Not found, return as is
	return currentKey
}

// pathsToKeys normalized paths to keys by splitting path with comma.
func pathsToKeys(paths []string) []string {
	keys := make([]string, 0)
	for _, path := range paths {
		for _, key := range splitWithEscape(path, ",", "\\") {
			keys = append(keys, strings.TrimSpace(key))
		}
	}

	return keys
}

// splitWithEscape like strings.Split with escape string supporting.
func splitWithEscape(s string, separator string, esc string) []string {
	a := strings.Split(s, separator)

	for i := len(a) - 2; i >= 0; i-- {
		if strings.HasSuffix(a[i], esc) {
			a[i] = a[i][:len(a[i])-len(esc)] + separator + a[i+1]
			a = append(a[:i+1], a[i+2:]...)
		}
	}
	return a
}

// valueToAny convert a Value to its current value(even unexported struct field) as an any.
func valueToAny(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}

	if value.CanInterface() {
		return value.Interface()
	} else {
		return valueInterface(value, false)
	}
}

// toConcreteElem convert an interface or pointer reflect.Value to its concret(underlying) element.
func toConcreteElem(value reflect.Value, safe bool, depth int) (reflect.Value, *QueryError) {
	if depth > 10 {
		return value, newQueryError(nil, ErrNotFound, "get concret elem exceed max depth %d", depth)
	}

	switch value.Kind() {
	case reflect.Invalid:
		return value, newQueryError(nil, ErrNotFound, "%s", value)

	case reflect.Interface:
		// Convert to the value's underlying value as an any.
		if value.CanInterface() {
			value = reflect.ValueOf(value.Interface())
		} else {
			if safe {
				return value, newQueryError(nil, ErrNotFound, "%s", value)
			}
			// When the Value is an unexported field (private field), two situations:
			// 1. The struct is defined within the user package: accessible to users, but inaccessible to goget.
			// 2. The struct is defined out of the user package: inaccessible to both users and goget.
			// But when goget searches for a field, it cannot know whether the struct is within the user package or not.
			// So, we have to access the private fields in all situations.
			value = reflect.ValueOf(valueInterface(value, false))
		}

	case reflect.Pointer:
		// Convert to the value that value points to.
		value = reflect.Indirect(value)
	}

	if value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		return toConcreteElem(value, safe, depth+1)
	}

	return value, nil
}
