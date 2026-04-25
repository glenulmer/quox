package htmlHelper

import (
	`strings`
	. `klpm/lib/output`
)

type isWrappable interface { wrappableType() }
type WrappableList []isWrappable
func (WrappableList)wrappableType() {}

type tContent string
func (tContent)wrappableType() {}
func content(items ...interface{}) tContent {
	var b Builder
	for _, s := range items { b.Add(s) }
	return tContent(b.String())
}

type Elem_t struct {
	tag string
	canWrap bool
	wrapped []isWrappable
	kvpairs, classes, styles tElemAttr
}

func (Elem_t)wrappableType() {}

func Elem(tag string) Elem_t {
	wrap := !noWrapTags[tag]

	e := Elem_t{
		tag: Lower(Trim(tag)),
		canWrap: wrap,
		kvpairs: ElemDeco(),
		classes: ElemDeco(),
		styles: ElemDeco(),
	}
	return e
}

func (e Elem_t)String() string { return e.StringLeft() + e.StringRight() }

type ElemList_t []Elem_t 
func (elems ElemList_t)String() string {
	var b Builder
	for _, e := range elems { b.Add(e) }
	return b.String()
}


func (e Elem_t)StringLeft() string {
	if e.tag == `` { return `` }

	var b Builder
	b.Add(`<`, e.tag)

	if e.classes.inOrder.Len() > 0 { b.Add(` class="`, e.classes.inOrder.Join(` `), `"`) }
	if e.kvpairs.inOrder.Len() > 0 {
		var list []string
		for _, key := range e.kvpairs.inOrder.Range() {
			val := e.kvpairs.data[key]
			if val != `` { val = Str(`=` + val) }
			list = append(list, key + val)
		}
		b.Add(` `, Join(list,` `))
	}
	if e.styles.inOrder.Len() > 0 {
		var list []string
		for _, key := range e.styles.inOrder.Range() {
			list = append(list, Str(key, `:`, e.styles.data[key]))
		}
		b.Add(` style="`, Join(list,`; `), `"`)
	}
	if !e.canWrap { b.Add(` /`) }
	b.Add(`>`)
	return b.String()
}

func (e Elem_t)StringRight() string {
	if !e.canWrap { return `` }
	var b Builder
	for _, w := range e.wrapped {
		switch v := w.(type) {
		case Elem_t: b.Add(Elem_t(v))
		case tContent: b.Add(tContent(v))
		}
	}
	b.Add(`</`, e.tag, `>`)
	return b.String() 
}

////////////////////////////////////////////////////////////////
////
////    ElemDeco

type tElemAttr struct {
	inOrder SliceString_t
	data map[string]string
}

func ElemDeco() tElemAttr {
	return tElemAttr {
		inOrder: SliceString(),
		data: make(map[string]string),
	}
}

func parts(s, sep string) (string, string, bool) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 { return parts[0], ``, false }
	return parts[0], parts[1], true
}

func guessAttrib(s string) string {
	if Contains(s, `=`) { return `kv` }
	if Contains(s, `:`) { return `style` }
	return `class`
}

func (in Elem_t)Set(list ...string) Elem_t {
	for _, item := range list {
		switch guessAttrib(item) {
		case `kv`: in = in.KVs(item)
		case `style`: in = in.Style(item)
		default: in = in.Class(item)
		}
	}
	return in
}

func (in Elem_t)Class(list ...string) Elem_t {
	for _, n := range list {
		n = Trim(n)
		if n == `` { continue }
		_, exists := in.classes.data[n]
		if exists { continue }
		in.classes.inOrder = in.classes.inOrder.Add(n)
		in.classes.data[n] = ``
	}
	return in
}

