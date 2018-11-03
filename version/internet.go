package version

import (
	"io"
	"net/http"
	"strings"
)

// InternetLike is an interface used to enable mocking http requests
type InternetLike interface {
	Get(url string) (*http.Response, error)
}

// real implementation

type realInternet struct{}

func (h *realInternet) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// mock implementation

type mockInternet struct {
	responses map[string]string
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error {
	return nil
}

func newMockInternet() *mockInternet {
	return &mockInternet{
		responses: map[string]string{},
	}
}

func (mock *mockInternet) addStringResponse(url string, response string) {
	mock.responses[url] = response
}

func (mock *mockInternet) Get(url string) (*http.Response, error) {
	stringResponse := mock.responses[url]

	response := &http.Response{
		Body: noopCloser{strings.NewReader(stringResponse)},
	}

	return response, nil
}
