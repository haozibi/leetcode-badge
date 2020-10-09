package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	BuildTime    = ""
	BuildVersion = ""
	BuildAppName = ""
	CommitHash   = ""
)

// Version Show version
func (a *APP) Version(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("X-Lc-BuildEnv", os.Getenv("LC_BUILD_ENV"))
	w.Header().Set("X-Lc-BuildTime", BuildTime)
	w.Header().Set("X-Lc-BuildVersion", BuildVersion)
	w.Header().Set("X-Lc-BuildName", BuildAppName)
	w.Header().Set("X-Lc-CommitHash", CommitHash)
	_, _ = io.WriteString(w, Version())
}

func Version() string {
	return fmt.Sprintf("%s \ntag: %s\nbuild: %s\nhash: %s\n",
		BuildAppName,
		BuildVersion,
		BuildTime,
		CommitHash,
	)
}
