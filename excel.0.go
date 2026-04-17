package main

import (
	"fmt"
	"net/http"
)

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	qvars := QuoteVars(&state)
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}
