package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-shiori/shiori/internal/http/webcontext"
	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name       string
		ok         bool
		message    any
		errParams  map[string]string
		statusCode int
	}{
		{
			name:       "successful response",
			ok:         true,
			message:    "success",
			errParams:  nil,
			statusCode: http.StatusOK,
		},
		{
			name:       "error response",
			ok:         false,
			message:    "error occurred",
			errParams:  map[string]string{"field": "invalid"},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewResponse(tt.message, tt.statusCode)
			assert.Equal(t, tt.message, resp.GetData())
			assert.Equal(t, tt.statusCode, resp.statusCode)
		})
	}
}

func TestResponse_IsError(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		want     bool
	}{
		{
			name:     "successful response",
			response: NewResponse("success", http.StatusOK),
			want:     false,
		},
		{
			name:     "error response",
			response: NewResponse("error", http.StatusBadRequest),
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.response.IsError())
		})
	}
}

func TestResponse_GetMessage(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		want     any
	}{
		{
			name:     "string message",
			response: NewResponse("test message", http.StatusOK),
			want:     "test message",
		},
		{
			name:     "struct message",
			response: NewResponse(struct{ Data string }{Data: "test"}, http.StatusOK),
			want:     struct{ Data string }{Data: "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.response.GetData())
		})
	}
}

func TestResponse_Send(t *testing.T) {
	tests := []struct {
		name           string
		response       *Response
		expectedStatus int
		expectedBody   any
	}{
		{
			name:           "plain response",
			response:       NewResponse("success", http.StatusOK),
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "json response",
			response: NewResponse(map[string]any{
				"message": "success",
			}, http.StatusOK),
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "success",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := webcontext.NewWebContext(w, r)

			err := tt.response.SendJSON(ctx)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var responseBody any
			err = json.NewDecoder(w.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
