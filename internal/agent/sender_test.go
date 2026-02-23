package agent

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
)

func TestSendMetric(t *testing.T) {
	mux := http.NewServeMux()
	var gotPaths string
	var host = "127.0.0.1:8080"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gotPaths = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	ln, err := net.Listen("tcp", host)
	require.NoError(t, err)
	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)
	defer srv.Close()

	tests := []struct {
		name  string
		value interface{}
		url   string
	}{
		{name: "Alloc", value: 45.34, url: "/update/gauge/Alloc/45.34"},
		{name: "PollCount", value: int64(1), url: "/update/counter/PollCount/1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMetric(tt.name, tt.value, host)
			require.NoError(t, err)
			assert.Equal(t, gotPaths, tt.url)
		})
	}
}
