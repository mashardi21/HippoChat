package settings

import (
	"testing"
)

func TestCheckForConfig(t *testing.T) {
	check := checkForConfig()

	if check {
		t.Errorf("Check: %v; Expected: false", check)
	}
}
