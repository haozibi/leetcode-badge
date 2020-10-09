package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// Badge badge
func (a *APP) Badge(badgeType BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {

	info, err := a.getUserProfile(name, isCN)
	if err != nil {
		log.Err(err).
			Str("BadgeType", badgeType.String()).
			Str("Name", name).
			Bool("IsCN", isCN).
			Msg("get user profile error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := a.getBadge(r.URL.Query(), isCN, badgeType, info)
	if err != nil {
		log.Err(err).
			Str("BadgeType", badgeType.String()).
			Str("Name", name).
			Bool("IsCN", isCN).
			Msg("get badge error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
}

func (a *APP) Chart(badgeType BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {
	day := 7
	dayStr := r.URL.Query().Get("day")
	if dayStr != "" {
		d, err := strconv.Atoi(dayStr)
		if err != nil {
			_, _ = fmt.Fprintf(w, "day params must number")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if d > 30 {
			d = 30
		}
		if d > day {
			day = d
		}
	}

	end := time.Now()
	start := time.Now().AddDate(0, 0, -1*(day-1))

	body, err := a.historyChart(badgeType, name, isCN, start, end)
	if err != nil {
		log.Err(err).
			Str("BadgeType", badgeType.String()).
			Str("Name", name).
			Bool("IsCN", isCN).
			Msg("get chart error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
}
