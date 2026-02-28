module github.com/jalocompany/orbit

go 1.25.5

require (
	github.com/coder/websocket v1.8.12
	github.com/fundathyon/base v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
)

require (
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

replace github.com/fundathyon/base => ./base
