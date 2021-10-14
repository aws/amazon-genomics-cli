package format

import (
	"reflect"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/stretchr/testify/assert"
)

type unnestedStruct struct {
	SampleInt    int
	SampleString string
}

type simpleSliceStruct struct {
	SampleInt         int
	SampleStringSlice []string
}

type nestedSliceStruct struct {
	SampleInt         int
	SampleNestedSlice []unnestedStruct
}

type invalidNestedStruct struct {
	SampleInt          int
	SampleNestedStruct unnestedStruct
}

type invalidSliceStruct struct {
	SampleInt           int
	InvalidNestedStruct []invalidNestedStruct
}

var (
	unnested1 = unnestedStruct{
		SampleInt:    1,
		SampleString: "Test1",
	}
	unnested2 = unnestedStruct{
		SampleInt:    2,
		SampleString: "Test2",
	}
	simpleSlice = simpleSliceStruct{
		SampleInt: 1,
		SampleStringSlice: []string{
			"hello",
			"world",
		},
	}
	nestedSlice = nestedSliceStruct{
		SampleInt: 1,
		SampleNestedSlice: []unnestedStruct{
			unnested1,
			unnested2,
		},
	}
	invalidNested = invalidNestedStruct{
		SampleInt:          1,
		SampleNestedStruct: unnested1,
	}
	invalidSlice = invalidSliceStruct{
		SampleInt: 1,
		InvalidNestedStruct: []invalidNestedStruct{
			invalidNested,
		},
	}
)

func TestTextTabular_Write(t *testing.T) {
	tests := []struct {
		name     string
		object   interface{}
		expected string
	}{
		{
			"test unnested struct",
			unnested1,
			"SampleInt\tSampleString\n1\t\tTest1",
		},
		{
			"test nested simple slice struct",
			simpleSlice,
			"SampleInt\tSampleStringSlice-string\n1\t\thello\n-\t\tworld\n",
		},
		{
			"test nested complex slice struct",
			nestedSlice,
			"SampleInt\tSampleNestedSlice-SampleInt\tSampleNestedSlice-SampleString\n1\t\t1\t\t\t\tTest1\n-\t\t2\t\t\t\tTest2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			textFormat := &Table{*tabwriter.NewWriter(buf, 0, 8, 0, '\t', 0)}
			textFormat.Write(tt.object)
			actual := buf.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTextTabular_Validate(t *testing.T) {
	tests := []struct {
		name     string
		object   reflect.Type
		expected bool
	}{
		{
			"test invalid nested struct",
			reflect.ValueOf(invalidNested).Type(),
			false,
		},
		{
			"test invalid nested slice struct",
			reflect.ValueOf(invalidSlice).Type(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			textFormat := &Table{*tabwriter.NewWriter(buf, 0, 8, 0, '\t', 0)}
			actual := textFormat.validate(tt.object)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
