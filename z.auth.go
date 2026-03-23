package main

import "net/http"

func FakeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := GetState(r)
		if !state.LoggedIn() {
			state.user = UserInfo_t{ 1, `Glen`, `glen.ulmer@gmail.com` }
			SetState(r, state)
		}
		next.ServeHTTP(w, r)
	}
}

func TrueAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !GetState(r).LoggedIn() {
			http.Redirect(w, r, `/signin`, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
