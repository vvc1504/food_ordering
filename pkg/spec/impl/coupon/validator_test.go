package coupon

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestValidator(t *testing.T) {
	tempDir := t.TempDir()

	// Create 3 test files
	// File 1: COUPON1, COUPON2, SHORT
	// File 2: COUPON1, COUPON3, INVALID_LEN
	// File 3: COUPON2, COUPON3, ALONE

	files := []struct {
		name    string
		content string
	}{
		{"f1.gz", "COUPON01 COUPON02 SHORT"},
		{"f2.gz", "COUPON01 COUPON03 LONGSTRINGHERE"},
		{"f3.gz", "COUPON02 COUPON03 ALONE"},
	}

	var paths []string
	for _, f := range files {
		path := filepath.Join(tempDir, f.name)
		paths = append(paths, path)

		file, _ := os.Create(path)
		gz := gzip.NewWriter(file)
		gz.Write([]byte(f.content))
		gz.Close()
		file.Close()
	}

	v, err := NewValidator(paths)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		code string
		want bool
	}{
		{"COUPON01", true},      // In file 1 and 2
		{"COUPON02", true},      // In file 1 and 3
		{"COUPON03", true},      // In file 2 and 3
		{"SHORT", false},        // Too short (5 chars)
		{"ALONE", false},        // Only in file 3 (short too)
		{"UNKNOWN", false},      // Not in any file
		{"ABC", false},          // Too short
		{"VERYLONGCODE", false}, // Too long
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, _ := v.Validate(tt.code)
			if got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
