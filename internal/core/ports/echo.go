package ports

type EchoService interface {
	ProcessMessage(msg any) any
}
