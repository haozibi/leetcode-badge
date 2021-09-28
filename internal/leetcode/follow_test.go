package leetcode

import (
	"fmt"
	"testing"
)

func TestGetFollow(t *testing.T) {

	var (
		name1 = "haozibi"
		name2 = "haozibi2"
		a     = -1
		b     = -1
		err   error
	)

	a, b, err = GetFollow(name1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(a, b)

	a, b, err = GetFollow(name2)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(a, b)
}
