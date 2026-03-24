package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"

	. "pm/lib/output"
)

func main() {
	Bootstrap()
	defer App.DB.Close()

	r := chi.NewRouter()
	r.Use(SessionMiddleware)
	r.Get(`/`, App.Auth(Page0Home))
	r.Get(`/signin`, SignInHandler)
	r.Post(`/signin`, SignInHandler)
	r.Get(`/signout`, SignOutHandler)

	r.Handle(`/static/*`, http.StripPrefix(`/static/`, http.FileServer(http.Dir(`./static`))))

	Log(`quo2 is running on :`, App.port)
	log.Fatal(http.ListenAndServe(Str(`:`, App.port), r))
}
