package chat

import "github.com/jalocompany/orbit/internal/core/ports"

type chatService struct{}

func NewChatService() ports.ChatService {
	return &chatService{}
}

func (s *chatService) ProcessMessage(msg any) any {
	// Simple echo logic for now: return the received message
	return msg
}
