package main

import (
	"container/heap"
	"fmt"
)

const MAX_CHAR = 26

type CharHeap []struct {
	freq int
	ch   rune
}

func (h CharHeap) Len() int           { return len(h) }
func (h CharHeap) Less(i, j int) bool { return h[i].freq > h[j].freq }
func (h CharHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CharHeap) Push(x interface{}) {
	*h = append(*h, x.(struct {
		freq int
		ch   rune
	}))
}

func (h *CharHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func rearrangeString(str string) string {
	N := len(str)

	count := make([]int, MAX_CHAR)
	for _, ch := range str {
		count[ch-'a']++
	}

	pq := &CharHeap{}
	heap.Init(pq)
	for c := 'a'; c <= 'z'; c++ {
		val := int(c - 'a')
		if count[val] > 0 {
			heap.Push(pq, struct {
				freq int
				ch   rune
			}{count[val], c})
		}
	}

	var result []rune
	prev := struct {
		freq int
		ch   rune
	}{-1, '#'}

	for pq.Len() > 0 {
		k := heap.Pop(pq).(struct {
			freq int
			ch   rune
		})
		result = append(result, k.ch)

		if prev.freq > 0 {
			heap.Push(pq, prev)
		}

		k.freq--
		prev = k
	}

	if N != len(result) {
		return ""
	} else {
		return string(result)
	}
}

func main() {
	fmt.Println(rearrangeString("ssskfsdcdsdsdddswdicddd"))

}
