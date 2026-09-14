// NO-AI TASK
package summary

import "fmt"

type Conversation struct {
	ID       string
	Customer string
	Priority string
	Status   string
}

func SummarizeConversations(conversations []Conversation) map[string]int {
	result := make(map[string]int)
	if len(conversations) != 0 {
		for _, conv := range conversations {
			if conv.Status != "" && conv.Priority != "" {
				key := fmt.Sprintf("%s_%s", conv.Status, conv.Priority)
				result[key]++
			}
		}	
	}
	return result
}