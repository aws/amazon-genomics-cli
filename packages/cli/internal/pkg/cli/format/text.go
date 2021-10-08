package format

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

const delimiter = "\t"

type Text struct {
	writer io.Writer
}

type indentedPair struct {
	key    string
	value  interface{}
	indent int
}

type valueSlice []reflect.Value

func (k valueSlice) Len() int {
	return len(k)
}
func (k valueSlice) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}
func (k valueSlice) Less(i, j int) bool {
	return k[i].String() < k[j].String()
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
	case reflect.Map:
		f.writeMap(value, 0)
	default:
		f.writeOther(value)
	}
}

func (f *Text) writeMap(value reflect.Value, indent int) {
	keys := value.MapKeys()
	sort.Sort(valueSlice(keys))
	for _, key := range keys {
		if key.Kind() != reflect.String {
			log.Warn().Msg("currently only able to format maps with string keys, using default formatting")
			f.writeOther(value)
		}
		i := value.MapIndex(key).Interface()
		f.writeIndentedPair(
			indentedPair{
				key:    key.String(),
				value:  i,
				indent: indent,
			})
	}
}

func (f *Text) writeIndentedPair(pair indentedPair) {
	indentStr := strings.Repeat(delimiter, pair.indent)
	f.write(fmt.Sprintf("%s%s:%s", indentStr, pair.key, delimiter))
	if reflect.TypeOf(pair.value).Kind() == reflect.Map {
		f.newLine()
		f.writeMap(reflect.ValueOf(pair.value), pair.indent+1)
	} else {
		f.write(fmt.Sprintf("%s", pair.value))
		f.newLine()
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
