# go-rabbit-example
go rabbit example

## Rabbit
`docker run -d --rm -p 15672:15672 -p 5672:5672 rabbitmq:3-management`

## Background worker (consumer)
go run receive.go

## Producer (web server)
go run send-server.go 
