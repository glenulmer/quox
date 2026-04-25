package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"

	"github.com/mssola/useragent"
)

const sessionCookie = `klpm_session`
const deviceCookie = `device`
const deviceMobile, deviceDesktop = `mobile`, `desktop`

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
		device: deviceDesktop,
		deviceConfirmed: false,
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

func NormalizeDeviceMode(raw string) (mode string, ok bool) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case deviceMobile:
		return deviceMobile, true
	case deviceDesktop:
		return deviceDesktop, true
	}
	return ``, false
}

func UAMode(r *http.Request) string {
	if useragent.New(r.UserAgent()).Mobile() { return deviceMobile }
	return deviceDesktop
}

func (x *SessionStore_t)EnsureToken(raw string) (token string, setCookie bool, created bool) {
	token = strings.TrimSpace(raw)

	x.mu.Lock()
	defer x.mu.Unlock()

	if token == `` {
		setCookie = true
		created = true
		for {
			token = NewSessionToken()
			if _, ok := x.byToken[token]; !ok { break }
		}
	}
	if _, ok := x.byToken[token]; !ok {
		x.byToken[token] = InitSessionVars()
		created = true
	}
	return token, setCookie, created
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
	if prior, ok := x.byToken[token]; ok {
		vars.deviceConfirmed = prior.deviceConfirmed
		vars.device = deviceDesktop
		if mode, modeOK := NormalizeDeviceMode(prior.device); modeOK { vars.device = mode }
	}
	x.byToken[token] = vars
	x.mu.Unlock()
}

func (x *SessionStore_t)GetDevice(token string) (mode string, confirmed bool) {
	token = strings.TrimSpace(token)
	if token == `` { return deviceDesktop, false }

	x.mu.RLock()
	vars, ok := x.byToken[token]
	x.mu.RUnlock()
	if !ok { return deviceDesktop, false }

	if mode, ok = NormalizeDeviceMode(vars.device); ok {
		return mode, vars.deviceConfirmed
	}
	return deviceDesktop, vars.deviceConfirmed
}

func (x *SessionStore_t)SetDevice(token, mode string, confirmed bool) {
	token = strings.TrimSpace(token)
	if token == `` { return }

	m := deviceDesktop
	if mode0, ok := NormalizeDeviceMode(mode); ok { m = mode0 }

	x.mu.Lock()
	vars, ok := x.byToken[token]
	if !ok { vars = InitSessionVars() }
	vars.device = m
	vars.deviceConfirmed = confirmed
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

		token, setCookie, created := App.sessionStore.EnsureToken(raw)
		if setCookie { SetSessionCookie(w, token) }

		if c, _ := r.Cookie(deviceCookie); c != nil {
			if mode, ok := NormalizeDeviceMode(c.Value); ok {
				App.sessionStore.SetDevice(token, mode, true)
			}
		}

		mode, confirmed := App.sessionStore.GetDevice(token)
		if !confirmed && created {
			if mode == deviceDesktop {
				App.sessionStore.SetDevice(token, UAMode(r), false)
			}
		}

		ctx := context.WithValue(r.Context(), sessionCtxKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DestroySession(r *http.Request) { App.sessionStore.Destroy(SessionToken(r)) }

func SetState(r *http.Request, state State_t) { App.sessionStore.SetState(SessionToken(r), state) }

func GetState(r *http.Request) State_t { return App.sessionStore.GetState(SessionToken(r)) }

func SessionDeviceMode(r *http.Request) string {
	mode, _ := App.sessionStore.GetDevice(SessionToken(r))
	return mode
}

func SetSessionDeviceMode(r *http.Request, mode string) {
	token := SessionToken(r)
	_, confirmed := App.sessionStore.GetDevice(token)
	App.sessionStore.SetDevice(token, mode, confirmed)
}

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
