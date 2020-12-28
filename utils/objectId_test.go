package utils

import (
	"sync"
	"testing"
)

func TestObjectID_Hex(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Log(NewObjectID().Hex())
		}()
	}
	wg.Wait()
}

func BenchmarkObjectID_Hex(b *testing.B) {
	b.Log(NewObjectID().Hex())
}
