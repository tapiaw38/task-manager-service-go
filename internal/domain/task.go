package domain

import (
	"strings"
	"time"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
