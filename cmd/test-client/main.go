package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type UserConfig struct {
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	TargetID string `json:"target_id"`
	Email    string `json:"email"`
}

type MessageRequest struct {
	Text     string `json:"text"`
	Type     string `json:"type"`
	TargetID string `json:"target_id"`
}

// Colores ANSI
const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run cmd/test-client/main.go <user_name>")
	}

	userName := os.Args[1]

	// 1. Leer archivo de configuración
	data, err := os.ReadFile("test_users.json")
	if err != nil {
		log.Fatalf("Error leyendo test_users.json: %v", err)
	}

	var configs map[string]UserConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		log.Fatalf("Error parseando JSON: %v", err)
	}

	user, ok := configs[userName]
	if !ok {
		log.Fatalf("Usuario '%s' no encontrado en test_users.json", userName)
	}

	// 2. Conectar al WebSocket de Orbit
	url := fmt.Sprintf("ws://localhost:8003/chat?access_token=%s", user.Token)
	fmt.Printf("%s[%s]%s Conectando a %s...\n", ColorYellow, user.Email, ColorReset, url)

	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Fatalf("Error en conexión WS: %v", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	// Mensaje de conexión premium con colores
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("%s[SESIÓN ACTIVA]%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sEmail:%s    %s\n", ColorCyan, ColorReset, user.Email)
	fmt.Printf("%sUser ID:%s  %s\n", ColorCyan, ColorReset, user.UserID)
	fmt.Printf("%sTarget:%s   %s\n", ColorCyan, ColorReset, user.TargetID)
	fmt.Println("----------------------------------------------------------------")

	// 3. Goroutine para leer mensajes entrantes
	go func() {
		for {
			var incoming any
			err := wsjson.Read(ctx, c, &incoming)
			if err != nil {
				log.Printf("Error leyendo del socket: %v", err)
				return
			}
			fmt.Printf("\n[MENSAJE RECIBIDO]: %+v\n> ", incoming)
		}
	}()

	// 4. Leer de terminal para enviar mensajes
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Escribe tu mensaje y presiona Enter (o 'exit' para salir):\n> ")
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "exit" {
			break
		}

		msg := MessageRequest{
			Text:     text,
			Type:     "chat",
			TargetID: user.TargetID,
		}

		err := wsjson.Write(ctx, c, msg)
		if err != nil {
			log.Printf("Error enviando mensaje: %v", err)
			break
		}
		fmt.Print("> ")
	}
}
