package chart

import (
	"os"
	"testing"
	"time"
)

func TestShowSolvedHistory(t *testing.T) {

	f, _ := os.Create("output.svg")
	defer f.Close()

	r1 := make([]SolvedHistory, 0)
	r2 := make([]SolvedHistory, 0)

	r1 = append(r1, SolvedHistory{
		Date: time.Now(),
		Num:  1,
	})
	r1 = append(r1, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 1),
		Num:  10,
	})
	r1 = append(r1, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 2),
		Num:  29,
	})
	r1 = append(r1, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 3),
		Num:  35,
	})
	r1 = append(r1, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 4),
		Num:  55,
	})

	r2 = append(r2, SolvedHistory{
		Date: time.Now(),
		Num:  2,
	})
	r2 = append(r2, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 1),
		Num:  23,
	})
	r2 = append(r2, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 2),
		Num:  23,
	})
	r2 = append(r2, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 3),
		Num:  23,
	})
	r2 = append(r2, SolvedHistory{
		Date: time.Now().AddDate(0, 0, 4),
		Num:  31,
	})

	ShowSolvedHistory(f, [][]SolvedHistory{r1, r2}, "abc", "uuu")
}
