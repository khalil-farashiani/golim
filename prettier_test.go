package main

import (
	"reflect"
	"testing"
)

type MockData struct {
	Name string
	Age  int
}

func TestMakeTable(t *testing.T) {
	data := []interface{}{
		MockData{Name: "John", Age: 30},
		MockData{Name: "Doe", Age: 40},
	}

	makeTable(data)
}

func TestToSlice(t *testing.T) {
	a := []MockData{
		{Name: "John", Age: 30},
		{Name: "Doe", Age: 40},
	}

	result := toSlice(a)
	if reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Errorf("Expected Slice got %v", reflect.TypeOf(result).Kind())
	}
}
