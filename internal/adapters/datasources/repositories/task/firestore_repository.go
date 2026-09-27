package task

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	firestoreclient "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/firestore"
	"github.com/tapiaw38/task-manager-service-go/internal/domain"
	apperrors "github.com/tapiaw38/task-manager-service-go/internal/platform/errors"
	"github.com/tapiaw38/task-manager-service-go/internal/platform/errors/mappings"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type firestoreRepository struct {
	client *firestoreclient.Client
	now    func() time.Time
	newID  func() string
}

type taskDocument struct {
	Title       string    `firestore:"title"`
	Description string    `firestore:"description"`
	Completed   bool      `firestore:"completed"`
	CreatedAt   time.Time `firestore:"createdAt"`
	UpdatedAt   time.Time `firestore:"updatedAt"`
}

func NewFirestoreRepository(client *firestoreclient.Client) Repository {
	return &firestoreRepository{
		client: client,
		now:    func() time.Time { return time.Now().UTC() },
		newID:  func() string { return uuid.NewString() },
	}
}

func (d taskDocument) toDomain(id string) domain.Task {
	return domain.Task{
		ID:          id,
		Title:       d.Title,
		Description: d.Description,
		Completed:   d.Completed,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func toDocument(task domain.Task) taskDocument {
	return taskDocument{
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func (r *firestoreRepository) Create(ctx context.Context, t domain.Task) (*domain.Task, apperrors.ApplicationError) {
	created := t
	created.ID = r.newID()
	created.CreatedAt = r.now()
	created.UpdatedAt = created.CreatedAt

	if _, err := r.client.Collection().Doc(created.ID).Set(ctx, toDocument(created)); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TaskCreateStoreError, err)
	}

	return &created, nil
}

func (r *firestoreRepository) List(ctx context.Context) ([]domain.Task, int, apperrors.ApplicationError) {
	documents := r.client.Collection().
		OrderBy("createdAt", firestore.Desc).
		Documents(ctx)
	defer documents.Stop()

	tasks := []domain.Task{}

	for {
		snapshot, err := documents.Next()
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, 0, apperrors.NewApplicationError(mappings.TaskListStoreError, err)
		}

		var document taskDocument
		if err := snapshot.DataTo(&document); err != nil {
			return nil, 0, apperrors.NewApplicationError(mappings.TaskListStoreError, err)
		}

		tasks = append(tasks, document.toDomain(snapshot.Ref.ID))
	}

	return tasks, len(tasks), nil
}

func (r *firestoreRepository) Get(ctx context.Context, id string) (*domain.Task, apperrors.ApplicationError) {
	snapshot, err := r.client.Collection().Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, notFoundError(id)
		}

		return nil, apperrors.NewApplicationError(mappings.TaskGetStoreError, err)
	}

	var document taskDocument
	if err := snapshot.DataTo(&document); err != nil {
		return nil, apperrors.NewApplicationError(mappings.TaskGetStoreError, err)
	}

	found := document.toDomain(snapshot.Ref.ID)

	return &found, nil
}

func (r *firestoreRepository) SetCompleted(ctx context.Context, id string, completed bool) (*domain.Task, apperrors.ApplicationError) {
	updatedAt := r.now()

	_, err := r.client.Collection().Doc(id).Update(ctx, []firestore.Update{
		{Path: "completed", Value: completed},
		{Path: "updatedAt", Value: updatedAt},
	})
	if err != nil {
		if isNotFound(err) {
			return nil, notFoundError(id)
		}

		return nil, apperrors.NewApplicationError(mappings.TaskCompleteStoreError, err)
	}

	return r.Get(ctx, id)
}

func (r *firestoreRepository) Delete(ctx context.Context, id string) apperrors.ApplicationError {
	if _, appErr := r.Get(ctx, id); appErr != nil {
		return appErr
	}

	if _, err := r.client.Collection().Doc(id).Delete(ctx); err != nil {
		return apperrors.NewApplicationError(mappings.TaskDeleteStoreError, err)
	}

	return nil
}
