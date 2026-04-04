package htmlHelper

import "strings"

type SliceString_t struct { x []string }

func SliceString() SliceString_t { return SliceString_t{} }

/*
func (in SliceString_t)Out() string {
	return "(SS|" + fmt.Sprint(len(in.x)) + "|" + fmt.Sprint(in.x) + ")"
}
*/

func fix(in []string) SliceString_t { return SliceString_t{x:in} }

func (in SliceString_t)Add(i string) SliceString_t {
	in.x = append(in.x, i)
	return in
}

func (in SliceString_t)Cut(name string) SliceString_t {
	fresh := SliceString()
	for _, item := range in.x {
		if name == item { continue }
		fresh = fresh.Add(item)
	}
	return fresh
}

func (in SliceString_t)Copy() SliceString_t {
	copy := SliceString()
	for _, item := range in.x { copy = copy.Add(item) }
	return copy
}

func (in SliceString_t)Len() int { return len(in.x) }

func (in SliceString_t)Join(sep string) string { return strings.Join(in.x, sep) }

func (in SliceString_t)Range() []string { return in.Copy().x }
