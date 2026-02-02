package agent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	assert.Equal(t, int64(0), c.counter)
	assert.Equal(t, 0, len(c.storage))
}

func TestCollector_Add(t *testing.T) {
	c := NewCollector()

	ctx, cancel := context.WithCancel(context.Background())
	go c.Add(ctx, 20)
	time.Sleep(22 * time.Millisecond)
	randomValue, ok := c.storage["RandomValue"]
	require.Equal(t, true, ok)
	assert.Equal(t, len(gaugeArea)+1, len(c.storage))
	assert.Equal(t, int64(1), c.counter)
	time.Sleep(22 * time.Millisecond)
	cancel()
	assert.Equal(t, int64(2), c.counter)
	assert.NotEqual(t, randomValue, c.storage["RandomValue"])
}
