package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type Foo struct {
	bar string
	baz int
}

func main() {
	blar := make(map[Foo]int)
	blar2 := make(map[[3]int]int)

	intArray := [3]int{10, 12, 14}

	one := &Foo{bar: "one", baz: 1}
	two := &Foo{bar: "two", baz: 2}
	three := &Foo{bar: "three", baz: 3}

	blar[*one] = 1
	blar[*two] = 2
	blar[*three] = 3

	blar2[intArray] = 16

	Display("foo", blar)
	Display("bar", blar2)
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

// Extend `Display` so that it can display maps whose keys are structs or arrays
func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		length := v.NumField()
		for i := 0; i < length; i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatComparable(key) /*the rendered key*/), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default:
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
