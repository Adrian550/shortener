package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestShortener проверяет работу хендлеров для сокращения ссылок
func TestShortener(t *testing.T) {
	// Тестовые данные

	type want struct {
		postStatus int
		getStatus  int
		content    string
	}

	tests := []struct {
		name   string
		method string
		body   string
		path   string
		want   want
	}{
		{
			name:   "Success URL shortening",
			method: http.MethodPost,
			body:   "http://example.com",
			path:   "/",
			want: want{
				postStatus: http.StatusCreated,
				getStatus:  http.StatusTemporaryRedirect,
				content:    "http://example.com",
			},
		},
		{
			name:   "Empty POST request",
			method: http.MethodPost,
			body:   "",
			path:   "/",
			want: want{
				postStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "Non existent short URL",
			method: http.MethodGet,
			path:   "/nonexistent",
			want: want{
				getStatus: http.StatusBadRequest,
			},
		},
		{
			name:   "Success redirect after URL shortening",
			method: http.MethodPost,
			body:   "http://example.com",
			path:   "/",
			want: want{
				postStatus: http.StatusCreated,
				getStatus:  http.StatusTemporaryRedirect,
				content:    "http://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.method == http.MethodPost {
				postReq := httptest.NewRequest(http.MethodPost, tt.path, bytes.NewBufferString(tt.body))
				postReq.Host = "localhost"
				postResp := httptest.NewRecorder()

				handleGenURL(postResp, postReq)
				postResult := postResp.Result()

				assert.Equal(t, tt.want.postStatus, postResult.StatusCode)

				if postResult.StatusCode == http.StatusCreated {
					postBody := postResp.Body.String()
					parts := strings.Split(postBody, "/")
					shortKey := parts[len(parts)-1]

					getReq := httptest.NewRequest(http.MethodGet, "/"+shortKey, nil)
					getResp := httptest.NewRecorder()

					handleRedirect(getResp, getReq)
					getResult := getResp.Result()

					assert.Equal(t, tt.want.getStatus, getResult.StatusCode)

					location := getResult.Header.Get("Location")
					assert.Equal(t, tt.want.content, location)
				}
			} else if tt.method == http.MethodGet {
				getReq := httptest.NewRequest(http.MethodGet, tt.path, nil)
				getResp := httptest.NewRecorder()

				handleRedirect(getResp, getReq)
				getResult := getResp.Result()

				assert.Equal(t, tt.want.getStatus, getResult.StatusCode)
			}
		})
	}
}
