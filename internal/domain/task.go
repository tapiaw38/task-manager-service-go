package domain

import (
	"strings"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (t *Task) Sanitize() {
	t.Title = collapseSpaces(t.Title)
	t.Description = collapseSpaces(t.Description)
}

func collapseSpaces(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func NormalizeTitle(title string) string {
	return strings.ToLower(collapseSpaces(title))
}
