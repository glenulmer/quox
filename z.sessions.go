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

const sessionHeader = `X-Session-Token`
const sessionStateKey = `state`

func NewSessions() Sessions_t {
	gob.Register(InitState())

	manager := scs.New()
	manager.Store = memstore.New()
	manager.Lifetime = 365 * 24 * time.Hour
	manager.IdleTimeout = 0
	return Sessions_t{
		manager: manager,
		header: sessionHeader,
	}
}

type tSessionWriter struct {
	http.ResponseWriter
	sessions Sessions_t
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
	token, _, e := x.sessions.manager.Commit(x.ctx)
	if e != nil { return }
	if token != `` { x.Header().Set(x.sessions.header, token) }
}

func (x Sessions_t) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.Header.Get(x.header))
		ctx, e := x.manager.Load(r.Context(), token)
		if e != nil {
			http.Error(w, `session load failed`, http.StatusBadRequest)
			return
		}

		if !x.manager.Exists(ctx, `session_started_unix`) {
			x.manager.Put(ctx, `session_started_unix`, time.Now().Unix())
		}
		SessionConfig(ctx)

		sw := &tSessionWriter{ ResponseWriter: w, sessions: x, ctx: ctx }
		next.ServeHTTP(sw, r.WithContext(ctx))
		sw.commitSession()
	})
}

func SessionConfig(ctx context.Context) {
	if App.sessions.manager.Exists(ctx, sessionStateKey) { return }
	App.sessions.manager.Put(ctx, sessionStateKey, InitState())
}

func SessionSet(key string, value any, r *http.Request) {
	App.sessions.manager.Put(r.Context(), key, value)
}

func SessionGetInt(key string, r *http.Request) int {
	return App.sessions.manager.GetInt(r.Context(), key)
}

func GetState(r *http.Request) State_t {
	v := App.sessions.manager.Get(r.Context(), sessionStateKey)
	state, ok := v.(State_t)
	if ok { return state }
	return InitState()
}

func SetState(r *http.Request, state State_t) {
	App.sessions.manager.Put(r.Context(), sessionStateKey, state)
}
