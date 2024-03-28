package leetcodecn

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func SendRaw(client *http.Client, uri string, method string, query string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, uri, strings.NewReader(query))
	if err != nil {
		return nil, errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
	}

	req.Header.Add("origin", "https://leetcode.cn")
	req.Header.Add("user-agent", getUserAgent())
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", "https://leetcode.cn")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("status code error, code: %d, uri: %s, method: %s, query: %s", resp.StatusCode, uri, method, query)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return body, errors.Wrapf(err, "uri: %s, method: %s, query: %s", uri, method, query)
}

func Send(client *http.Client, uri string, method string, query string, ptr interface{}) error {
	if t := reflect.ValueOf(ptr); t.Kind() != reflect.Ptr {
		return errors.Errorf("ptr must be ptr")
	}

	body, err := SendRaw(client, uri, method, query)
	if err != nil {
		return err
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

var uaList = []string{
	"Mozilla/5.0 (Linux; U; Android 5.1.1; zh-cn; Redmi 3 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/42.0.0.0 Mobile Safari/537.36 XiaoMi/MiuiBrowser/2.1.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.93 Safari/537.36",
	"Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14931",
	"Chrome (AppleWebKit/537.1; Chrome50.0; Windows NT 6.3) AppleWebKit/537.36 (KHTML like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/7046A194A",
	"Mozilla/5.0 (iPad; CPU OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5355d Safari/8536.25",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.117 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.0 Mobile/14G60 Safari/602.1",
}

func getUserAgent() string {

	i := rand.Intn(len(uaList))

	return uaList[i]
}
