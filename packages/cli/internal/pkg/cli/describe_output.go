package cli

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func DescribeOutput(output interface{}) string {
	buf := new(strings.Builder)

	_, _ = fmt.Fprintln(buf, "Output of the command has following format:")
	value := reflect.ValueOf(output)
	describeType(buf, value.Type())

	return buf.String()
}

func describeType(buf *strings.Builder, refType reflect.Type) {
	switch refType.Kind() {
	case reflect.Struct:
		describeStructType(buf, refType)
	case reflect.Ptr, reflect.Slice:
		describeType(buf, refType.Elem())
	default:
		describeOtherType(buf, refType)
	}
}

func describeOtherType(buf *strings.Builder, refType reflect.Type) {
	_, _ = fmt.Fprintln(buf, strings.ToUpper(refType.Name()))
}

func describeStructType(buf *strings.Builder, refType reflect.Type) {
	id := strings.ToUpper(refType.Name())
	_, _ = fmt.Fprintf(buf, "%s:", id)
	names := getFieldNamesSortedByName(refType)
	describeSimpleFields(buf, refType, names)
	_, _ = fmt.Fprintln(buf)
	describeCompoundFields(buf, refType, names)

}

func describeSimpleFields(buf *strings.Builder, refType reflect.Type, names []string) {
	for _, name := range names {
		f, ok := refType.FieldByName(name)
		if !ok {
			continue
		}
		if isSimpleType(f.Type) {
			_, _ = fmt.Fprintf(buf, " %s", name)
		}
	}
}

func describeCompoundFields(buf *strings.Builder, refType reflect.Type, names []string) {
	for _, name := range names {
		f, ok := refType.FieldByName(name)
		if !ok {
			continue
		}
		if !isSimpleType(f.Type) {
			describeType(buf, f.Type)
		}
	}
}

func isSimpleType(t reflect.Type) bool {
	return !(t.Kind() == reflect.Struct || t.Kind() == reflect.Slice)
}

func getFieldNamesSortedByName(refType reflect.Type) []string {
	var names []string
	for i := 0; i < refType.NumField(); i++ {
		names = append(names, refType.Field(i).Name)
	}
	sort.Strings(names)
	return names
}