func (in Elem_t)Style(list ...string) Elem_t {
	for _, nstyle := range list {
		left, right, has2 := parts(nstyle, `:`)
		left = Trim(left)
		if !has2 { delete(in.styles.data, left); continue }
		_, exists := in.styles.data[left]
		if !exists { in.styles.inOrder = in.styles.inOrder.Add(left) }
		in.styles.data[left] = Trim(right)
	}
	return in
}

func (in Elem_t)KVs(list ...string) Elem_t {
	for _, nstyle := range list {
		left, right, _ := parts(nstyle, `=`)
		left = Trim(left)
		_, exists := in.kvpairs.data[left]
		if !exists { in.kvpairs.inOrder = in.kvpairs.inOrder.Add(left) }
		in.kvpairs.data[left] = Trim(right)
	}
	return in
}

func (in Elem_t)KV(key string, value ...interface{}) Elem_t { 
	if len(value) > 0 { key = key + `=` + Q(Str(value[0])) }
	return in.KVs([]string{key}...)
}
func (e Elem_t)Id(a ...any) Elem_t { return e.KV(`id`,Str(a...)) }
func (e Elem_t)Name(a ...any) Elem_t { return e.KV(`name`, Str(a...)) }
func (e Elem_t)Value(a ...any) Elem_t { return e.KV(`value`, Str(a...)) }
func (e Elem_t)Type(a ...any) Elem_t { return e.KV(`type`, Str(a...)) }
func (e Elem_t)Place(a ...any) Elem_t { return e.KV(`placeholder`, Str(a...)) }
func (e Elem_t)NCols(a any) Elem_t { return e.KV(`colspan`, Str(a)) }
func (e Elem_t)Width(a any) Elem_t { return e.KV(`width`, Str(a)) }
func (e Elem_t)Len(i int) Elem_t { return e.KV(`maxlength`, i) }

func (e Elem_t)Orig(a ...any) Elem_t { return e.KV(`data-orig`,Str(a...)) }
func (e Elem_t)VO(a ...any) Elem_t { return e.Value(a...).Orig(a...) }
func (e Elem_t)CO(b bool) Elem_t { return e.Check(b).Orig(Bit(b)) }
func (e Elem_t)SelO(a any) Elem_t { return e.Choose(a).Orig(a) }
func (e Elem_t)Args(a ...any) Elem_t { return e.KV(`data-record`, Str(a...)) }
func (e Elem_t)Post(s string) Elem_t { return e.KV(`data-post`, s) }
func (e Elem_t)Opt(b bool) Elem_t { return e.KV(`data-optional`, b) }
func (e Elem_t)Optional(b bool) Elem_t { return e.Opt(b) }
func (e Elem_t)Data(key string, val ...interface{}) Elem_t { return e.KV(Str(`data-`, key), val...) }

func (e Elem_t)VOPl(a ...any) Elem_t { return e.Value(a...).Orig(a...).Place(a...) }

func (in tElemAttr)Copy() tElemAttr {
	out := tElemAttr{ inOrder: in.inOrder.Copy() }
	out.data = make(map[string]string, len(in.data))
	for k, v := range in.data {
		out.data[k] = v
	}
	return out
}

func (in tElemAttr)Cut(name string) tElemAttr {
	name = Trim(name)
	_, exists := in.data[name]
	if !exists { return in }
	delete(in.data, name)
	in.inOrder = in.inOrder.Cut(name)
	return in
}

func (e Elem_t)CutAttrib(name string) Elem_t { e.kvpairs = e.kvpairs.Cut(name); return e }
func (e Elem_t)CutClass(name string) Elem_t { e.classes = e.classes.Cut(name); return e }
func (e Elem_t)CutStyle(name string) Elem_t { e.styles = e.styles.Cut(name); return e }

func (e Elem_t)Text(items ...interface{}) Elem_t {
	return e.Wrap(content(items...))
}

func Wrappable(list []Elem_t) WrappableList {
    wrappable := make(WrappableList, len(list))
    for i, elem := range list {
        wrappable[i] = elem
    }
    return wrappable
}
