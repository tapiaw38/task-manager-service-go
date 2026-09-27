package task_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
	repository "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/repositories/task"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
)

func newStoreFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "bd.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

func newRepository(t *testing.T, path string) repository.Repository {
	t.Helper()

	store, err := jsonstore.New(path)
	require.NoError(t, err)

	return repository.NewRepository(store)
}

func newInput(title string) domain.Task {
	return domain.Task{
		Title:       title,
		Description: "Description of " + title,
	}
}

func readStoreFile(t *testing.T, path string) []domain.Task {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var tasks []domain.Task
	require.NoError(t, json.Unmarshal(content, &tasks))

	return tasks
}

func TestStoreLoading(t *testing.T) {
	t.Run("when the file does not exist it starts empty and is created on the first write", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nested", "bd.json")
		repo := newRepository(t, path)

		_, appErr := repo.Create(context.Background(), newInput("First task"))
		require.Nil(t, appErr)

		assert.Len(t, readStoreFile(t, path), 1)
	})

	t.Run("when the file is empty it starts empty", func(t *testing.T) {
		repo := newRepository(t, newStoreFile(t, ""))

		tasks, total, appErr := repo.List(context.Background())
		require.Nil(t, appErr)
		assert.Empty(t, tasks)
		assert.Zero(t, total)
	})

	t.Run("when the file contains invalid json it fails fast", func(t *testing.T) {
		_, err := jsonstore.New(newStoreFile(t, "{not json"))
		assert.Error(t, err)
	})
}

func TestRepositoryCreate(t *testing.T) {
	path := newStoreFile(t, "[]")
	repo := newRepository(t, path)

	created, appErr := repo.Create(context.Background(), newInput("First task"))
	require.Nil(t, appErr)

	assert.NotEmpty(t, created.ID)
	assert.False(t, created.Completed)
	assert.False(t, created.CreatedAt.IsZero())
	assert.Equal(t, created.CreatedAt, created.UpdatedAt)

	persisted := readStoreFile(t, path)
	require.Len(t, persisted, 1)
	assert.Equal(t, created.ID, persisted[0].ID)
}

func TestRepositoryList(t *testing.T) {
	repo := newRepository(t, newStoreFile(t, "[]"))

	for _, title := range []string{"First task", "Second task", "Third task"} {
		_, appErr := repo.Create(context.Background(), newInput(title))
		require.Nil(t, appErr)
	}

	tasks, total, appErr := repo.List(context.Background())
	require.Nil(t, appErr)

	assert.Equal(t, 3, total)
	require.Len(t, tasks, 3)

	for i := 1; i < len(tasks); i++ {
		assert.False(t, tasks[i].CreatedAt.After(tasks[i-1].CreatedAt))
	}
}

func TestRepositoryGet(t *testing.T) {
	repo := newRepository(t, newStoreFile(t, "[]"))

	created, appErr := repo.Create(context.Background(), newInput("First task"))
	require.Nil(t, appErr)

	t.Run("when the task exists", func(t *testing.T) {
		found, appErr := repo.Get(context.Background(), created.ID)
		require.Nil(t, appErr)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("when the task does not exist", func(t *testing.T) {
		_, appErr := repo.Get(context.Background(), "missing")
		require.NotNil(t, appErr)
		assert.True(t, appErr.IsBasedOn(mappings.TaskNotFoundError))
	})
}

func TestRepositorySetCompleted(t *testing.T) {
	path := newStoreFile(t, "[]")
	repo := newRepository(t, path)

	created, appErr := repo.Create(context.Background(), newInput("First task"))
	require.Nil(t, appErr)

	t.Run("when the task is marked as completed it keeps its identity", func(t *testing.T) {
		updated, appErr := repo.SetCompleted(context.Background(), created.ID, true)
		require.Nil(t, appErr)

		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, created.CreatedAt, updated.CreatedAt)
		assert.True(t, updated.Completed)
		assert.False(t, updated.UpdatedAt.Before(updated.CreatedAt))

		persisted := readStoreFile(t, path)
		require.Len(t, persisted, 1)
		assert.True(t, persisted[0].Completed)
	})

	t.Run("when the task is marked as pending again", func(t *testing.T) {
		updated, appErr := repo.SetCompleted(context.Background(), created.ID, false)
		require.Nil(t, appErr)
		assert.False(t, updated.Completed)
	})

	t.Run("when the task does not exist", func(t *testing.T) {
		_, appErr := repo.SetCompleted(context.Background(), "missing", true)
		require.NotNil(t, appErr)
		assert.True(t, appErr.IsBasedOn(mappings.TaskNotFoundError))
	})
}

func TestRepositoryDelete(t *testing.T) {
	path := newStoreFile(t, "[]")
	repo := newRepository(t, path)

	created, appErr := repo.Create(context.Background(), newInput("First task"))
	require.Nil(t, appErr)

	t.Run("when the task exists it is removed from the file", func(t *testing.T) {
		require.Nil(t, repo.Delete(context.Background(), created.ID))
		assert.Empty(t, readStoreFile(t, path))
	})

	t.Run("when the task does not exist", func(t *testing.T) {
		appErr := repo.Delete(context.Background(), "missing")
		require.NotNil(t, appErr)
		assert.True(t, appErr.IsBasedOn(mappings.TaskNotFoundError))
	})
}

func TestRepositoryConcurrentCreates(t *testing.T) {
	path := newStoreFile(t, "[]")
	repo := newRepository(t, path)

	const attempts = 25

	var wg sync.WaitGroup

	for i := range attempts {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if _, appErr := repo.Create(context.Background(), newInput(fmt.Sprintf("Task %d", i))); appErr != nil {
				t.Errorf("unexpected error: %s", appErr.Message())
			}
		}(i)
	}

	wg.Wait()

	persisted := readStoreFile(t, path)
	assert.Len(t, persisted, attempts)

	ids := make(map[string]struct{}, len(persisted))
	for _, task := range persisted {
		ids[task.ID] = struct{}{}
	}

	assert.Len(t, ids, attempts)
}
