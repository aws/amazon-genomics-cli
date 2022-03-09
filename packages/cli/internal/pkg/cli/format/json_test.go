package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson_WriteStruct(t *testing.T) {

	tests := []struct {
		name     string
		object   interface{}
		expected string
	}{
		{
			"Struct with no fields",
			testEmptyStruct{},
			"{}\n",
		},
		{
			"int",
			123,
			"123\n",
		},
		{
			"string",
			"foo bar",
			"\"foo bar\"\n",
		},
		{
			"Simple fields",
			testSimpleFields{
				AIntField:    654,
				BStringField: "Some string",
				CBoolField:   true,
			},
			`{
	"AIntField": 654,
	"BStringField": "Some string",
	"CBoolField": true
}
`,
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
			`{
	"AName": "This is name",
	"BItems1": [
		{
			"AIntField": 1,
			"BStringField": "First",
			"CBoolField": false
		},
		{
			"AIntField": 2,
			"BStringField": "Second",
			"CBoolField": true
		}
	],
	"CItems2": [
		{},
		{}
	],
	"DSomeNumber": -88
}
`,
		},
		{
			"Collection of strings",
			[]string{"One", "Two", "Three"},
			`[
	"One",
	"Two",
	"Three"
]
`,
		},
		{
			"Collection of structs",
			[]testNestedStruct{
				{1, "First"},
				{2, "Second"},
			},
			`[
	{
		"AId": 1,
		"BName": "First"
	},
	{
		"AId": 2,
		"BName": "Second"
	}
]
`,
		},
		{
			"Nested Struct",
			[]testStructWithNestedStruct{
				{1, testNestedStruct{100, "First Nested"}},
				{2, testNestedStruct{200, "Second Nested"}},
			},
			`[
	{
		"AId": 1,
		"BSubStruct": {
			"AId": 100,
			"BName": "First Nested"
		}
	},
	{
		"AId": 2,
		"BSubStruct": {
			"AId": 200,
			"BName": "Second Nested"
		}
	}
]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			jsonFormat := &Json{buf}
			jsonFormat.Write(tt.object)
			actual := buf.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
