package request

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// SendRequest send request
func SendRequest(req *http.Request) ([]byte, *http.Response, error) {

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp, errors.Errorf("StatusCode not eq 200: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read body")
	}

	return body, resp, nil
}
