package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

//go:generate mockgen -source client.go -destination mock/http_requester.go -package mock HttpRequester
type HttpRequester interface {
	Get(ctx context.Context, endpoint string, header http.Header) (*http.Response, error)
	Post(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error)
	Put(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error)
	Patch(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error)
	Delete(ctx context.Context, endpoint string, header http.Header) (*http.Response, error)
	Close()
}

type Http struct {
	client  *http.Client
	baseUrl string
}

// NewHttpClient returns a new Http instance
func NewHttpClient(baseUrl string) *Http {
	return &Http{
		client: &http.Client{
			Timeout: time.Second * 2,
		},
		baseUrl: baseUrl,
	}
}

func (b *Http) Get(ctx context.Context, endpoint string, header http.Header) (*http.Response, error) {
	return b.makeRequest(ctx, endpoint, nil, header, http.MethodGet)
}

func (b *Http) Post(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error) {
	return b.makeRequest(ctx, endpoint, data, header, http.MethodPost)
}

func (b *Http) Put(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error) {
	return b.makeRequest(ctx, endpoint, data, header, http.MethodPut)
}
func (b *Http) Patch(ctx context.Context, endpoint string, data any, header http.Header) (*http.Response, error) {
	return b.makeRequest(ctx, endpoint, data, header, http.MethodPatch)
}

func (b *Http) Delete(ctx context.Context, endpoint string, header http.Header) (*http.Response, error) {
	return b.makeRequest(ctx, endpoint, nil, header, http.MethodDelete)
}

func (b *Http) makeRequest(ctx context.Context, endpoint string, data any, header http.Header, method string) (*http.Response, error) {
	buf := &bytes.Buffer{}
	if data != nil {
		if err := json.NewEncoder(buf).Encode(data); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, b.baseUrl+endpoint, buf)
	if err != nil {
		return nil, err
	}
	if header != nil {
		request.Header = header
	}
	request.Header.Set("Content-Type", "application/json")
	return b.client.Do(request)
}

func (b *Http) Close() {
	b.client.CloseIdleConnections()
}
