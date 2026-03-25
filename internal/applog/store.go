package applog

import (
	"sync"
	"time"
)

const maxEntries = 200

type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Method     string    `json:"method,omitempty"`
	Path       string    `json:"path,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
	DurationMs int64     `json:"duration_ms,omitempty"`
	RemoteIP   string    `json:"remote_ip,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}

var (
	mu      sync.RWMutex
	entries []Entry
)

func Add(entry Entry) {
	mu.Lock()
	defer mu.Unlock()

	entries = append([]Entry{entry}, entries...)
	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}
}

func List() []Entry {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Entry, len(entries))
	copy(result, entries)
	return result
}
