package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"

	. "quo2/lib/output"
)

func main() {
	Bootstrap()
	defer App.DB.Close()
	if true {
		SetForcedQuoteDefaults()
	}

	r := chi.NewRouter()
	r.Use(SessionMiddleware)
	r.Get(`/`, App.Auth(Page2EditQ))

	r.Get(`/quote`, App.Auth(Page1Quote))
	r.Post(`/quote-info-change`, App.Auth(Page1QuoteChange))
	
	r.Get(`/quote-review`, App.Auth(Page2EditQ))
	r.Post(`/quote-review`, App.Auth(Page2EditQEntry))
	r.Post(`/quote-review-change`, App.Auth(Page2EditQChange))

	r.Get(`/signin`, SignInHandler)
	r.Post(`/signin`, SignInHandler)
	r.Get(`/signout`, SignOutHandler)

	r.Handle(`/static/*`, http.StripPrefix(`/static/`, http.FileServer(http.Dir(`./static`))))

	Log(`quo2 is running on :`, App.port)
	log.Fatal(http.ListenAndServe(Str(`:`, App.port), r))
}
