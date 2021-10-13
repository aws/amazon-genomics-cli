package format

import (
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TextTabular struct {
	writer tabwriter.Writer
}

const (
	noPrefix          = ""
	emptyVal          = "-"
	prefixNameDivider = "-"
	tabularDelimiter  = "\t"
)

func (f *TextTabular) Write(o interface{}) {
	val := reflect.ValueOf(o)
	if !f.validate(val.Type()) {
		textWriter.Write(o)
	} else {
		f.writeHeader(val)
		f.writeValue(val)
		f.writer.Flush()
	}
}

func (f *TextTabular) validate(val reflect.Type) bool {
	if f.isSlice(val) {
		return f.validateStruct(val.Elem())
	} else if f.isStruct(val) {
		return f.validateStruct(val)
	} else {
		return true
	}
}

func (f *TextTabular) validateStruct(val reflect.Type) bool {
	for i := 0; i < val.NumField(); i++ {
		structFieldType := val.Field(i).Type
		if f.isStruct(structFieldType) {
			return false
		} else if f.isSlice(structFieldType) && !f.validateSlice(structFieldType.Elem()) {
			return false
		}
	}
	return true
}

func (f *TextTabular) validateSlice(val reflect.Type) bool {
	if f.isSimpleType(val) {
		return true
	} else if f.isStruct(val) {
		for i := 0; i < val.NumField(); i++ {
			structFieldType := val.Field(i).Type
			if !f.isSimpleType(structFieldType) {
				return false
			}
		}
	}
	return true
}

func (f *TextTabular) writeHeader(value reflect.Value) {
	f.write(strings.Join(f.getHeader(noPrefix, value.Type()), tabularDelimiter))
	f.newLine()
}

func (f *TextTabular) getHeader(prefix string, value reflect.Type) []string {
	if f.isStruct(value) {
		headers := f.getStructHeaders(prefix, value)
		return headers
	} else if f.isSlice(value) {
		headers := f.getHeader(prefix, value.Elem())
		return headers
	} else {
		return []string{prefix + value.Name()}
	}
}

func (f *TextTabular) getStructHeaders(prefix string, val reflect.Type) []string {
	headers := make([]string, 0)
	for i := 0; i < val.NumField(); i++ {
		structFieldType := val.Field(i).Type
		if f.isSlice(structFieldType) {
			headers = append(headers, f.getHeader(prefix+val.Field(i).Name+prefixNameDivider, val.Field(i).Type.Elem())...)
		} else {
			headers = append(headers, prefix+val.Field(i).Name)
		}
	}
	return headers
}

func (f *TextTabular) writeValue(value reflect.Value) {
	val := reflect.Indirect(value)
	if f.isSimpleType(value.Type()) {
		f.write(fmt.Sprint(value))
	} else if f.isStruct(value.Type()) {
		maxCollectionSize := f.getMaxSliceSize(value)
		if maxCollectionSize == 0 {
			f.write(f.getSimpleStructOutput(val))
		} else {
			f.write(f.getSliceOutput(maxCollectionSize, value))
		}
	} else {
		f.writeCollection(value)
	}
}

func (f *TextTabular) writeCollection(value reflect.Value) {
	for i := 0; i < value.Len(); i++ {
		maxCollectionSize := f.getMaxSliceSize(value.Index(i))
		if maxCollectionSize == 0 {
			f.write(f.getSimpleStructOutput(value.Index(i)))
		} else {
			f.write(f.getSliceOutput(maxCollectionSize, value.Index(i)))
		}
		f.newLine()
	}
}

func (f *TextTabular) getMaxSliceSize(value reflect.Value) int {
	maxSize := 0
	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Field(i)
		if f.isSlice(field.Type()) && field.Len() > maxSize {
			maxSize = field.Len()
		}
	}
	return maxSize
}

func (f *TextTabular) getSliceOutput(maxCollectionSize int, value reflect.Value) string {
	val := reflect.Indirect(value)
	outputRowList := make([]string, 0)
	output := ""

	for rowIndex := 0; rowIndex < maxCollectionSize; rowIndex++ {
		for structIndex := 0; structIndex < val.Type().NumField(); structIndex++ {
			structField := val.Field(structIndex)
			if f.isSimpleType(structField.Type()) {
				outputRowList = f.appendSimpleValue(structField, outputRowList, rowIndex)
			} else if f.isStruct(structField.Type().Elem()) {
				outputRowList = f.appendSliceStructRow(structField, outputRowList, rowIndex)
			} else if f.isSlice(structField.Type()) {
				outputRowList = f.appendSimpleSlice(structField, outputRowList, rowIndex)
			} else {
				outputRowList = append(outputRowList, fmt.Sprint(structField))
			}
		}
		output += strings.Join(outputRowList, tabularDelimiter) + "\n"
		outputRowList = outputRowList[:0]
	}

	return output
}

func (f *TextTabular) appendSliceStructRow(field reflect.Value, outputList []string, row int) []string {
	for k := 0; k < field.Type().Elem().NumField(); k++ {
		if row >= field.Len() {
			outputList = append(outputList, emptyVal)
		} else {
			outputList = append(outputList, fmt.Sprint(field.Index(row).Field(k)))
		}
	}
	return outputList
}

func (f *TextTabular) appendSimpleSlice(field reflect.Value, outputList []string, row int) []string {
	if row >= field.Len() {
		outputList = append(outputList, emptyVal)
	} else {
		outputList = append(outputList, fmt.Sprint(field.Index(row)))
	}
	return outputList
}

func (f *TextTabular) appendSimpleValue(field reflect.Value, outputList[]string, row int) []string {
	if row > 0 {
		outputList = append(outputList, emptyVal)
	} else {
		outputList = append(outputList, fmt.Sprint(field))
	}
	return outputList
}

func (f *TextTabular) getSimpleStructOutput(val reflect.Value) string {
	outputList := make([]string, 0)
	for i := 0; i < val.Type().NumField(); i++ {
		outputList = append(outputList, fmt.Sprint(val.Field(i)))
	}
	return strings.Join(outputList, tabularDelimiter)
}

func (f *TextTabular) isSimpleType(value reflect.Type) bool {
	return !(value.Kind() == reflect.Slice || value.Kind() == reflect.Struct)
}

func (f *TextTabular) isSlice(value reflect.Type) bool {
	return value.Kind() == reflect.Slice
}

func (f *TextTabular) isStruct(value reflect.Type) bool {
	return value.Kind() == reflect.Struct
}

func (f *TextTabular) write(v interface{}) {
	_, _ = fmt.Fprint(&f.writer, v)
}

func (f *TextTabular) newLine() {
	_, _ = fmt.Fprintln(&f.writer)
}
