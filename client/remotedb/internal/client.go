package internal

import (
	"context"

	"cloud.google.com/go/firestore"
)

type RemoteDBClient struct {
	client *firestore.Client
}

func (rdbc *RemoteDBClient) Client(
	ctx context.Context,
) (*firestore.Client, error) {
	return rdbc.client, nil
	// panic("unimplemented")
}
