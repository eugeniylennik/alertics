package handlers

import (
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/stretchr/testify/assert"
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
		target metrics.Data
		want   want
	}{
		{
			name: "Positive add metric to map",
			target: metrics.Data{
				Name:  "Alloc",
				Type:  "gauge",
				Value: 33812.12,
			},
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "Negative test data metrics != 3",
			target: metrics.Data{
				Name: "Alloc",
				Type: "gauge",
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/update/%s/%s/%.2f", tt.target.Type, tt.target.Name, tt.target.Value),
				nil,
			)
			//h.RecordMetrics(tt.args.w, tt.args.r)
			//создаем новый Recorder
			w := httptest.NewRecorder()
			//определяем хэндлер
			h := NewStorage()

			hf := http.HandlerFunc(h.RecordMetrics)
			hf.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tt.target.Value, h.Gauge[tt.target.Name])
		})
	}
}
