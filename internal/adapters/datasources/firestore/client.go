package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type Client struct {
	client     *firestore.Client
	collection string
}

func New(ctx context.Context, projectID, collection string) (*Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client for project %q: %w", projectID, err)
	}

	return &Client{client: client, collection: collection}, nil
}

func (c *Client) Collection() *firestore.CollectionRef {
	return c.client.Collection(c.collection)
}

func (c *Client) Close() error {
	return c.client.Close()
}
