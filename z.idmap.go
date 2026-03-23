package main

import (
	"iter"
)

type IdMap_t[T any] struct {
	sort []int
	byId map[int]T
}

func IdMap[T any]() IdMap_t[T] {
	return IdMap_t[T]{
		sort: []int{},
		byId: make(map[int]T),
	}
}

func (in *IdMap_t[T])Add(id int, item T) {  //IdMap_t[T] {
	_, ok := in.byId[id]
	if !ok { in.sort = append(in.sort, id) }
	in.byId[id] = item
//	return in
}

func (m IdMap_t[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for _, id := range m.sort {
			v, ok := m.byId[id]
			if !ok { continue } // future-proof if delete ever exists
			if !yield(id, v) { return }
		}
	}
}
