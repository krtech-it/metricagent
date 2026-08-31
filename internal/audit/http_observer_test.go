package audit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPObserverPostsEvent(t *testing.T) {
	var gotMethod string
	var gotEvent Event

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &gotEvent))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	observer := NewHTTPObserver(server.URL)
	event := Event{Timestamp: 42, Metrics: []string{"Alloc", "Frees"}, IPAddress: "192.168.0.42"}
	require.NoError(t, observer.Notify(event))

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, event, gotEvent)
}

func TestHTTPObserverErrorOnUnreachable(t *testing.T) {
	observer := NewHTTPObserver("http://127.0.0.1:0")
	err := observer.Notify(Event{Timestamp: 1, Metrics: []string{"Alloc"}, IPAddress: "127.0.0.1"})
	assert.Error(t, err)
}
