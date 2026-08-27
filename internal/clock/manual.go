package clock

import (
	"sync"
	"time"
)

type Manual struct {
	mu    sync.RWMutex
	value time.Time
}

func NewManual(value time.Time) *Manual { return &Manual{value: value.UTC()} }

func (m *Manual) Now() time.Time { m.mu.RLock(); defer m.mu.RUnlock(); return m.value }

func (m *Manual) Advance(delta time.Duration) {
	m.mu.Lock()
	m.value = m.value.Add(delta)
	m.mu.Unlock()
}

func (m *Manual) Set(value time.Time) { m.mu.Lock(); m.value = value.UTC(); m.mu.Unlock() }
