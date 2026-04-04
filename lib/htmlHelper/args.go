package htmlHelper

import (
	. `quo2/lib/output`
)

type args_t struct {
	data map[string]any
}

func Args() args_t {
	return args_t{data: make(map[string]any)}
}

func (a args_t) Add(key string, value any) args_t {
	a.data[key] = value
	return a
}

func (a args_t) String() string {
	var parts []string
	for k, v := range a.data { parts = append(parts, Str(k,`:`,v)) }
	return Join(parts, `,`)
}

func (e Elem_t) Record(args args_t) Elem_t { 
	return e.Data(`record`, args.String())
}
