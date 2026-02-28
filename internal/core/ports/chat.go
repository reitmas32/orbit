package ports

type ChatService interface {
	ProcessMessage(msg any) any
}
