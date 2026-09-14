package main

import "time"

type Status string
type Priority string

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusResolved   Status = "RESOLVED"
)

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

func (s Status) Valid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusResolved:
		return true
	default:
		return false
	}
}

func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

type Conversation struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customerName"`
	CustomerEmail string    `json:"customerEmail"`
	Subject       string    `json:"subject"`
	Status        Status    `json:"status"`
	Priority      Priority  `json:"priority"`
	CreatedAt     time.Time `json:"createdAt"`
}
