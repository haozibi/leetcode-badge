package chart

import (
	"os"
	"testing"
	"time"
)

func TestShowRankHistory(t *testing.T) {

	f, _ := os.Create("output.svg")
	defer f.Close()

	r1 := make([]RankHistory, 0)
	r2 := make([]RankHistory, 0)

	r1 = append(r1, RankHistory{
		Date: time.Now(),
		Rank: 29999,
	})
	r1 = append(r1, RankHistory{
		Date: time.Now().AddDate(0, 0, 1),
		Rank: 26989,
	})
	r1 = append(r1, RankHistory{
		Date: time.Now().AddDate(0, 0, 2),
		Rank: 22999,
	})
	r1 = append(r1, RankHistory{
		Date: time.Now().AddDate(0, 0, 3),
		Rank: 20321,
	})
	r1 = append(r1, RankHistory{
		Date: time.Now().AddDate(0, 0, 5),
		Rank: 19321,
	})

	r2 = append(r2, RankHistory{
		Date: time.Now(),
		Rank: 26989,
	})
	r2 = append(r2, RankHistory{
		Date: time.Now().AddDate(0, 0, 1),
		Rank: 16321,
	})
	r2 = append(r2, RankHistory{
		Date: time.Now().AddDate(0, 0, 2),
		Rank: 20321,
	})
	r2 = append(r2, RankHistory{
		Date: time.Now().AddDate(0, 0, 3),
		Rank: 16321,
	})
	r2 = append(r2, RankHistory{
		Date: time.Now().AddDate(0, 0, 6),
		Rank: 18321,
	})

	ShowRankHistory(f, [][]RankHistory{r1, r2}, "abc", "uuu")
}

func TestSameValue(t *testing.T) {

	f, _ := os.Create("output.svg")
	defer f.Close()

	r1 := make([]SolvedHistory, 0)
	r1 = append(r1, SolvedHistory{
		Date: time.Now(),
		Num:  62,
	})
	r1 = append(r1, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 1),
		Num:  62,
	})

	ShowSolvedHistory(f, [][]SolvedHistory{r1}, "aaa")
}
