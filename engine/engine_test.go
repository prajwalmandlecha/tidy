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

func TestIsIgnored(t *testing.T) {
	ignores := []string{"*.crdownload", "*.part", ".DS_Store"}

	tests := []struct {
		path string
		want bool
	}{
		{"/downloads/ubuntu.iso.crdownload", true},
		{"/downloads/movie.mkv.part", true},
		{"/downloads/.DS_Store", true},
		{"/downloads/report.pdf", false},
	}

	for _, tt := range tests {
		got := IsIgnored(tt.path, ignores)
		if got != tt.want {
			t.Errorf("IsIgnored(%q) = %v; want %v", tt.path, got, tt.want)
		}
	}
}
