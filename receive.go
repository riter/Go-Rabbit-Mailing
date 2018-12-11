package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const directoryPath = "./public"

func receiveMailing() {
	conn, err := amqp.Dial(urlServe)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueRequest, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	e := ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	failOnError(e, "Failed to declare a exchange")

	ec := ch.QueueBind(
		queueRequest,
		queueRequest,
		exchange,
		false,
		nil,
		)

	failOnError(ec, "Failed to declare a exchange")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {

			var requestMails RequestMailing
			errDecode := json.Unmarshal(d.Body, &requestMails)

			failOnError(errDecode, "Failed decode message")

			switch requestMails.Action {
			case "post":
				processMessageMailCreateFile(requestMails)
				break
			case "delete":
				processMessageMailDeleteFile(requestMails)
				break
			}

			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}

func processMessageMailCreateFile(data RequestMailing) {

	for index := range data.Users {
		user := data.Users[index]

		log.Printf("mail sended to " + user.FullName + " with Email : " + user.Email + " and body: " + "File " + data.Filename + " created")
	}

	response, _ := json.Marshal(ResponseMailing{Status: 200, Message:"Mails sended", Action:"create"})
	go sendResponseStorage(response)
}

func processMessageMailDeleteFile(data RequestMailing) {
	for index := range data.Users {
		user := data.Users[index]

		log.Printf("mail sended to " + user.FullName + " with Email : " + user.Email + " and body: " + "File " + data.Filename + " deleted")
	}

	response, _ := json.Marshal(ResponseMailing{Status: 200, Message:"Mails sended", Action:"delete"})
	go sendResponseStorage(response)
}