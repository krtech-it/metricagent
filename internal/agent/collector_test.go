package agent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	assert.Equal(t, int64(0), c.counter)
	assert.Equal(t, 0, len(c.Storage))
}

func TestCollector_Add(t *testing.T) {
	c := NewCollector()

	ctx, cancel := context.WithCancel(context.Background())
	go c.Add(ctx, 20)
	time.Sleep(22 * time.Millisecond)
	randomValue, _ := c.Storage["RandomValue"]
	assert.Equal(t, len(gaugeArea)+1, len(c.Storage))
	assert.Equal(t, int64(1), c.counter)
	time.Sleep(22 * time.Millisecond)
	cancel()
	assert.Equal(t, int64(2), c.counter)
	assert.NotEqual(t, randomValue, c.Storage["RandomValue"])
}
