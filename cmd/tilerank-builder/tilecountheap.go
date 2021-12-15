package main

import (
	"bufio"
	"container/heap"
	"sort"
)

type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
}

type TileCountHeap []*TileCountReader

type TileCountReader struct {
	tc      TileCount
	scanner bufio.Scanner
	index   int
}

func (tch TileCountHeap) Len() int {
	return len(tch)
}

func (tch TileCountHeap) Less(i, j int) bool {
	return TileCountLess(tch[i].tc, tch[j].tc)
}

func (tch TileCountHeap) Swap(i, j int) {
	tch[i], tch[j] = tch[j], tch[i]
	tch[i].index = i
	tch[j].index = j
}

func (tch *TileCountHeap) Push(x interface{}) {
	n := len(*tch)
	item := x.(*TileCountReader)
	item.index = n
	*tch = append(*tch, item)
}

func (tch *TileCountHeap) Pop() interface{} {
	old := *tch
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*tch = old[0 : n-1]
	return item
}

func (tch *TileCountHeap) fix(item *TileCountReader, tc TileCount, scanner bufio.Scanner) {
	item.scanner = scanner
	item.tc = tc
	heap.Fix(tch, item.index)
}
