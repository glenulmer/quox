package main

import (
	"net/http"

	. "quo2/lib/output"
)

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	vars := QuoteVars(&state)
	slim := IsSlimXl(req)

	path, fname, ok := CreateXlQuote(vars, slim)
	if !ok {
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, path, fname, xlsx)
}
