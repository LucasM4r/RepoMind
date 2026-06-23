package rpc

import (
	"context"

	"github.com/LucasM4r/repomind/internal/pb"
	"google.golang.org/grpc"
)

type AIClient struct {
	EmbeddingClient pb.EmbeddingServiceClient
	LLMClient       pb.LLMServiceClient
}

func NewAIClient(conn *grpc.ClientConn) *AIClient {
	return &AIClient{
		EmbeddingClient: pb.NewEmbeddingServiceClient(conn),
		LLMClient:       pb.NewLLMServiceClient(conn),
	}
}

func (c *AIClient) GetEmbeddings(ctx context.Context, texts []string) (*pb.EmbeddingResponse, error) {
	return c.EmbeddingClient.GenerateEmbeddings(ctx, &pb.EmbeddingRequest{Texts: texts})
}
