package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func timeNow() time.Time {
	return time.Now()
}

func sessionTitle(msgs []message) string {
	for _, m := range msgs {
		if m.Content != "" {
			title := strings.TrimSpace(m.Content)
			if len(title) > 60 {
				title = title[:60] + "..."
			}
			return title
		}
	}
	return "Empty conversation"
}

func loadHistory(path string) []session {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var store struct {
		Sessions []session `json:"sessions"`
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return nil
	}
	if store.Sessions == nil {
		return nil
	}
	// newest first
	for i, j := 0, len(store.Sessions)-1; i < j; i, j = i+1, j-1 {
		store.Sessions[i], store.Sessions[j] = store.Sessions[j], store.Sessions[i]
	}
	return store.Sessions
}

func saveSession(path string, s session) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	existing := loadHistory(path)
	found := false
	for i, es := range existing {
		if es.ID == s.ID {
			existing[i] = s
			found = true
			break
		}
	}
	if !found {
		existing = append(existing, s)
	}
	if len(existing) > 50 {
		existing = existing[len(existing)-50:]
	}
	var store struct {
		Sessions []session `json:"sessions"`
	}
	store.Sessions = existing
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(path, data, 0644)
}

func deleteSession(path, id string) {
	sessions := loadHistory(path)
	var kept []session
	for _, s := range sessions {
		if s.ID != id {
			kept = append(kept, s)
		}
	}
	var store struct {
		Sessions []session `json:"sessions"`
	}
	store.Sessions = kept
	data, _ := json.MarshalIndent(store, "", "  ")
	os.WriteFile(path, data, 0644)
}
