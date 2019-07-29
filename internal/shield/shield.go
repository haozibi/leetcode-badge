package shield

import (
	"fmt"
	"net/http"

	"github.com/haozibi/leetcode-badge/internal/request"

	"github.com/pkg/errors"
)

// Badge get badge
// https://img.shields.io/badge/leetcode-haozibi-green.svg?color=red
func Badge(r *http.Request, left, right, color string) ([]byte, error) {

	query := r.URL.Query()

	uri := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s", left, right, color)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http request")
	}
	req.URL.RawQuery = query.Encode()

	body, _, err := request.SendRequest(req)
	if err != nil {
		return nil, err
	}

	return body, nil
}
