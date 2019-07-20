package confluence

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/naminomare/gogutil/network"
)

// https://docs.atlassian.com/atlassian-confluence/REST/6.6.0/
// を参考に

// PageType ページタイプ
type PageType string

var (
	// PageTypePage page
	PageTypePage PageType = "page"

	// PageTypeBlog blog
	PageTypeBlog PageType = "blog"
)

// Client コンフルアクセス用クライアント
type Client struct {
	baseURL    string
	httpClient *network.HTTPWaitClient
}

// NewClient クライアント作成
func NewClient(
	baseURL,
	serverName,
	userName,
	password string,
) *Client {
	ret := Client{
		baseURL:    baseURL,
		httpClient: network.NewHTTPWaitClient(1000, serverName),
	}
	ret.httpClient.SetAuth(userName, password)
	return &ret
}

// CreateContent コンテンツ作成
func (t *Client) CreateContent(
	spaceKey,
	ancestorsID,
	title,
	content string,
	pagetype PageType,
) (*http.Response, error) {
	targetURL := t.baseURL + "/content"
	postMap := map[string]interface{}{
		"type":  pagetype,
		"title": title,
		"ancestors": []interface{}{
			map[string]string{
				"id": ancestorsID,
			},
		},
		"space": map[string]string{
			"key": spaceKey,
		},
		"body": map[string]interface{}{
			"storage": map[string]string{
				"value":          content,
				"representation": "storage",
			},
		},
	}
	reader := toJSONReader(postMap)
	resp, err := t.httpClient.DoRequest(
		http.MethodPost,
		targetURL,
		reader,
		map[string]string{
			network.ContentType: network.ApplicationJSON,
		},
	)
	return resp, err
}

// FetchPage ページ内容を取得する
func (t *Client) FetchPage(
	query map[string]string,
) (*http.Response, error) {
	targetURL := t.baseURL + "/content"
	qStr := ""
	for k, v := range query {
		qStr += url.QueryEscape(k) + "=" + url.QueryEscape(v)
	}
	if qStr != "" {
		targetURL += "?" + qStr
	}

	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		targetURL,
		nil,
		nil,
	)
	return resp, err
}

// FetchPageByID IDでページを取得
func (t *Client) FetchPageByID(ID string) (*http.Response, error) {
	targetURL := t.baseURL + "/content/" + ID
	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		targetURL,
		nil,
		nil,
	)
	return resp, err
}

// AddAttachments ページにファイルを添付する
func (t *Client) AddAttachments(pageID string, files []string) (*http.Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	for _, file := range files {
		fh, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		fw, err := w.CreateFormFile("file", file)
		if err != nil {
			return nil, err
		}
		io.Copy(fw, fh)
	}
	w.Close()

	targetURL := t.baseURL + "/content/" + pageID + "/child/attachment"
	resp, err := t.httpClient.DoRequest(
		http.MethodPost,
		targetURL,
		&buf,
		map[string]string{
			network.ContentType: w.FormDataContentType(),
			"X-Atlassian-Token": "no-check",
		},
	)
	return resp, err
}

func toJSONReader(mapobj map[string]interface{}) *bytes.Reader {
	bin, err := json.Marshal(mapobj)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(bin)

	return reader
}

// UpdatePageByID IDでページを更新
// func (t *Client) UpdatePageByID(ID string) (*http.Response, error) {
// 	targetURL := t.baseURL + "/content/" + ID

// 	resp, err := t.FetchPageByID(ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	respMap, err := network.ResponseToMap(resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	respMap[]
// }
