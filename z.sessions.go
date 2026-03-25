package main

import (
	"context"
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

const sessionStateKey = `state`
const sessionCookie = `quo2_session`

func NewSessionManager() *scs.SessionManager {
	gob.Register(InitState())

	manager := scs.New()
	manager.Store = memstore.New()
	manager.Lifetime = 365 * 24 * time.Hour
	manager.IdleTimeout = 0
	return manager
}

type tSessionWriter struct {
	http.ResponseWriter
	manager *scs.SessionManager
	ctx context.Context
	done bool
}

func (x *tSessionWriter) WriteHeader(code int) {
	x.commitSession()
	x.ResponseWriter.WriteHeader(code)
}

func (x *tSessionWriter) Write(b []byte) (int, error) {
	x.commitSession()
	return x.ResponseWriter.Write(b)
}

func (x *tSessionWriter) commitSession() {
	if x.done { return }
	x.done = true
	token, _, _ := x.manager.Commit(x.ctx)
	if token != `` {
		SetSessionCookie(x.ResponseWriter, token)
	}
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ``
		cookie, _ := r.Cookie(sessionCookie)
		if cookie != nil { token = strings.TrimSpace(cookie.Value) }
		ctx, e := App.sessionManager.Load(r.Context(), token)
		if e != nil {
			http.Error(w, `session load failed`, http.StatusBadRequest)
			return
		}

		if !App.sessionManager.Exists(ctx, `session_started_unix`) {
			App.sessionManager.Put(ctx, `session_started_unix`, time.Now().Unix())
		}
		SessionConfig(ctx)

		sw := &tSessionWriter{ ResponseWriter: w, manager: App.sessionManager, ctx: ctx }
		next.ServeHTTP(sw, r.WithContext(ctx))
		sw.commitSession()
	})
}

func SessionConfig(ctx context.Context) {
	if App.sessionManager.Exists(ctx, sessionStateKey) { return }
	App.sessionManager.Put(ctx, sessionStateKey, InitState())
}

func SessionSet(key string, value any, r *http.Request) { App.sessionManager.Put(r.Context(), key, value) }

func SessionGetInt(key string, r *http.Request) int { return App.sessionManager.GetInt(r.Context(), key) }

func SetState(r *http.Request, state State_t) { App.sessionManager.Put(r.Context(), sessionStateKey, state) }

func GetState(r *http.Request) State_t {
	v := App.sessionManager.Get(r.Context(), sessionStateKey)
	state, ok := v.(State_t)
	if ok { return state }
	return InitState()
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
