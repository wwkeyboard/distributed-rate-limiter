package limiter

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// This will fail if run on the minute boundery. This should be fixed but may be outside the scope of this project.
func Test_key(t *testing.T) {
	minute := time.Now().Minute()
	slug := "/something"
	k := key(slug)

	if !strings.Contains(k, slug) {
		t.Errorf("key didn't include slug")
	}
	if !strings.Contains(k, fmt.Sprint(minute)) {
		t.Errorf("key didn't include minute")
	}
}
