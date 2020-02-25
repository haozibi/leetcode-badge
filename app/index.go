package app

import (
	"fmt"
	"net/http"
)

// IndexPage index page
func (a *APP) IndexPage(w http.ResponseWriter, r *http.Request) {

	githubPage := "https://github.com/haozibi/leetcode-badge"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<meta http-equiv=refresh content=0;url="%s">`, githubPage)
}
