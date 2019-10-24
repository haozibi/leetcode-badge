package app

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

var (
	BuildTime    = ""
	BuildVersion = ""
	BuildAppName = "lcbadge"
	CommitHash   = ""
)

// Version show version
func (a *APP) Version(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("X-Lc-AppName", os.Getenv("APPNAME"))

	w.Header().Set("X-Lc-BuildTime", BuildTime)
	w.Header().Set("X-Lc-BuildVersion", BuildVersion)
	w.Header().Set("X-Lc-BuildName", BuildAppName)

	w.Header().Set("X-User-Num", strconv.Itoa(len(a.userMap)))
	w.Header().Set("X-Chart-Num", strconv.Itoa(len(a.recordMap)))
	io.WriteString(w, "lc ok")
}
