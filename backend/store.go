package main

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrNotFound = errors.New("conversation not found")

type ListFilter struct {
	Status   string
	Priority string
	Search   string
}

type UpdatePatch struct {
	Status   *Status
	Priority *Priority
}

type Store interface {
	List(f ListFilter) []Conversation
	Get(id string) (Conversation, error)
	Update(id string, p UpdatePatch) (Conversation, error)
}

type memoryStore struct {
	mu   sync.RWMutex
	data map[string]Conversation
}

func NewMemoryStore() *memoryStore {
	s := &memoryStore{data: make(map[string]Conversation)}
	for _, c := range seedConversations() {
		s.data[c.ID] = c
	}
	return s
}

func (s *memoryStore) List(f ListFilter) []Conversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	search := strings.ToLower(strings.TrimSpace(f.Search))
	out := make([]Conversation, 0, len(s.data))

	for _, c := range s.data {
		if f.Status != "" && string(c.Status) != f.Status {
			continue
		}
		if f.Priority != "" && string(c.Priority) != f.Priority {
			continue
		}
		if search != "" && !matchesSearch(c, search) {
			continue
		}
		out = append(out, c)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func matchesSearch(c Conversation, term string) bool {
	return strings.Contains(strings.ToLower(c.CustomerName), term) ||
		strings.Contains(strings.ToLower(c.CustomerEmail), term) ||
		strings.Contains(strings.ToLower(c.Subject), term)
}

func (s *memoryStore) Get(id string) (Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.data[id]
	if !ok {
		return Conversation{}, ErrNotFound
	}
	return c, nil
}

func (s *memoryStore) Update(id string, p UpdatePatch) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.data[id]
	if !ok {
		return Conversation{}, ErrNotFound
	}
	if p.Status != nil {
		c.Status = *p.Status
	}
	if p.Priority != nil {
		c.Priority = *p.Priority
	}
	s.data[id] = c
	return c, nil
}

func seedConversations() []Conversation {
	base := time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)
	mk := func(id, name, email, subject string, st Status, pr Priority, daysAgo int) Conversation {
		return Conversation{
			ID:            id,
			CustomerName:  name,
			CustomerEmail: email,
			Subject:       subject,
			Status:        st,
			Priority:      pr,
			CreatedAt:     base.AddDate(0, 0, -daysAgo),
		}
	}

	return []Conversation{
		mk("c1", "John Carter", "john.carter@example.com", "Cannot log into my account", StatusOpen, PriorityHigh, 1),
		mk("c2", "Sarah Kim", "sarah.kim@example.com", "Refund not received", StatusInProgress, PriorityHigh, 2),
		mk("c3", "Mike Johnson", "mike.johnson@acme.io", "Question about billing cycle", StatusOpen, PriorityLow, 3),
		mk("c4", "Emily Davis", "emily.davis@example.com", "Feature request: dark mode", StatusOpen, PriorityMedium, 4),
		mk("c5", "John Smith", "john.smith@globex.com", "App crashes on startup", StatusInProgress, PriorityHigh, 5),
		mk("c6", "Laura Chen", "laura.chen@example.com", "How to export my data?", StatusResolved, PriorityLow, 6),
		mk("c7", "David Brown", "david.brown@example.com", "Password reset email never arrives", StatusOpen, PriorityMedium, 7),
		mk("c8", "Anna White", "anna.white@initech.com", "Duplicate charge on my invoice", StatusResolved, PriorityHigh, 8),
		mk("c9", "Tom Wilson", "tom.wilson@example.com", "Cannot update profile picture", StatusInProgress, PriorityMedium, 9),
		mk("c10", "Rachel Green", "rachel.green@example.com", "Account deletion request", StatusOpen, PriorityLow, 10),
	}
}
