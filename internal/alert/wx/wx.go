package wx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/haozibi/leetcode-badge/internal/alert"
)

type wx struct {
	uri string
}

func New(uri string) alert.Alert {
	return &wx{uri: uri}
}

func (w *wx) Send(header, body string) error {
	content := fmt.Sprintf("# %s\n\n%s", header, body)
	return send(w.uri, content)
}

func send(uri string, content string) error {
	n := make(map[string]interface{})
	n["content"] = content

	m := make(map[string]interface{})
	m["msgtype"] = "markdown"
	m["markdown"] = n

	body, _ := json.Marshal(m)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return errors.Wrapf(err, "uri: %s, content: %s", uri, content)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "uri: %s, content: %s", uri, content)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("status code error, code: %d, uri: %s, content: %s", resp.StatusCode, uri, content)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "uri: %s, content: %s", uri, content)
	}

	var p Resp
	if err = json.Unmarshal(respBody, &p); err != nil {
		return errors.Wrapf(err, "uri: %s, body: %s", uri, string(respBody))
	}

	if p.Code != 0 {
		return errors.Errorf("response code error, code: %d, message: %s", p.Code, p.Message)
	}
	return nil
}

type Resp struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}
