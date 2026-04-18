package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
)

const sessionCookie = `quo2_session`

type tSessionCtxKey struct{}

type SessionStore_t struct {
	mu sync.RWMutex
	byToken map[string]SessionVars_t
}

var sessionCtxKey = tSessionCtxKey{}

func NewSessionStore() *SessionStore_t {
	return &SessionStore_t{
		byToken: make(map[string]SessionVars_t),
	}
}

func InitSessionVars() SessionVars_t {
	return SessionVarsFromState(InitState())
}

func SessionVarsFromState(state State_t) SessionVars_t {
	work := InitState()
	work.user = state.user
	work.quote = CloneUIBagVars(state.quote)
	return SessionVars_t{
		user: work.user,
		quote: QuoteVars(&work),
	}
}

func StateFromSessionVars(vars SessionVars_t) State_t {
	out := QuoteStateFromQuoteVars(vars.quote)
	out.user = vars.user
	return out
}

func NewSessionToken() string {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil { panic(e) }
	return hex.EncodeToString(b)
}

func (x *SessionStore_t)EnsureToken(raw string) (token string, setCookie bool) {
	token = strings.TrimSpace(raw)

	x.mu.Lock()
	defer x.mu.Unlock()

	if token == `` {
		setCookie = true
		for {
			token = NewSessionToken()
			if _, ok := x.byToken[token]; !ok { break }
		}
	}
	if _, ok := x.byToken[token]; !ok {
		x.byToken[token] = InitSessionVars()
	}
	return token, setCookie
}

func (x *SessionStore_t)GetSessionVars(token string) SessionVars_t {
	token = strings.TrimSpace(token)
	if token == `` { return InitSessionVars() }

	x.mu.RLock()
	vars, ok := x.byToken[token]
	x.mu.RUnlock()
	if ok { return vars }
	return InitSessionVars()
}

func (x *SessionStore_t)GetState(token string) State_t {
	return StateFromSessionVars(x.GetSessionVars(token))
}

func (x *SessionStore_t)SetState(token string, state State_t) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	vars := SessionVarsFromState(state)

	x.mu.Lock()
	x.byToken[token] = vars
	x.mu.Unlock()
}

func (x *SessionStore_t)Destroy(token string) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	x.mu.Lock()
	delete(x.byToken, token)
	x.mu.Unlock()
}

func SessionToken(r *http.Request) string {
	v := r.Context().Value(sessionCtxKey)
	token, ok := v.(string)
	if !ok { return `` }
	return strings.TrimSpace(token)
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := ``
		cookie, _ := r.Cookie(sessionCookie)
		if cookie != nil { raw = cookie.Value }

		token, setCookie := App.sessionStore.EnsureToken(raw)
		if setCookie { SetSessionCookie(w, token) }

		ctx := context.WithValue(r.Context(), sessionCtxKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DestroySession(r *http.Request) { App.sessionStore.Destroy(SessionToken(r)) }

func SetState(r *http.Request, state State_t) { App.sessionStore.SetState(SessionToken(r), state) }

func GetState(r *http.Request) State_t { return App.sessionStore.GetState(SessionToken(r)) }

func SetSessionCookie(w http.ResponseWriter, token string) {
	token = strings.TrimSpace(token)
	if token == `` { return }
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie,
		Value: token,
		Path: `/`,
		MaxAge: 60 * 60 * 24 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie,
		Value: ``,
		Path: `/`,
		MaxAge: -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
