package globals
import (
	. `quo2/lib/output`
    "regexp"
)

var removeTags *regexp.Regexp

func init() {
    removeTags = regexp.MustCompile(`<[^>]*>`)
}

func EscapeSelector(s string) string {
    // Handle quotes in selectors for JSON string embedding
    return Replace(s, `"`, `\"`)
}

func DQtoQuot(s string) string   { return Replace(s, `"`, `&quot;`) }
func QuottoDQ(s string) string { return Replace(s, `&quot;`, `"`) }

func EscapeHTML(s string) string {
    s = Replace(s, `\`, `\\`)
    s = Replace(s, `"`, `\"`)     // for JSON string embedding
    s = Replace(s, "<", `\u003c`) // HTML entities
    s = Replace(s, ">", `\u003e`)
    s = Replace(s, "\n", `\n`)
    return s
}

func RemoveTags(s string) string {
    return removeTags.ReplaceAllString(s, "")
}

func RewriteHtml(s string) string {
    s = Replace(s, `&`, `&amp;`)
    s = Replace(s, `<`, `&lt;`)
    s = Replace(s, `>`, `&gt;`)
    s = Replace(s, `"`, `&quot;`)
    s = Replace(s, `'`, `&#39;`)
    s = Replace(s, "<br>", `\n`)
    return s
/*
    return t
        .replace(/&/g, '')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;')
        .replace(/\n/g, '<br>');
*/
}
