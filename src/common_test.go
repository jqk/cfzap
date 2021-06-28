package cfzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareStringArray(t *testing.T) {
	assert.True(t, CompareStringArray(nil, nil), "two nil string array should be equal.")
	assert.False(t, CompareStringArray(nil, []string{}), "one nil string array should not be equal.")
	assert.True(t, CompareStringArray([]string{}, []string{}), "two empty string array should be equal.")
	assert.False(t, CompareStringArray([]string{"abc"}, []string{}), "one empty string array should not be equal.")
	assert.False(t, CompareStringArray([]string{"abc"}, []string{"def"}), "two different string array should not be equal.")
	assert.True(t, CompareStringArray([]string{"abc"}, []string{"abc"}), "two string array with same content should be equal.")
}
