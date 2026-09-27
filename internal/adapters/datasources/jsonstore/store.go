package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tapiaw38/task-manager-service-go/internal/domain"
)

type Store struct {
	path string
	mu   sync.RWMutex
	data []domain.Task
}

func New(path string) (*Store, error) {
	store := &Store{path: path, data: []domain.Task{}}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) load() error {
	content, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("reading store file %q: %w", s.path, err)
	}

	if len(content) == 0 {
		return nil
	}

	var tasks []domain.Task
	if err := json.Unmarshal(content, &tasks); err != nil {
		return fmt.Errorf("decoding store file %q: %w", s.path, err)
	}

	s.data = tasks

	return nil
}

func (s *Store) Snapshot() []domain.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make([]domain.Task, len(s.data))
	copy(snapshot, s.data)

	return snapshot
}

func (s *Store) Mutate(fn func(data []domain.Task) ([]domain.Task, error)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	working := make([]domain.Task, len(s.data))
	copy(working, s.data)

	updated, err := fn(working)
	if err != nil {
		return err
	}

	if err := s.persist(updated); err != nil {
		return err
	}

	s.data = updated

	return nil
}

func (s *Store) persist(data []domain.Task) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding store file %q: %w", s.path, err)
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating store directory %q: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, filepath.Base(s.path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temporary store file: %w", err)
	}
	tmpName := tmp.Name()

	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return fmt.Errorf("writing temporary store file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temporary store file: %w", err)
	}

	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("replacing store file %q: %w", s.path, err)
	}

	return nil
}
