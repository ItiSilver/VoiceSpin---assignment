package main

import "testing"

func TestMemoryStore_List_Filters(t *testing.T) {
	store := NewMemoryStore()

	tests := []struct {
		name   string
		filter ListFilter
		check func(c Conversation) bool
		wantAtLeastOne bool
	}{
		{
			name:           "status filter returns only that status",
			filter:         ListFilter{Status: "OPEN"},
			check:          func(c Conversation) bool { return c.Status == StatusOpen },
			wantAtLeastOne: true,
		},
		{
			name:           "priority filter returns only that priority",
			filter:         ListFilter{Priority: "HIGH"},
			check:          func(c Conversation) bool { return c.Priority == PriorityHigh },
			wantAtLeastOne: true,
		},
		{
			name:           "status + priority filters combine with AND",
			filter:         ListFilter{Status: "OPEN", Priority: "LOW"},
			check:          func(c Conversation) bool { return c.Status == StatusOpen && c.Priority == PriorityLow },
			wantAtLeastOne: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := store.List(tc.filter)
			if tc.wantAtLeastOne && len(got) == 0 {
				t.Fatalf("expected at least one result, got none")
			}
			for _, c := range got {
				if !tc.check(c) {
					t.Errorf("conversation %s does not match filter %+v", c.ID, tc.filter)
				}
			}
		})
	}
}

func TestMemoryStore_List_Search(t *testing.T) {
	store := NewMemoryStore()
	
	got := store.List(ListFilter{Search: "JoHn"})

	gotIDs := map[string]bool{}
	for _, c := range got {
		gotIDs[c.ID] = true
	}

	for _, want := range []string{"c1", "c3", "c5"} {
		if !gotIDs[want] {
			t.Errorf("expected search %q to include %s; got %v", "JoHn", want, keys(gotIDs))
		}
	}
	if gotIDs["c2"] {
		t.Errorf("search %q should not match c2 (Sarah Kim / refund)", "JoHn")
	}
}

func TestMemoryStore_List_NoFilterReturnsAll(t *testing.T) {
	store := NewMemoryStore()
	if got := len(store.List(ListFilter{})); got != len(seedConversations()) {
		t.Errorf("expected %d conversations, got %d", len(seedConversations()), got)
	}
}

func TestMemoryStore_List_NewestFirst(t *testing.T) {
	store := NewMemoryStore()
	got := store.List(ListFilter{})
	for i := 1; i < len(got); i++ {
		if got[i-1].CreatedAt.Before(got[i].CreatedAt) {
			t.Errorf("results not sorted newest-first at index %d", i)
		}
	}
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
