package leetcodecn

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetUserQuestionProgress(t *testing.T) {
	names := []string{
		"haozibi",
		//"oooooooooooxxxxxx",
		//"ac_oier",
	}

	for _, v := range names {
		p, err := GetUserQuestionProgress(v)
		fmt.Println(err)
		spew.Dump(p)
	}
}
