package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Demo struct
type Demo struct {
	A int `json:"a,omitempty"`
	B int `json:"b,omitempty"`
}

var (
	o = &Demo{
		A: 1,
	}
	n = &Demo{
		A: 2,
		B: 1,
	}
)

func TestMergeStructByMarshal(t *testing.T) {
	err := MergeStructByMarshal(o, n)
	require.Equal(t, nil, err)
	t.Logf("%+v", o)
}

func BenchmarkMergeStructByMarshal(b *testing.B) {
	MergeStructByMarshal(o, n)
}
