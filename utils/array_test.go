package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayIndexOf(t *testing.T) {
	assert.Equal(t, 1, ArrayIndexOf([]string{"1", "2", "3"}, "2"))
	assert.Equal(t, -1, ArrayIndexOf([]string{"1", "2", "3"}, "0"))
}
