package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

func (a *APP) addUser(name string, isCN bool) {
	t := time.Now()

	key := name + "_" + fmt.Sprintf("%v", isCN)

	a.userInfoMu.Lock()
	a.userInfo[key] = t
	a.userInfoMu.Unlock()
}

func (a *APP) runInterHTTP() error {
	var (
		port = a.config.DebugAddress
	)

	r := mux.NewRouter()
	r.HandleFunc("/user_info", a.ShowUserInfo)

	return a.runhttp(fmt.Sprintf("127.0.0.1:%d", port), r)
}

type MetricUserInfo struct {
	Name string `json:"name"`
	// IsCN string `json:"is_cn"`
	Time string `json:"last_time"`
}

func (a *APP) ShowUserInfo(w http.ResponseWriter, r *http.Request) {

	a.userInfoMu.Lock()
	list := make([]MetricUserInfo, 0, len(a.userInfo))
	for k, v := range a.userInfo {
		list = append(list, MetricUserInfo{
			Name: k,
			Time: v.Format("2006-01-02 15:04:05"),
		})
	}
	a.userInfoMu.Unlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	body, _ := json.Marshal(list)
	w.Write(body)
}
