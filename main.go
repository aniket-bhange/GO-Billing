package main

import (
	database "billing-gorilla/db"
	"billing-gorilla/model"
	"bytes"
	"log"
	"net/http"
	"time"

	amq "billing-gorilla/amq"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file present")
	}
}

func main() {

	queConfig := amq.ConnectRBMQ()
	msgs, err := amq.ConsumeMessage(queConfig)
	if err != nil {
		log.Printf("Error while consuming message")
	}

	forever := make(chan bool)
	db := database.ConnectDB()
	model.LoadModels(db.Db)
	router := NewRouter()

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router))
	// http.ListenAndServe(":8080", (router))

	<-forever

}
