package leetcode

import (
	"fmt"
	"testing"
)

func TestGetFollow(t *testing.T) {

	var (
		name1 = "haozibi"
		name2 = "haozibi2"
	)

	a, err := GetFollow(name1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(a)

	b, err := GetFollow(name2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(b)
}
