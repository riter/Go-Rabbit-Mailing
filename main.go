package main

import "log"

type ResponseMailing struct {
	Status int
	Message string
	Action string
}

type User struct {
	FullName string
	Email string
}

type RequestMailing struct {
	Users  []User
	Filename string
	Action string
}

const exchange = "riter-mailing"
const queueRequest = "mailing-request"
const queueResponse = "mailing-response"

const urlServe = "amqp://guest:guest@localhost:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	receiveMailing()
}