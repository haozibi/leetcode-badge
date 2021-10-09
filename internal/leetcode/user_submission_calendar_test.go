package leetcode

import "testing"

func TestGetSubCal(t *testing.T) {

	var err error

	if _, err = GetSubCal("haozibi"); err != nil {
		t.Error(err)
	}

	if _, err = GetSubCal("zerotrac2"); err != nil {
		t.Error(err)
	}
}
