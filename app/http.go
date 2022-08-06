package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/haozibi/leetcode-badge/internal/statics"
	"github.com/rs/zerolog/log"
)

func (a *APP) Badge(bt BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {
	if !isCN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	info, err := a.userFollow(name)
	if err != nil {
		log.Err(err).
			Str("BadgeType", bt.String()).
			Str("Name", name).
			Bool("IsCN", isCN).
			Msg("get user follow error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := a.followBadge(r.URL.Query(), bt, info)
	if err != nil {
		log.Err(err).
			Str("BadgeType", bt.String()).
			Str("Name", name).
			Bool("IsCN", isCN).
			Msg("get badge error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.write(w, body)
}

// Basic info
func (a *APP) Basic(badgeType BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {

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

	body, err := a.basicBadge(r.URL.Query(), isCN, badgeType, info)
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

// SubCal SubmissionCalendar
func (a *APP) SubCal(_ BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {
	if !isCN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	body, err := a.getSubCal(name, r)
	if err != nil {
		if err == ErrUserNotSupport {
			a.write(w, statics.SVGNotFound())
		} else {
			log.Err(err).
				Str("Name", name).
				Bool("IsCN", isCN).
				Msg("get subcal error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	a.write(w, body)
}

func (a *APP) Card(badgeType BadgeType, name string, isCN bool, w http.ResponseWriter, r *http.Request) {
	if !isCN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
		return
	}

	body, err := a.getCard(badgeType, name, r)
	if err != nil {
		if err == ErrUserNotSupport {
			a.write(w, statics.SVGNotFound())
		} else {
			log.Err(err).
				Str("Name", name).
				Bool("IsCN", isCN).
				Msg("get subcal error")
			w.WriteHeader(http.StatusInternalServerError)
		}
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
