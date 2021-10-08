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

func TestText_writeMap(t *testing.T) {

	type args struct {
		value  interface{}
		indent int
	}

	emptyMap := make(map[string]interface{})

	simpleMap := map[string]interface{}{"a": "abcd", "b": "bcde"}

	mapWithArray := make(map[string]interface{})
	mapWithArray["a"] = []string{"a", "b", "c"}

	nestedMap := make(map[string]interface{})
	nestedMap["a"] = "abcd"
	nestedMap["b"] = map[string]interface{}{"foo": "foz", "baa": "baz"}
	nestedMap["c"] = "bcde"

	deeplyNestedMap := make(map[string]interface{})
	deeplyNestedMap["a"] = "abcd"
	deeplyNestedMap["b"] = nestedMap
	deeplyNestedMap["c"] = mapWithArray
	deeplyNestedMap["d"] = simpleMap

	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "simple map",
			args: args{
				value:  simpleMap,
				indent: 0,
			},
			expected: `a:	abcd
b:	bcde
`,
		},
		{
			name: "empty map",
			args: args{
				value:  emptyMap,
				indent: 0,
			},
			expected: "",
		},
		{
			name: "map with array",
			args: args{
				value:  mapWithArray,
				indent: 0,
			},
			expected: `a:	[a b c]
`,
		},
		{
			name: "nested map",
			args: args{
				value:  nestedMap,
				indent: 0,
			},
			expected: `a:	abcd
b:	
	baa:	baz
	foo:	foz
c:	bcde
`,
		},
		{
			name: "deeply nested map",
			args: args{
				value:  deeplyNestedMap,
				indent: 0,
			},
			expected: `a:	abcd
b:	
	a:	abcd
	b:	
		baa:	baz
		foo:	foz
	c:	bcde
c:	
	a:	[a b c]
d:	
	a:	abcd
	b:	bcde
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			textFormat := &Text{buf}
			textFormat.Write(tt.args.value)
			actual := buf.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
