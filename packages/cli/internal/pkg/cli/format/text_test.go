package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_WriteStruct(t *testing.T) {

	tests := []struct {
		name     string
		object   interface{}
		expected string
	}{
		{
			"Struct with no fields",
			testEmptyStruct{},
			"TESTEMPTYSTRUCT\n",
		},
		{
			"int",
			123,
			"123\n",
		},
		{
			"string",
			"foo bar",
			"foo bar\n",
		},
		{
			"Simple fields",
			testSimpleFields{
				AIntField:    654,
				BStringField: "Some string",
				CBoolField:   true,
			},
			"TESTSIMPLEFIELDS\t654\tSome string\ttrue\n",
		},
		{
			"Struct with collections",
			testStructWithCollections{
				AName: "This is name",
				BItems1: []testSimpleFields{
					{
						AIntField:    1,
						BStringField: "First",
						CBoolField:   false,
					},
					{
						AIntField:    2,
						BStringField: "Second",
						CBoolField:   true,
					},
				},
				CItems2: []testEmptyStruct{
					{},
					{},
				},
				DSomeNumber: -88,
			},
			`TESTSTRUCTWITHCOLLECTIONS	This is name	-88
TESTSIMPLEFIELDS	1	First	false
TESTSIMPLEFIELDS	2	Second	true
TESTEMPTYSTRUCT
TESTEMPTYSTRUCT
`,
		},
		{
			"Collection of strings",
			[]string{"One", "Two", "Three"},
			"One\nTwo\nThree\n",
		},
		{
			"Collection of structs",
			[]testNestedStruct{
				{1, "First"},
				{2, "Second"},
			},
			`TESTNESTEDSTRUCT	1	First
TESTNESTEDSTRUCT	2	Second
`,
		},
		{
			"Nested Struct",
			[]testStructWithNestedStruct{
				{1, testNestedStruct{100, "First Nested"}},
				{2, testNestedStruct{200, "Second Nested"}},
			},
			`TESTSTRUCTWITHNESTEDSTRUCT	1
TESTNESTEDSTRUCT	100	First Nested
TESTSTRUCTWITHNESTEDSTRUCT	2
TESTNESTEDSTRUCT	200	Second Nested
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			textFormat := &Text{buf}
			textFormat.Write(tt.object)
			actual := buf.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
