package format

type testEmptyStruct struct{}

// field name prefix enforcing the ordering
type testSimpleFields struct {
	AIntField    int
	BStringField string
	CBoolField   bool
}

type testStructWithCollections struct {
	AName       string
	BItems1     []testSimpleFields
	CItems2     []testEmptyStruct
	DSomeNumber int
}

type testNestedStruct struct {
	AId   int
	BName string
}

//nolint:structcheck
type testStructWithNestedStruct struct {
	AId        int
	BSubStruct testNestedStruct
}
