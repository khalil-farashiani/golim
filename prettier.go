package main

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

func makeTable(data []interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if len(data) > 0 {
		var tr = table.Row{}
		switch reflect.TypeOf(data[0]).Kind() {
		case reflect.Struct:
			val := reflect.ValueOf(data[0])
			for i := 0; i < val.NumField(); i++ {
				tr = append(tr, val.Type().Field(i).Name)
			}
			t.AppendHeader(tr)
			t.AppendSeparator()
		}
	}
	for _, d := range data {
		var tr = table.Row{}
		val := reflect.ValueOf(d)
		for i := 0; i < val.NumField(); i++ {
			tr = append(tr, val.Field(i).Interface())
		}
		t.AppendHeader(tr)
		t.AppendSeparator()
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func toSlice(i interface{}) []interface{} {
	var out []interface{}

	rv := reflect.ValueOf(i)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
	}
	return out
}
