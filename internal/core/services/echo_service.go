package services

import "github.com/jalocompany/orbit/internal/core/ports"

type echoService struct{}

func NewEchoService() ports.EchoService {
	return &echoService{}
}

func (s *echoService) ProcessMessage(msg any) any {
	// Simple echo logic: return the received message
	return msg
}
