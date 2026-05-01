package main

import (
	"net/http"
	"os"

	sky "github.com/xuri/excelize/v2"

	. "klpm/lib/output"
)

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	state := GetState(req)
	vars := QuoteVars(&state)

	path, fname, ok := CreateXlQuote(vars)
	if !ok {
		http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
		return
	}

	SendFileToClient(w, path, fname, xlsx)
}

func CreateXlQuote(vars QuoteVars_t) (path, fname string, ok bool) {
	slim := vars.slim
	ex, e := sky.OpenFile(template)
	if e != nil {
		Log(e)
		return ``, ``, false
	}
	defer ex.Close()

	styles := LoadXlStyles(ex)

	if e = WriteXlLayout(ex, styles, vars); e != nil {
		Log(e)
		return ``, ``, false
	}

	path = workDir
	if e = os.MkdirAll(path, 0o775); e != nil {
		Log(e)
		return ``, ``, false
	}

	fname = XlFileName(ClientName(vars), slim)
	_ = ex.DeleteSheet(xlStyleSheet)
	SetXlDefaultCell(ex, quoteSheet, nameCell)
	if e = ex.SaveAs(Str(path, `/`, fname)); e != nil {
		Log(e)
		return ``, ``, false
	}

	return path, fname, true
}
