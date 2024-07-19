package stores

import (
	"sync"

	"github.com/rcbadiale/go-rate-limiter/pkg/status"
)

// Memory represents a memory store for rate limiter statuses.
type Memory struct {
	statuses map[string]*status.Status
	mu       sync.Mutex
}

// NewMemory returns a new memory store.
func NewMemory() *Memory {
	return &Memory{
		statuses: make(map[string]*status.Status),
	}
}

// Get returns the status of a key.
//
// If the key does not exist, it resets the status.
func (m *Memory) Get(key string) *status.Status {
	s, ok := m.statuses[key]
	if !ok {
		return m.Reset(key)
	}
	return s
}

// Increment increments the count of a key.
//
// If the key does not exist, it resets the status.
func (m *Memory) Increment(key string) *status.Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.statuses[key]
	if !ok {
		return m.Reset(key)
	}
	s.Count += 1
	return s
}

// Reset resets the status of a key.
//
// If the key does not exist, it creates a new status.
func (m *Memory) Reset(key string) *status.Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := status.NewStatus()
	m.statuses[key] = s
	return s
}
