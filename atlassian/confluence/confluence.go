package confluence

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/naminomare/gogutil/fileio"

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

	// ErrInvalidArguments 入力値が不正の時
	ErrInvalidArguments = errors.New("入力値が不正です")
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
	targetURL := t.baseURL + "/rest/api/content"
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
	targetURL := t.baseURL + "/rest/api/content"
	qStr := ""
	for k, v := range query {
		if qStr != "" {
			qStr += "&"
		}
		qStr += url.PathEscape(k) + "=" + url.PathEscape(v)
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
	targetURL := t.baseURL + "/rest/api/content/" + ID
	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		targetURL,
		nil,
		nil,
	)
	return resp, err
}

// FetchContentByTitle タイトルでページのコンテンツを取得
func (t *Client) FetchContentByTitle(spaceKey, title string) (*http.Response, error) {
	targetURL := t.baseURL + "/rest/api/content?spaceKey=" + url.PathEscape(spaceKey) + "&title=" + url.PathEscape(title) + "&expand=body.storage"
	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		targetURL,
		nil,
		nil,
  )
	return resp, err
}

// MovePage ページの移動
func (t *Client) MovePage(srcPageID, dstParentPageID string) (*http.Response, error) {
	targetURL := t.baseURL + "/rest/api/content/" + srcPageID
	resp, err := t.FetchPageByID(srcPageID)
	if err != nil {
		return nil, err
	}
	srcPageMap, err := network.ResponseToMap(resp)
	if err != nil {
		return nil, err
	}
	newver := int(srcPageMap["version"].(map[string]interface{})["number"].(float64) + 1)
	putMap := map[string]interface{}{
		"version": map[string]interface{}{
			"number": newver,
		},
		"type":  srcPageMap["type"],
		"space": srcPageMap["space"],
		"title": srcPageMap["title"],
		"ancestors": []map[string]interface{}{
			map[string]interface{}{
				"id": dstParentPageID,
			},
		},
	}
	reader := toJSONReader(putMap)
	resp, err = t.httpClient.DoRequest(
		http.MethodPut,
		targetURL,
		reader,
		map[string]string{
			network.ContentType: network.ApplicationJSON,
		},
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

		fw, err := w.CreateFormFile("file", fileio.FileName(file))
		if err != nil {
			return nil, err
		}
		io.Copy(fw, fh)
	}
	w.Close()

	targetURL := t.baseURL + "/rest/api/content/" + pageID + "/child/attachment"
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

// AddAttachmentsByIO readerとそれに応じたfilenamesを使って書き込む
// len(readers) != len(filenames) の時は ErrInvalidArgumentsを返す
func (t *Client) AddAttachmentsByIO(pageID string, readers []io.Reader, filenames []string) (*http.Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if len(readers) != len(filenames) {
		return nil, ErrInvalidArguments
	}

	for i, reader := range readers {
		fw, err := w.CreateFormFile("file", fileio.FileName(filenames[i]))
		if err != nil {
			return nil, err
		}
		io.Copy(fw, reader)
	}
	w.Close()

	targetURL := t.baseURL + "/rest/api/content/" + pageID + "/child/attachment"
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

// MoveAttachment pageIDのattachmentIDのattachmentをdstPageIDへ
func (t *Client) MoveAttachment(pageID, attachmentID, dstPageID string) (*http.Response, error) {
	targetURL := t.baseURL + "/rest/api/content/" + pageID + "/child/attachment/" + attachmentID
	jsonObj := map[string]interface{}{
		"id":     attachmentID,
		"type":   "attachment",
		"status": "current",
		"version": map[string]interface{}{
			"number": 1,
		},
		"container": map[string]string{
			"id":   dstPageID,
			"type": "attachment",
		},
	}
	reader := toJSONReader(jsonObj)
	resp, err := t.httpClient.DoRequest(
		http.MethodPut,
		targetURL,
		reader,
		map[string]string{
			network.ContentType: network.ApplicationJSON,
			"X-Atlassian-Token": "no-check",
		},
	)
	return resp, err
}

// MoveAttachmentsFromPage fromPageIDに添付されているファイルをdstPageIDに移す
func (t *Client) MoveAttachmentsFromPage(fromPageID, dstPageID string) ([]*http.Response, error) {
	metadata, err := t.FetchAttachmentMetaData(fromPageID)
	if err != nil {
		return nil, err
	}
	var ret []*http.Response
	for _, v := range metadata.Results {
		resp, err := t.MoveAttachment(fromPageID, v.ID, dstPageID)
		if err != nil {
			return nil, err
		}
		ret = append(ret, resp)
	}
	return ret, nil
}

// FetchAttachmentMetaData pageIDに添付されたファイルのデータを取得する
func (t *Client) FetchAttachmentMetaData(pageID string) (*AttachmentResults, error) {
	targetURL := t.baseURL + "/rest/api/content/" + pageID + "/child/attachment"
	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		targetURL,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	bin, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var res AttachmentResults
	err = json.Unmarshal(bin, &res)

	return &res, err
}

// DownloadAttachmentsFromPage ページに添付してあるファイルをダウンロードする
func (t *Client) DownloadAttachmentsFromPage(pageID, directory string) error {
	res, err := t.FetchAttachmentMetaData(pageID)
	if err != nil {
		return err
	}

	os.MkdirAll(directory, os.ModePerm)
	for _, v := range res.Results {
		downloadURL := t.baseURL + v.Links.Download
		path, err := fileio.GetNonExistFileName(filepath.Join(directory, v.Title), 1000)
		if err != nil {
			return err
		}
		t.DownloadFromURL(downloadURL, path)
	}
	return nil
}

// DownloadFromURL ダウンロードする
func (t *Client) DownloadFromURL(url, outputFilepath string) error {
	resp, err := t.httpClient.DoRequest(
		http.MethodGet,
		url,
		nil,
		nil,
	)
	if err != nil {
		return err
	}
	fh, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer fh.Close()
	defer resp.Body.Close()
	_, err = io.Copy(fh, resp.Body)

	return err
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
