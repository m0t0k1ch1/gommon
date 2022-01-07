package testutils

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal(t *testing.T, want, got interface{}) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff: %s", diff)
	}
}

func Contains(t *testing.T, s string, substr string) {
	t.Helper()

	if ok := strings.Contains(s, substr); !ok {
		t.Errorf(`"%s" does not contain "%s"`, s, substr)
	}
}
