package f3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Client for the Form3 API.
type Client struct {
	baseURL        *url.URL
	servicePathURL string
	pagination     string
	userAgent      string

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

// WithServicePath sets the service path in the url.
func WithServicePath(serviceUrlPath string) func(*Client) {
	return func(c *Client) {
		c.servicePathURL = serviceUrlPath
	}
}

// WithPagination page number and page size used in Pagination
func WithPagination(pageNumber int, pageSize int) func(*Client) {
	return func(c *Client) {
		c.pagination = "?page[number]=" + strconv.Itoa(pageNumber) + "&page[size]=" + strconv.Itoa(pageSize)
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
