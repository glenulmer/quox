package main

import "net/http"

func DownloadExcel(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, `/quote-review`, http.StatusSeeOther)
}
