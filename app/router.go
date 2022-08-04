package app

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Router(r *mux.Router, a *APP, w io.Writer) {

	r.HandleFunc("/version", a.Version)

	apiCNV1 := r.PathPrefix("/v1cn").Subrouter()
	apiENV1 := r.PathPrefix("/v1").Subrouter()

	router := make([]*mux.Router, 0, 5)
	router = append(router, apiCNV1)
	router = append(router, apiENV1)

	for k, api := range router {

		isCN := k == 0

		// [basic] 排名
		api.Methods(http.MethodGet).Path("/ranking/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeRanking, isCN)),
		)

		// [basic] 通过的题目/问题总数
		api.Methods(http.MethodGet).Path("/solved/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeSolved, isCN)),
		)

		// [basic] 通过的题目/问题总数 rate
		api.Methods(http.MethodGet).Path("/solved-rate/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeSolvedRate, isCN)),
		)

		// [basic] 通过提交/提交的总数
		api.Methods(http.MethodGet).Path("/accepted/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeAccepted, isCN)),
		)

		// [basic] 通过提交/提交的总数 rate
		api.Methods(http.MethodGet).Path("/accepted-rate/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeAcceptedRate, isCN)),
		)

		// // [badge] 关注
		// api.Methods(http.MethodGet).Path("/badge/following/{name:.+}.svg").Handler(
		// 	handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeFollowing, isCN)),
		// )
		//
		// // [badge] 被关注数
		// api.Methods(http.MethodGet).Path("/badge/followers/{name:.+}.svg").Handler(
		// 	handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeFollowers, isCN)),
		// )

		// [chart] 排名记录图表
		api.Methods(http.MethodGet).Path("/chart/ranking/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeChartRanking, isCN)),
		)

		// [chart] 答题数量图表
		api.Methods(http.MethodGet).Path("/chart/solved/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeChartSolved, isCN)),
		)

		// [chart] 获得答题日历
		api.Methods(http.MethodGet).Path("/chart/submission-calendar/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeChartSubmissionCalendar, isCN)),
		)

		// [basic] 获得个人信息
		api.Methods(http.MethodGet).Path("/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, a.HandlerFunc(BadgeTypeProfile, isCN)),
		)
	}
	r.PathPrefix("/").Handler(http.HandlerFunc(IndexPage))
}

func (a *APP) HandlerFunc(badgeType BadgeType, isCN bool) http.Handler {

	var f func(badgeType BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request)

	switch badgeType {
	case BadgeTypeProfile, BadgeTypeRanking, BadgeTypeSolved, BadgeTypeSolvedRate, BadgeTypeAccepted, BadgeTypeAcceptedRate:
		f = a.Basic
	case BadgeTypeChartRanking, BadgeTypeChartSolved:
		f = a.Chart
	// case BadgeTypeFollowing, BadgeTypeFollowers:
	// 	f = a.Badge
	case BadgeTypeChartSubmissionCalendar:
		f = a.SubCal
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if f == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		vars := mux.Vars(r)
		name := vars["name"]
		if name == "" || f == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		f(badgeType, name, isCN, w, r)
	})
}

// IndexPage index page
func IndexPage(w http.ResponseWriter, r *http.Request) {

	githubPage := "https://github.com/haozibi/leetcode-badge"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<meta http-equiv=refresh content=0;url="%s">`, githubPage)
}
