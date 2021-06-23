package cfzap

import (
	"testing"
)

func TestCompareStringArray(t *testing.T) {
	if !CompareStringArray(nil, nil) {
		t.Error("two nil string array should be equal.")
	}

	if CompareStringArray(nil, []string{}) {
		t.Error("one nil string array should not be equal.")
	}

	if !CompareStringArray([]string{}, []string{}) {
		t.Error("two empty string array should be equal.")
	}

	if CompareStringArray([]string{"abc"}, []string{}) {
		t.Error("one empty string array should not be equal.")
	}

	if CompareStringArray([]string{"abc"}, []string{"def"}) {
		t.Error("two diffrent string array should not be equal.")
	}

	if !CompareStringArray([]string{"abc"}, []string{"abc"}) {
		t.Error("two string array with same content should be equal.")
	}
}
