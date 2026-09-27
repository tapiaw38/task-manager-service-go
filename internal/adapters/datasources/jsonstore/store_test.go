package jsonstore_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
)

func newStore(t *testing.T, content string) (*jsonstore.Store, string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "bd.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	store, err := jsonstore.New(path)
	require.NoError(t, err)

	return store, path
}

func readFile(t *testing.T, path string) []domain.Task {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var tasks []domain.Task
	require.NoError(t, json.Unmarshal(content, &tasks))

	return tasks
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		content       string
		writeFile     bool
		expectsError  bool
		expectedTasks int
	}{
		"when the file contains tasks": {
			content:       `[{"id":"task-1","title":"First task"}]`,
			writeFile:     true,
			expectedTasks: 1,
		},
		"when the file is empty": {
			content:       "",
			writeFile:     true,
			expectedTasks: 0,
		},
		"when the file does not exist": {
			writeFile:     false,
			expectedTasks: 0,
		},
		"when the file contains invalid json": {
			content:      "{not json",
			writeFile:    true,
			expectsError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "bd.json")

			if tt.writeFile {
				require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o600))
			}

			store, err := jsonstore.New(path)

			if tt.expectsError {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Len(t, store.Snapshot(), tt.expectedTasks)
		})
	}
}

func TestSnapshotReturnsACopy(t *testing.T) {
	store, _ := newStore(t, `[{"id":"task-1","title":"First task"}]`)

	snapshot := store.Snapshot()
	snapshot[0].Title = "Mutated outside the store"

	assert.Equal(t, "First task", store.Snapshot()[0].Title)
}

func TestMutate(t *testing.T) {
	t.Run("when the mutator succeeds it persists the new collection", func(t *testing.T) {
		store, path := newStore(t, "[]")

		err := store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
			return append(data, domain.Task{ID: "task-1", Title: "First task"}), nil
		})

		require.NoError(t, err)
		assert.Len(t, store.Snapshot(), 1)
		assert.Len(t, readFile(t, path), 1)
	})

	t.Run("when the mutator fails neither memory nor disk change", func(t *testing.T) {
		store, path := newStore(t, `[{"id":"task-1","title":"First task"}]`)
		failure := errors.New("rejected by the repository")

		err := store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
			return append(data, domain.Task{ID: "task-2"}), failure
		})

		assert.ErrorIs(t, err, failure)
		assert.Len(t, store.Snapshot(), 1)
		assert.Len(t, readFile(t, path), 1)
	})

	t.Run("when the target directory does not exist it is created", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nested", "deep", "bd.json")
		store, err := jsonstore.New(path)
		require.NoError(t, err)

		require.NoError(t, store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
			return append(data, domain.Task{ID: "task-1"}), nil
		}))

		assert.Len(t, readFile(t, path), 1)
	})

	t.Run("when the mutator receives a copy the store is untouched until it returns", func(t *testing.T) {
		store, _ := newStore(t, `[{"id":"task-1","title":"First task"}]`)

		err := store.Mutate(func(data []domain.Task) ([]domain.Task, error) {
			data[0].Title = "Mutated inside the mutator"

			return nil, errors.New("aborted")
		})

		assert.Error(t, err)
		assert.Equal(t, "First task", store.Snapshot()[0].Title)
	})
}
