package applog

import (
	"sync"
	"time"
)

const maxEntries = 300

type Entry struct {
	Timestamp     time.Time `json:"timestamp"`
	Level         string    `json:"level"`
	Message       string    `json:"message"`
	Method        string    `json:"method,omitempty"`
	Path          string    `json:"path,omitempty"`
	Query         string    `json:"query,omitempty"`
	FullPath      string    `json:"full_path,omitempty"`
	StatusCode    int       `json:"status_code,omitempty"`
	DurationMs    int64     `json:"duration_ms,omitempty"`
	RemoteIP      string    `json:"remote_ip,omitempty"`
	UserAgent     string    `json:"user_agent,omitempty"`
	Host          string    `json:"host,omitempty"`
	Referer       string    `json:"referer,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	UserNIM       string    `json:"user_nim,omitempty"`
	UserRole      string    `json:"user_role,omitempty"`
	ResponseBytes int       `json:"response_bytes,omitempty"`
	ContentLength int64     `json:"content_length,omitempty"`
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
