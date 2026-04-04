package htmlHelper

/*
func (e Elem_t) Wrap(items ...isWrappable) Elem_t {
    e.wrapped = append(e.wrapped, items...)
    return e
}
*/

func (e Elem_t) Wrap(items ...any) Elem_t {
    for _, item := range items {
        switch v := item.(type) {

        case []string:
            for _, s := range v {
                e.wrapped = append(e.wrapped, content(s))
            }

        case Elem_t:
            e.wrapped = append(e.wrapped, v)
        case []Elem_t:
            for _, elem := range v {
                e.wrapped = append(e.wrapped, elem)
            }

        case tContent:
            e.wrapped = append(e.wrapped, v)

        case []tContent:
            for _, content := range v {
                e.wrapped = append(e.wrapped, content)
            }

        default:
            e.wrapped = append(e.wrapped, content(v))

        }
    }
    return e
}