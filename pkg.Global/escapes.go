package globals
import (
	. `klpm/lib/output`
)

func EscapeSelector(s string) string {
    // Handle quotes in selectors for JSON string embedding
    return Replace(s, `"`, `\"`)
}

func EscapeHTML(s string) string {
    s = Replace(s, `\`, `\\`)
    s = Replace(s, `"`, `\"`)     // for JSON string embedding
    s = Replace(s, "<", `\u003c`) // HTML entities
    s = Replace(s, ">", `\u003e`)
    s = Replace(s, "\n", `\n`)
    return s
}
