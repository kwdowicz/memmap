package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Function to calculate size, starting address, and end address of variables
func printVariableDetails(variables map[string]interface{}) {
	fmt.Println("Variable details (Start Address, Size, End Address):")
	for name, value := range variables {
		// Get the type and pointer of the value using reflection
		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Ptr {
			fmt.Printf("Error: %s is not a pointer\n", name)
			continue
		}

		ptr := unsafe.Pointer(v.Pointer()) // Get the raw pointer
		startAddress := uintptr(ptr)       // Convert to uintptr

		// Handle special types (like strings) separately
		if v.Elem().Kind() == reflect.String {
			str := v.Elem().Interface().(string)
			headerSize := unsafe.Sizeof(str) // String header size
			contentSize := uintptr(len(str)) // Actual string length
			totalSize := headerSize + contentSize
			endAddress := startAddress + totalSize

			fmt.Printf("%s: Start = %d, Size = %d bytes (Header: %d, Content: %d), End = %d\n",
				name, startAddress, totalSize, headerSize, contentSize, endAddress)
			continue
		}

		// For structs, include size of all embedded strings' content
		if v.Elem().Kind() == reflect.Struct {
			size := v.Elem().Type().Size() // Base size of the struct
			for i := 0; i < v.Elem().NumField(); i++ {
				field := v.Elem().Field(i)
				fieldType := v.Elem().Type().Field(i)

				// Access unexported fields using unsafe
				if !fieldType.IsExported() {
					fieldPtr := unsafe.Pointer(field.UnsafeAddr())
					if field.Kind() == reflect.String {
						str := *(*string)(fieldPtr) // Access the unexported string
						size += uintptr(len(str))  // Add string content size
					}
					continue
				}

				// Access exported fields normally
				if field.Kind() == reflect.String {
					str := field.Interface().(string)
					size += uintptr(len(str)) // Add string content size
				}
			}
			endAddress := startAddress + size
			fmt.Printf("%s: Start = %d, Size = %d bytes (Includes string content), End = %d\n",
				name, startAddress, size, endAddress)
			continue
		}

		// For other types, calculate size directly
		size := v.Elem().Type().Size() // Get the size of the dereferenced type
		endAddress := startAddress + size

		fmt.Printf("%s: Start = %d, Size = %d bytes, End = %d\n",
			name, startAddress, size, endAddress)
	}
}

func main() {
	// Define local variables
	intVar := 42
	floatVar := 3.14
	boolVar := true
	arrayVar := [3]int{1, 2, 3}
	type str struct {
		VarIntName    int    // Exported field
		VarStringName string // Exported field
		unexported    string // Unexported field
	}
	structVar := str{
		VarIntName:    43,
		VarStringName: "exported string contenta",
		unexported:    "unexported string content",
	}
	stringVar := "a"

	// Create a map with pointers to the variables
	variables := map[string]interface{}{
		"intVar":     &intVar,
		"floatVar":   &floatVar,
		"boolVar":    &boolVar,
		"arrayVar":   &arrayVar,
		"structName": &structVar,
		"stringVar":  &stringVar,
	}

	// Print details about the variables
	printVariableDetails(variables)
}
