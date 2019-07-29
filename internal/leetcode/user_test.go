package leetcode

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetUserProfile(t *testing.T) {

	tests := []struct {
		userName   string
		wantErrStr string
	}{
		{"haozibi", ""},
		// {"haozibi222", "not found haozibi222"},
	}

	for k, v := range tests {

		got, err := GetUserProfile(v.userName, true)
		if err != nil {
			if err.Error() != v.wantErrStr {
				t.Errorf("k: %v,want: %v got: %v\n", k, v.wantErrStr, err)
			}
		}
		spew.Dump(got)
	}
}

func TestGetUserProfileEN(t *testing.T) {

	info, err := getUserProfile("haozi2123123bi")
	if err != nil {
		panic(err)
	}

	spew.Dump(info)
}
