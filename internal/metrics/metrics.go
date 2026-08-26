package metrics

import "sync"

type Registry struct {
	mu     sync.RWMutex
	values map[string]int64
}

func New() *Registry { return &Registry{values: map[string]int64{}} }

func (r *Registry) Add(key string, delta int64) { r.mu.Lock(); r.values[key] += delta; r.mu.Unlock() }

func (r *Registry) Get(key string) int64 { return r.values[key] }

func (r *Registry) Snapshot() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make(map[string]int64, len(r.values))
	for k, v := range r.values {
		items[k] = v
	}
	return items
}
