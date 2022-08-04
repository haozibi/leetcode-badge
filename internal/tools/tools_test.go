package tools

import (
	"fmt"
	"testing"
	"time"
)

func TestZeroTime(t *testing.T) {

	now := time.Now()

	day := ZeroTime(now.AddDate(0, 0, -7))
	fmt.Println(day.Unix())
}

type A struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
