package chat

import (
	"github.com/fundathyon/base/core/customctx"
	"github.com/fundathyon/base/core/eventbus"
	"github.com/jalocompany/orbit/internal/core/entities"
	"github.com/jalocompany/orbit/internal/core/ports"
)

type chatService struct {
	eb eventbus.EventBus
}

func NewChatService(eb eventbus.EventBus) ports.ChatService {
	return &chatService{
		eb: eb,
	}
}

func (s *chatService) PublishMessage(ctx *customctx.CustomContext, msg entities.SessionMessage) error {
	ctx.Logger().Info("[Chat Service] Publicando mensaje hacia el destinatario %s vía RabbitMQ", msg.RecipientID)

	// Publicar el evento usando el EventBus de base
	if err := s.eb.Publish(ctx, msg); err != nil {
		ctx.Logger().Error("[Chat Service] Error al publicar en RabbitMQ: %v", err)
		return err
	}

	return nil
}
