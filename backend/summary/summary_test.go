package summary

import (
	"reflect"
	"testing"
)

func TestSummarizeConversations(t *testing.T) {
    tests := []struct{
		name string
		input []Conversation
		want map[string]int
	}{
		{
			name: "return all the status_priority combinations number",
			input: []Conversation{
				{ID: "1", Customer: "a", Priority: "HIGH", Status: "OPEN"},
				{ID: "2", Customer: "C", Priority: "MEDIUM", Status: "OPEN"},
				{ID: "10", Customer: "J", Priority: "HIGH", Status: "OPEN"},
				{ID: "7", Customer: "G", Priority: "LOW", Status: "RESOLVED"},
				{ID: "8", Customer: "H", Priority: "LOW", Status: "RESOLVED"},
			},	
			want: map[string]int{
    			"OPEN_HIGH":        2,
    			"RESOLVED_LOW": 	2,
    			"OPEN_MEDIUM":       1,
				},
		},
		{
			name: "skip when priority or status are empty",
			input:  []Conversation{
				{ID: "11", Customer: "K", Priority: "", Status: "OPEN"},
				{ID: "12", Customer: "L", Priority: "HIGH", Status: ""},
				{ID: "13", Customer: "M", Priority: "", Status: ""},
			},
			want: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			got := SummarizeConversations(tt.input)
			if !reflect.DeepEqual(got, tt.want){
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
