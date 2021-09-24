package confluence

import (
	"bytes"
	"net/http"
)

type Client struct {
	Client     *http.Client
	URL        string
	AuthBase64 string
}

func (c Client) Push(body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err

	}
	req.Header.Set("Authorization", c.AuthBase64)

	return c.Client.Do(req)
}
