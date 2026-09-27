package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
)

func TestTaskSanitize(t *testing.T) {
	tests := map[string]struct {
		task     domain.Task
		expected domain.Task
	}{
		"when the fields have surrounding whitespace": {
			task:     domain.Task{Title: "  Buy milk  ", Description: "  Go shopping  "},
			expected: domain.Task{Title: "Buy milk", Description: "Go shopping"},
		},
		"when the fields have repeated inner whitespace": {
			task:     domain.Task{Title: "Buy   milk", Description: "Go\tto  the   shop"},
			expected: domain.Task{Title: "Buy milk", Description: "Go to the shop"},
		},
		"when the fields only contain whitespace": {
			task:     domain.Task{Title: "   ", Description: "\t\n"},
			expected: domain.Task{Title: "", Description: ""},
		},
		"when the fields are already clean": {
			task:     domain.Task{Title: "Buy milk", Description: "Go shopping"},
			expected: domain.Task{Title: "Buy milk", Description: "Go shopping"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			task := tt.task
			task.Sanitize()

			assert.Equal(t, tt.expected.Title, task.Title)
			assert.Equal(t, tt.expected.Description, task.Description)
		})
	}
}

func TestNormalizeTitle(t *testing.T) {
	tests := map[string]struct {
		title    string
		expected string
	}{
		"when the title has mixed case":            {title: "Buy MILK", expected: "buy milk"},
		"when the title has surrounding spaces":    {title: "  Buy milk  ", expected: "buy milk"},
		"when the title has repeated inner spaces": {title: "Buy   MILK", expected: "buy milk"},
		"when the title is empty":                  {title: "   ", expected: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.NormalizeTitle(tt.title))
		})
	}
}
