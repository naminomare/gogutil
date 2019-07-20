package network

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/naminomare/gogutil/timer"
)

var (
	// ContentType Content-Type
	ContentType = "Content-Type"

	// ApplicationJSON application/json
	ApplicationJSON = "application/json"
)

// HTTPWaitClient 一定時間必ず待つ様なクライアント
type HTTPWaitClient struct {
	waitTimer  *timer.WaitTimer
	password   string
	username   string
	servername string
	intervalMS int
}

// NewHTTPWaitClient 一定時間必ず待つ様なクライアントを返す
func NewHTTPWaitClient(intervalMS int, servername string) *HTTPWaitClient {
	return &HTTPWaitClient{
		intervalMS: intervalMS,
		waitTimer:  timer.NewWaitTimer(),
		servername: servername,
	}
}

// SetAuth Authをセット
func (t *HTTPWaitClient) SetAuth(username, password string) {
	t.password = password
	t.username = username
}

// SetServerName サーバー名をいれる。tlsの都合
func (t *HTTPWaitClient) SetServerName(servername string) {
	t.servername = servername
}

func createTLSVerifySkipClient(serverName string) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: serverName, InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

// DoRequest リクエスト
func (t *HTTPWaitClient) DoRequest(
	method,
	url string,
	body io.Reader,
	header map[string]string,
) (*http.Response, error) {
	client := createTLSVerifySkipClient(t.servername)
	req, err := http.NewRequest(method, url, body)
	if t.username != "" {
		req.SetBasicAuth(t.username, t.password)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	t.waitTimer.Wait()
	res, err := client.Do(req)
	t.waitTimer.Start(t.intervalMS)

	return res, err
}

// ResponseToMap httpResponseをmapにして返します
func ResponseToMap(resp *http.Response) (map[string]interface{}, error) {
	bin, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var ret map[string]interface{}
	err = json.Unmarshal(bin, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
