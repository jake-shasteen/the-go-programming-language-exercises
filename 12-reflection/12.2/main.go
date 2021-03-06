package main

import (
	"fmt"
	"reflect"
	"strconv"
)

// PtrLoopExample has a comment
type PtrLoopExample struct {
	name string
	ptr  *PtrLoopExample
}

const LoopsAllowed = 8

func main() {
	// cyclic data structures:

	// a struct with a pointer that points to itself
	loop := &PtrLoopExample{
		name: "loop"}
	loop.ptr = loop

	// loop5 := &PtrLoopExample{name: "loop5"}
	// loop4 := &PtrLoopExample{
	// 	name: "loop4",
	// 	ptr:  loop5}
	// loop3 := &PtrLoopExample{
	// 	name: "loop3",
	// 	ptr:  loop4}
	loop2 := &PtrLoopExample{
		name: "loop2"}
	loop1 := &PtrLoopExample{
		name: "loop1",
		ptr:  loop2}
	// loop5.ptr = loop1

	Display("loop", loop)

	Display("loop chain", loop1)
}

// Display has a comment
func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(LoopsAllowed, name, reflect.ValueOf(x))
}

// Make `display` safe to use on cyclic data structures by bounding the number of steps it takes before abandoning the recursion.
func display(count int, path string, v reflect.Value) {
	if count == 0 {
		fmt.Printf("Truncated: %s = %s\n", path, formatAtom(v))
		return
	}
	switch v.Kind() {
	case reflect.Invalid:
		// terminal
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			// non-terminal
			display(count-1, fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		length := v.NumField()
		for i := 0; i < length; i++ {
			// non-terminal
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(count-1, fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			// non-terminal
			display(count-1, fmt.Sprintf("%s[%s]", path, formatComparable(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			// terminal
			fmt.Printf("%s = nil\n", path)
		} else {
			// non-terminal
			display(count-1, fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			// terminal
			fmt.Printf("%s = nil\n", path)
		} else {
			// terminal
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			// non-terminal
			display(count-1, path+".value", v.Elem())
		}
	default:
		// terminal
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

// formatComparable formats Comparable type variables
// https://golang.org/ref/spec#Comparison_operators
func formatComparable(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return "nil"
		} else {
			return fmt.Sprintf("%x", v.Pointer())
		}
	case reflect.Slice, reflect.Array:
		var result string
		length := v.Len()
		result = result + "["
		for i := 0; i < length; i++ {
			result = result + formatAtom(v.Index(i))
			if i != length-1 {
				result = result + ", "
			}
		}
		result = result + "]"
		return result
	case reflect.Struct:
		var result string
		length := v.Type().NumField()
		result = result + "{"
		for i := 0; i < length; i++ {
			result += v.Type().Field(i).Name + ": " + formatAtom(v.Field(i))
			if i != length-1 {
				result = result + ","
			}
		}
		result = result + "}"
		return result
	default:
		return formatAtom(v)
	}

}

// formatAtom formats a value without inspecting its internal structure.
// It is a copy of the the function in gopl.io/ch11/format.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}
