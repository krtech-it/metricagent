package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// HTTPObserver отправляет событие аудита POST-запросом на удалённый сервер.
type HTTPObserver struct {
	url    string
	client *http.Client
}

func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (o *HTTPObserver) Notify(event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	resp, err := o.client.Post(o.url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
