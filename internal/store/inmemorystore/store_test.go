package inmemorystore_test

import (
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/store/inmemorystore"
	"github.com/stretchr/testify/assert"
)

func TestNew_WithDefaultStorage(t *testing.T) {
	store := inmemorystore.New()

	assert.NotNil(t, store)
	assert.NotNil(t, store.S)
}

func TestNew_WithCustomStorage(t *testing.T) {
	customStorage := &inmemorystore.MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	store := inmemorystore.New(inmemorystore.WithMemStorage(customStorage))

	assert.NotNil(t, store)
	assert.Equal(t, customStorage, store.S)
}
