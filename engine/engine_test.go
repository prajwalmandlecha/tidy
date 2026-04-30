package engine

import (
	"testing"

	"github.com/prajwalmandlecha/tidy/config"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		rule     config.Rule
		want     bool
	}{
		{
			name:     "matches by extension",
			filePath: "/downloads/budget.pdf",
			rule:     config.Rule{Extensions: []string{".pdf"}},
			want:     true,
		},
		{
			name:     "no match — wrong extension",
			filePath: "/downloads/photo.png",
			rule:     config.Rule{Extensions: []string{".pdf"}},
			want:     false,
		},
		{
			name:     "matches by pattern",
			filePath: "/downloads/Screenshot 2025.png",
			rule:     config.Rule{Pattern: "Screenshot*"},
			want:     true,
		},
		{
			name:     "no match — wrong pattern",
			filePath: "/downloads/photo.png",
			rule:     config.Rule{Pattern: "Screenshot*"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Matches(tt.filePath, tt.rule)
			if got != tt.want {
				t.Errorf("Matches(%q, %v) = %v; want %v", tt.filePath, tt.rule, got, tt.want)
			}
		})
	}
}
