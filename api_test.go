package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatRequestJSON(t *testing.T) {
	body := chatRequest{
		Model:     "test-model",
		Messages:  []message{{Role: "user", Content: "hi"}},
		Stream:    true,
		MaxTokens: 100,
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"model":"test-model"`) {
		t.Errorf("missing model field")
	}
	if !strings.Contains(string(data), `"stream":true`) {
		t.Errorf("missing stream field")
	}
}

func TestSendMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "" {
			t.Errorf("unexpected Authorization header")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	cmd := sendMessage(srv.URL, "", "test-model", []message{{Role: "user", Content: "hi"}})
	msg := cmd()
	if err, ok := msg.(errMsg); ok {
		t.Fatalf("unexpected error: %v", err.err)
	}
	done, ok := msg.(streamDoneMsg)
	if !ok {
		t.Fatalf("expected streamDoneMsg, got %T", msg)
	}
	if done.fullText != "Hello" {
		t.Errorf("expected Hello, got %q", done.fullText)
	}
}

func TestSendMessageWithAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk-test123" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer srv.Close()

	cmd := sendMessage(srv.URL, "sk-test123", "test-model", []message{{Role: "user", Content: "hi"}})
	msg := cmd()
	if err, ok := msg.(errMsg); ok {
		t.Fatalf("unexpected error: %v", err.err)
	}
	done, ok := msg.(streamDoneMsg)
	if !ok {
		t.Fatalf("expected streamDoneMsg, got %T", msg)
	}
	if done.fullText != "ok" {
		t.Errorf("expected ok, got %q", done.fullText)
	}
}

func TestSendMessageAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer srv.Close()

	cmd := sendMessage(srv.URL, "", "test-model", []message{{Role: "user", Content: "hi"}})
	msg := cmd()
	err, ok := msg.(errMsg)
	if !ok {
		t.Fatalf("expected errMsg, got %T", msg)
	}
	if !strings.Contains(err.err.Error(), "400") {
		t.Errorf("expected 400 error, got %v", err.err)
	}
}

func TestSendMessageRequestError(t *testing.T) {
	cmd := sendMessage("http://invalid.local:1", "", "test-model", []message{{Role: "user", Content: "hi"}})
	msg := cmd()
	if _, ok := msg.(errMsg); !ok {
		t.Fatalf("expected errMsg, got %T", msg)
	}
}
