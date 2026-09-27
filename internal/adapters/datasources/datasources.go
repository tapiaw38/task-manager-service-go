package datasources

import (
	firestoreclient "github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/firestore"
	"github.com/tapiaw38/task-manager-service-go/internal/adapters/datasources/jsonstore"
)

type Datasources struct {
	TaskStore       *jsonstore.Store
	FirestoreClient *firestoreclient.Client
}

func CreateDatasources(store *jsonstore.Store, firestore *firestoreclient.Client) *Datasources {
	return &Datasources{
		TaskStore:       store,
		FirestoreClient: firestore,
	}
}
