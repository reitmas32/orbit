package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jalocompany/orbit/internal/adapters/inbound/middlewares"
	"github.com/jalocompany/orbit/internal/adapters/inbound/websocket/chat"
	"github.com/jalocompany/orbit/internal/adapters/inbound/websocket/echo"

	"github.com/fundathyon/base/core/eventbus"
	"github.com/fundathyon/base/infrastructure/eventbus/rabbit"
	"github.com/jalocompany/orbit/internal/core/entities"
	chatSvc "github.com/jalocompany/orbit/internal/core/services/chat"
	echoSvc "github.com/jalocompany/orbit/internal/core/services/echo"
	"github.com/jalocompany/orbit/internal/core/services/message_store"
	"github.com/jalocompany/orbit/internal/core/services/multiplexer"
	notifierSvc "github.com/jalocompany/orbit/internal/core/services/notifier"

	"github.com/jalocompany/orbit/internal/core/settings"
)

func main() {
	settings.LoadDotEnv()
	settings.LoadEnvs()

	// 1. Core Infrastructure Setup (RabbitMQ)
	rabbitConfig := rabbit.DefaultConfig(settings.Settings.RABBITMQ_URL)
	rabbitConfig.ExchangeName = settings.Settings.RABBITMQ_EXCHANGE

	// Fábrica de eventos para reconstruir SessionMessage al consumir
	eventFactory := func(eventName string) eventbus.Event {
		if eventName == "session.message" {
			return &entities.SessionMessage{}
		}
		return nil
	}

	eventBus, err := rabbit.NewEventBus(rabbitConfig, eventFactory)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer eventBus.Close()

	// 2. Core Logic Setup (Hexagonal Service)
	echoService := echoSvc.NewEchoService()
	chatService := chatSvc.NewChatService(eventBus)
	notifierService := notifierSvc.UserNotifierFactory()
	storageService := message_store.MessageStoreFactory()
	multiplexerService := multiplexer.MultiplexerFactory()

	// 3. Adapter Setup (HTTP/WebSocket Handler)
	echoHandler := echo.NewHandler(echoService)
	chatHandler := chat.NewHandler(chatService, notifierService, storageService, multiplexerService, eventBus)

	// 3. Configure API Endpoints
	http.Handle(settings.Settings.ROOT_PATH+"/echo", echoHandler)

	http.HandleFunc(settings.Settings.ROOT_PATH+"/chat", middlewares.NewChain(
		middlewares.WebSocketLoggerMiddleware,
		middlewares.JWTAuthMiddleware,
	).Then(chatHandler.ServeHTTP))

	port := fmt.Sprintf(":%s", settings.Settings.PORT)
	log.Printf("Starting WebSocket server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
