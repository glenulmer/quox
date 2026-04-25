package globals

import (
	"net/http"
	. "klpm/lib/output"
)

// include navbar.css

var navbarUserResolver = func(r *http.Request) string {
    session, _ := r.Cookie(`session`)
    if session != nil {
        return session.Value
    }
    return `User`
}

func SetNavbarUserResolver(fn func(*http.Request) string) {
    if fn != nil { navbarUserResolver = fn }
}

func GetUsernameFromCookie(r *http.Request) string {
    return navbarUserResolver(r)
}

func Navbar(req *http.Request, trail ...string) string {
	return `
	<nav class="navbar is-light is-fixed-top" role="navigation" aria-label="main navigation">
	<div class="navbar-brand">
		<a class="navbar-item has-text-link" href="/">KL Editor (Home)</a>` + NL +
		NavTrail(trail...) + NL +
`	</div>
	<div class="navbar-end">
		<div class="navbar-item">
			<span class="mr-2">
` +  GetUsernameFromCookie(req) + `
			</span>
			<a href="/signout" class="button is-light">Log out</a>
		</div>
	</div>
</nav>` + NL
}

func NavTrail(items ...string) string {
	if len(items) == 0 { return `` }
	
	var b Builder
	for _, item := range items {
		b.Add(`<span class="navbar-item has-text-grey-light">|</span>`)
		b.Add(item)
	}
	return b.String()
}
