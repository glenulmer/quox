module quo2

go 1.25.3

require (
	github.com/alexedwards/scs/v2 v2.9.0
	github.com/go-chi/chi/v5 v5.0.8
	golang.org/x/crypto v0.44.0
	pm v0.0.0
)

require github.com/go-sql-driver/mysql v1.7.1 // indirect

replace pm => ../pm
