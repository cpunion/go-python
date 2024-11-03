package gp

import (
	"testing"
)

func TestMakeList(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		wantLen  int
		wantVals []any
	}{
		{
			name:     "empty list",
			args:     []any{},
			wantLen:  0,
			wantVals: []any{},
		},
		{
			name:     "integers",
			args:     []any{1, 2, 3},
			wantLen:  3,
			wantVals: []any{1, 2, 3},
		},
		{
			name:     "mixed types",
			args:     []any{1, "hello", 3.14},
			wantLen:  3,
			wantVals: []any{1, "hello", 3.14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := MakeList(tt.args...)

			if got := list.Len(); got != tt.wantLen {
				t.Errorf("MakeList() len = %v, want %v", got, tt.wantLen)
			}

			for i, want := range tt.wantVals {
				got := list.GetItem(i).String()
				if got != From(want).String() {
					t.Errorf("MakeList() item[%d] = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestList_SetItem(t *testing.T) {
	list := MakeList(1, 2, 3)
	list.SetItem(1, From("test"))

	// Get the raw value without quotes for comparison
	got := list.GetItem(1).String()

	if got != "test" {
		t.Errorf("List.SetItem() = %v, want %v", got, "test")
	}
}

func TestList_Append(t *testing.T) {
	list := MakeList(1, 2)
	initialLen := list.Len()

	list.Append(From(3))

	if got := list.Len(); got != initialLen+1 {
		t.Errorf("List.Append() length = %v, want %v", got, initialLen+1)
	}

	if got := list.GetItem(2).String(); got != From(3).String() {
		t.Errorf("List.Append() last item = %v, want %v", got, From(3).String())
	}
}

func TestList_Len(t *testing.T) {
	tests := []struct {
		name string
		args []any
		want int
	}{
		{"empty list", []any{}, 0},
		{"single item", []any{1}, 1},
		{"multiple items", []any{1, 2, 3}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := MakeList(tt.args...)
			if got := list.Len(); got != tt.want {
				t.Errorf("List.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
