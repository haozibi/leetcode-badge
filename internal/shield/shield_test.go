package shield

import (
	"fmt"
	"net/url"
	"testing"
)

func TestBadge(t *testing.T) {
	names := []string{
		"mr-j001",
		"ac_oier",
	}

	for _, v := range names {
		b, err := Badge(url.Values{}, "LeetCode CN", v, "green")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(b))
	}

}
