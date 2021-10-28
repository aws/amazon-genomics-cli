package slices

import (
	"reflect"
	"testing"
)

func TestDeDuplicateStrings(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want []string
	}{
		{
			name: "Should gracefully handle empty slice",
			strs: []string{},
			want: []string{},
		},
		{
			name: "Should success when no duplicates",
			strs: []string{"C", "B", "A"},
			want: []string{"A", "B", "C"},
		},
		{
			name: "Should remove duplicates",
			strs: []string{"C", "B", "C", "B", "A"},
			want: []string{"A", "B", "C"},
		},
		{
			name: "Should gracefully handle nil",
			strs: nil,
			want: nil,
		},
		{
			name: "Should be case sensitive",
			strs: []string{"C", "B", "A", "a"},
			want: []string{"A", "B", "C", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeDuplicateStrings(tt.strs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeDuplicateStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
