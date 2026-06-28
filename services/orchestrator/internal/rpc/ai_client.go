package rpc

import (
	"context"
	"fmt"

	"github.com/LucasM4r/repomind/internal/domain"
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
	resp, err := c.EmbeddingClient.GenerateEmbeddings(ctx, &pb.EmbeddingRequest{Texts: texts})
	if err != nil {
		return nil, fmt.Errorf("[AIClient] failed to generate embeddings: %w", err)
	}
	return resp, nil
}

func (c *AIClient) GenerateText(ctx context.Context, history []domain.ChatMessage) (domain.ChatMessage, error) {
	if len(history) == 0 {
		return domain.ChatMessage{}, fmt.Errorf("[AIClient] history is empty, cannot generate text")
	}
	var chatHistory []*pb.Message

	for _, msg := range history {
		chatHistory = append(chatHistory, &pb.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	response, err := c.LLMClient.GenerateCompletion(ctx, &pb.CompletionRequest{History: chatHistory})
	if err != nil {
		return domain.ChatMessage{}, fmt.Errorf("[AIClient] failed to generate completion: %w", err)
	}

	return domain.ChatMessage{
		Role:    "assistant",
		Content: response.GetText(),
	}, nil
}
