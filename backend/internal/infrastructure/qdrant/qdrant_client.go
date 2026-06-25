package qdrant

import (
	"fmt"

	"jarvis/config"

	qdrantpb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewQdrantClient creates and returns a new Qdrant gRPC client.
func NewQdrantClient(cfg *config.QdrantConfig) (qdrantpb.QdrantClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Qdrant: %w", err)
	}

	client := qdrantpb.NewQdrantClient(conn)
	return client, nil
}
