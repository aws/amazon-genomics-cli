package format

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

const delimiter = "\t"

type Text struct {
	writer io.Writer
}

func (f *Text) Write(o interface{}) {
	f.writeValue(reflect.ValueOf(o))
}

func (f *Text) writeValue(value reflect.Value) {
	switch value.Kind() {
	case reflect.Struct:
		f.writeStruct(value)
	case reflect.Ptr:
		f.writeValue(value.Elem())
	case reflect.Slice:
		f.writeSliceValue(value)
	default:
		f.writeOther(value)
	}
}

func (f *Text) writeOther(value reflect.Value) {
	f.write(value)
	f.newLine()
}

func (f *Text) writeSliceValue(sliceValue reflect.Value) {
	for i := 0; i < sliceValue.Len(); i++ {
		itemValue := sliceValue.Index(i)
		f.writeValue(itemValue)
	}
}

func (f *Text) writeStruct(value reflect.Value) {
	f.writeSimpleFields(value)
	f.writeCompoundFields(value)
}

func (f *Text) writeSimpleFields(value reflect.Value) {
	id := strings.ToUpper(value.Type().Name())
	f.write(id)
	fieldNames := getFieldNames(value)
	for _, fieldName := range fieldNames {
		field := value.FieldByName(fieldName)
		if isSimpleType(field) {
			f.write(delimiter)
			f.write(fmt.Sprint(field))
		}
	}
	f.newLine()
}

func isSimpleType(value reflect.Value) bool {
	return !(value.Kind() == reflect.Struct || value.Kind() == reflect.Slice)
}

func (f *Text) write(v interface{}) {
	_, _ = fmt.Fprint(f.writer, v)
}

func (f *Text) newLine() {
	_, _ = fmt.Fprintln(f.writer)
}

func getFieldNames(value reflect.Value) []string {
	var names []string
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		names = append(names, valueType.Field(i).Name)
	}
	sort.Strings(names)
	return names
}

func (f Text) writeCompoundFields(value reflect.Value) {
	fieldNames := getFieldNames(value)
	for _, fieldName := range fieldNames {
		field := value.FieldByName(fieldName)
		if !isSimpleType(field) {
			f.writeValue(field)
		}
	}
}
