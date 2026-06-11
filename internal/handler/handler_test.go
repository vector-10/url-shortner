package handler

import (
	"strings"
	"testing"
	"time"
)

func TestExpiryForType(t *testing.T) {
	tests := []struct {
		linkType      string
		wantNil       bool
		expectedHours float64
	}{
		{"general", true, 0},
		{"payment", false, 48},
		{"kyc", false, 168},
		{"onboarding", false, 720},
	}

	for _, tt := range tests {
		createdAt := time.Now()
		result := expiryForType(tt.linkType, createdAt)

		if tt.wantNil && result != nil {
			t.Errorf("linkType %s: expected nil, got %v", tt.linkType, result)
		}

		if !tt.wantNil && result == nil {
			t.Errorf("linkType %s: expected a time, got nil", tt.linkType)
		}

		if result != nil {
			got := result.Sub(createdAt).Hours()
			if got != tt.expectedHours {
				t.Errorf("linkType %s: expected %.0f hours, got %.0f", tt.linkType, tt.expectedHours, got)
			}
		}
	}
}

func TestGenerateSlug(t *testing.T) {
	slug := generateSlug()

	if len(slug) != 10 {
		t.Errorf("expected length 10, got %d", len(slug))
	}

	for _, c := range slug {
		if !strings.Contains(slugCharset, string(c)) {
			t.Errorf("unexpected character %q in slug", c)
		}
	}
}
