package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_RecordMetrics(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		want   want
	}{
		{
			name:   "Positive add metric to map",
			target: "/gauge/Alloc/33733",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/update"+tt.target, nil)
			//h.RecordMetrics(tt.args.w, tt.args.r)
			//создаем новый Recorder
			w := httptest.NewRecorder()
			//определяем хэндлер
			h := NewHandler()

			hf := http.HandlerFunc(h.RecordMetrics)
			hf.ServeHTTP(w, r)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
		})
	}
}
