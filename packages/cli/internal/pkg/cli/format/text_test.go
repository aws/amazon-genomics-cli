package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEmptyStruct struct{}

// field name prefix enforcing the ordering
type testSimpleFields struct {
	aIntField    int
	bStringField string
	cBoolField   bool
}

type testStructWithCollections struct {
	aName       string
	bItems1     []testSimpleFields
	cItems2     []testEmptyStruct
	dSomeNumber int
}

type testNestedStruct struct {
	aId   int
	bName string
}

//nolint:structcheck
type testStructWithNestedStruct struct {
	aId        int
	bSubStruct testNestedStruct
}

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
				aIntField:    654,
				bStringField: "Some string",
				cBoolField:   true,
			},
			"TESTSIMPLEFIELDS\t654\tSome string\ttrue\n",
		},
		{
			"Struct with collections",
			testStructWithCollections{
				aName: "This is name",
				bItems1: []testSimpleFields{
					{
						aIntField:    1,
						bStringField: "First",
						cBoolField:   false,
					},
					{
						aIntField:    2,
						bStringField: "Second",
						cBoolField:   true,
					},
				},
				cItems2: []testEmptyStruct{
					{},
					{},
				},
				dSomeNumber: -88,
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
