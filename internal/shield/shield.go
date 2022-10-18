package shield

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/request"
)

// Badge get badge
// https://img.shields.io/badge/leetcode-haozibi-green.svg?color=red
func Badge(query url.Values, left, right, color string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	param := encoding(left) + "-" + encoding(right) + "-" + color
	uri := "https://img.shields.io/badge/" + param
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http request")
	}
	req.URL.RawQuery = query.Encode()

	body, _, err := request.SendRequest(http.DefaultClient, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Dashes --	→	- Dash
// Underscores __	→	_ Underscore
// _ or Space  	→	  Space
func encoding(s string) string {
	s = strings.ReplaceAll(s, "-", "--")
	s = strings.ReplaceAll(s, "_", "__")
	return s
}
