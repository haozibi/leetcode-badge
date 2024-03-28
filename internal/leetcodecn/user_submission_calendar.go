package leetcodecn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func GetSubCal(name string) (map[string]int, error) {
	if name == "" {
		return nil, errors.Errorf("miss user name")
	}

	uri := fmt.Sprintf("https://leetcode.cn/api/user_submission_calendar/%s/", name)

	body, err := SendRaw(http.DefaultClient, uri, http.MethodGet, "")
	if err != nil {
		return nil, err
	}

	body = bytes.Replace(body, []byte(`\"`), []byte(`"`), -1)
	body = bytes.TrimPrefix(body, []byte(`"`))
	body = bytes.TrimSuffix(body, []byte(`"`))

	var p map[string]int

	if err = json.Unmarshal(body, &p); err != nil {
		return nil, errors.Wrapf(err, "body: %s", string(body))
	}

	return p, nil
}
