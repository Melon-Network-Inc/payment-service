package utils

import "testing"

func TestString(t *testing.T) {
	type args struct {
		id uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test_20", args{id: uint(20)}, "20"},
		{"test_0", args{id: 0}, "0"},
		{"test_1000", args{id: 1000}, "1000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.id); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
