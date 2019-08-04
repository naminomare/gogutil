package confluence

// AttachmentResults Results
type AttachmentResults struct {
	Results []AttachmentFetchResult `json:"results"`
	Start   float64                 `json:"start"`
	Limit   float64                 `json:"limit"`
	Size    float64                 `json:"size"`
	Links   map[string]string       `json:"_links"`
}

// AttachmentFetchResult 添付ファイルのメタデータ
type AttachmentFetchResult struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Status     string               `json:"status"`
	Title      string               `json:"title"`
	MetaData   AttachmentMetaData   `json:"metadata"`
	Extensions AttachmentExtensions `json:"extensions"`
	Expandable AttachmentExpandable `json:"_expandable"`
	Links      AttachmentLinks      `json:"_links"`
}

// AttachmentMetaData メタデータ
type AttachmentMetaData struct {
	MediaType  string                 `json:"mediaType"`
	Labels     AttachmentLabels       `json:"labels"`
	Expandable map[string]interface{} `json:"_expandable"`
}

// AttachmentLabels ラベル
type AttachmentLabels struct {
	Results []interface{}     `json:"results"`
	Start   float64           `json:"start"`
	Limit   float64           `json:"limit"`
	Size    float64           `json:"size"`
	Links   map[string]string `json:"_links"`
}

// AttachmentExtensions Extensions
type AttachmentExtensions struct {
	MediaType string  `json:"mediaType"`
	FileSize  float64 `json:"fileSize"`
	Comment   string  `json:"comment"`
}

// AttachmentExpandable expandable
type AttachmentExpandable struct {
	Container    string `json:"container"`
	Operations   string `json:"operations"`
	Children     string `json:"children"`
	Restrictions string `json:"restrictions"`
	History      string `json:"history"`
	// Ancestors string `json:"ancestors"`
	// Body string `json:"body"`
	// Version string `json:"version"`
	Descendants string `json:"descendants"`
	Space       string `json:"space"`
}

// AttachmentLinks links
type AttachmentLinks struct {
	Self      string `json:"self"`
	Webui     string `json:"webui"`
	Download  string `json:"download"`
	Thumbnail string `json:"thumbnail"`
}
