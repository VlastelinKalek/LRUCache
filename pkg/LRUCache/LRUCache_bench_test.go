package lrucache

import (
	"fmt"
	"testing"
)

func BenchmarkLRUCache(b *testing.B) {
	var c = NewLRUCache(3)
	for i := 0; i < b.N; i++ {
		c.Add("1", "A")
		c.Add("2", "B")
		c.Add("3", "C")
		c.Get("2")
		c.Add("4", "D")
		c.Add("3", "E")
		c.Remove("4")
		c.RemoveFront()
		c.RemoveBack()
	}
}

func BenchmarkAdd(b *testing.B) {
	var c = NewLRUCache(b.N)
	for i := 0; i < b.N; i++ {
		c.Add(fmt.Sprintf("%d", i), fmt.Sprintf("%d", i))
	}
}
