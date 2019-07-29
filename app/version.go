package app

import (
	"io"
	"net/http"
)

var (
	BuildTime    = ""
	BuildVersion = ""
	BuildAppName = "lcbadge"
)

func (a *APP) Version(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("X-Lc-BuildTime", BuildTime)
	w.Header().Set("X-Lc-BuildVersion", BuildVersion)
	w.Header().Set("X-Lc-BuildName", BuildAppName)

	io.WriteString(w, "lc ok")
}
