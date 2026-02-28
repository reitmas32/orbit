package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jalocompany/orbit/internal/adapters/inbound/middlewares"
	"github.com/jalocompany/orbit/internal/adapters/inbound/websocket/chat"
	"github.com/jalocompany/orbit/internal/adapters/inbound/websocket/echo"
	"github.com/jalocompany/orbit/internal/core/services"
	"github.com/jalocompany/orbit/internal/core/settings"
)

func main() {
	settings.LoadDotEnv()
	settings.LoadEnvs()

	// 1. Core Logic Setup (Hexagonal Service)
	echoService := services.NewEchoService()
	chatService := services.NewChatService()

	// 2. Adapter Setup (HTTP/WebSocket Handler)
	echoHandler := echo.NewHandler(echoService)
	chatHandler := chat.NewHandler(chatService)

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
