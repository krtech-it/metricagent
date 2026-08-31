package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileObserverAppendsLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	observer, err := NewFileObserver(path)
	require.NoError(t, err)

	require.NoError(t, observer.Notify(Event{Timestamp: 1, Metrics: []string{"Alloc"}, IPAddress: "1.2.3.4"}))
	require.NoError(t, observer.Notify(Event{Timestamp: 2, Metrics: []string{"PollCount"}, IPAddress: "1.2.3.4"}))

	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.Len(t, lines, 2)

	var first Event
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &first))
	assert.Equal(t, int64(1), first.Timestamp)
	assert.Equal(t, []string{"Alloc"}, first.Metrics)
	assert.Equal(t, "1.2.3.4", first.IPAddress)

	var second Event
	require.NoError(t, json.Unmarshal([]byte(lines[1]), &second))
	assert.Equal(t, []string{"PollCount"}, second.Metrics)
}

func TestFileObserverInvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "audit.log")
	_, err := NewFileObserver(path)
	assert.Error(t, err)
}
