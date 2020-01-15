package f3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Client for the Form3 API.
type Client struct {
	baseURL   *url.URL
	userAgent string

	accessKey       string
	secretAccessKey string

	httpClient *http.Client
}

// NewClient returns a new Form3 Client.
func NewClient(opts ...func(*Client)) *Client {
	c := &Client{
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.form3.tech",
		},
		userAgent:  "client-go/0.1.0",
		httpClient: &http.Client{},
	}

	for _, f := range opts {
		f(c)
	}

	return c
}

// WithBaseURL sets the Client base URL.
func WithBaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// WithAccessKey sets the Client access key for authentication.
func WithAccessKey(key string) func(*Client) {
	return func(c *Client) {
		c.accessKey = key
	}
}

// WithSecretAccessKey sets the Client secret access key for authentication.
func WithSecretAccessKey(key string) func(*Client) {
	return func(c *Client) {
		c.secretAccessKey = key
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Serialized-Access-Key", c.accessKey)
	req.Header.Set("Serialized-Secret-Access-Key", c.secretAccessKey)

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil
		}
	}

	return resp, err
}
