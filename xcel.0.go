package main

import (
	"net/http"

	. "klpm/lib/output"
)

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	vars := QuoteVars(&state)
	slim := IsSlimXl(req)
	lang := vars.lang
	if lang <= 0 { lang = English }
	if x := Atoi(req.FormValue(`lang`)); x > 0 { lang = LangId_t(x) }
	vars.lang = lang

	path, fname, ok := CreateXlQuote(vars, slim, lang)
	if !ok {
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, path, fname, xlsx)
}
