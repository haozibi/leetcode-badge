package shield

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/request"
)

// Badge get badge
// https://img.shields.io/badge/leetcode-haozibi-green.svg?color=red
func Badge(query url.Values, left, right, color string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uri := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s", left, right, color)

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
