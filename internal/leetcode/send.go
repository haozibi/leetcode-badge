package leetcode

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func Send(client *http.Client, uri string, method string, query string, ptr interface{}) error {
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr {
		return errors.Errorf("ptr must be ptr")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, uri, strings.NewReader(query))
	if err != nil {
		return errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
	}

	req.Header.Add("origin", "https://leetcode-cn.com")
	req.Header.Add("user-agent", GetUserAgent())
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", "https://leetcode-cn.com")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("status code error, code: %d, uri: %s, method: %s, query: %s", resp.StatusCode, uri, method, query)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
	}

	var p Data
	if err = json.Unmarshal(body, &p); err != nil {
		return errors.Wrapf(err, "uri: %s, method: %s, resp body: %s", uri, method, string(body))
	}

	// ugly
	if len(p.Errors) != 0 {
		ee := ""
		for _, v := range p.Errors {
			if v.Message == ErrUserNotExist.Error() {
				return nil
			}
			ee += " " + v.Message
		}
		return errors.Errorf("error: %v, uri: %s, method: %s, query: %s", ee, uri, method, query)
	}

	err = json.Unmarshal(p.Data, ptr)
	return errors.Wrapf(err, "uri: %s, method: %s, resp data: %s", uri, method, string(p.Data))
}
