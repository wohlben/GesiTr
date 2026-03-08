package slug

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Bench Press", "bench-press"},
		{"3/4 Sit-Up", "3-4-sit-up"},
		{"Dumbbell Fly (Incline)", "dumbbell-fly-incline"},
		{"Café Crème", "cafe-creme"},
		{"  Leading & Trailing  ", "leading-trailing"},
		{"already-a-slug", "already-a-slug"},
		{"UPPER CASE", "upper-case"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Generate(tt.input)
			if got != tt.want {
				t.Errorf("Generate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
