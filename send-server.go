package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	http.HandleFunc("/favicon.ico", handleFavicon)
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {

}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request received, %s!", r.URL.Path[1:])

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := filepath.Clean(r.URL.EscapedPath())
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
