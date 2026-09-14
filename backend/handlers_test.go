package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(NewAPI(NewMemoryStore()).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestGetConversation_NotFound(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/api/conversations/does-not-exist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
}

func TestPatchConversation_PersistsChange(t *testing.T) {
	srv := newTestServer(t)
	patch := patch(t, srv.URL+"/api/conversations/c1", `{"status":"RESOLVED"}`)
	if patch.StatusCode != http.StatusOK {
		t.Fatalf("PATCH want 200, got %d", patch.StatusCode)
	}
	var updated Conversation
	decode(t, patch, &updated)
	if updated.Status != StatusResolved {
		t.Errorf("PATCH response status = %q, want RESOLVED", updated.Status)
	}
	if updated.Priority != PriorityHigh {
		t.Errorf("priority should be untouched, got %q", updated.Priority)
	}
	resp, err := http.Get(srv.URL + "/api/conversations/c1")
	if err != nil {
		t.Fatal(err)
	}
	var reFetched Conversation
	decode(t, resp, &reFetched)
	if reFetched.Status != StatusResolved {
		t.Errorf("after PATCH, GET status = %q, want RESOLVED", reFetched.Status)
	}
}

func TestPatchConversation_RejectsInvalidStatus(t *testing.T) {
	srv := newTestServer(t)

	resp := patch(t, srv.URL+"/api/conversations/c1", `{"status":"BANANA"}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 for invalid status, got %d", resp.StatusCode)
	}
}

func TestPatchConversation_RejectsEmptyBody(t *testing.T) {
	srv := newTestServer(t)

	resp := patch(t, srv.URL+"/api/conversations/c1", `{}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 for empty patch, got %d", resp.StatusCode)
	}
}

func TestListConversations_StatusFilter(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/api/conversations?status=RESOLVED")
	if err != nil {
		t.Fatal(err)
	}
	var got []Conversation
	decode(t, resp, &got)
	if len(got) == 0 {
		t.Fatal("expected some RESOLVED conversations")
	}
	for _, c := range got {
		if c.Status != StatusResolved {
			t.Errorf("got %s with status %q via ?status=RESOLVED", c.ID, c.Status)
		}
	}
}

func patch(t *testing.T, url, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response, into any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
