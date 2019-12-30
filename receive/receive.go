package main

import (
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s",msg,err)
	}
}

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	failOnError(err, "Gagal terkoneksi dengan RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err,"Gagal membuka Channel")
	defer ch.Close()

	// "message-from-php" : Merupakan key dari pesan yang akan diterima dari mana 
    // Key pada Receiver harus sama dan sesuai dengan Key yang diberikan oleh Sender
    // Notes: Sender boleh bersalah dari program dengan bahasa apa pun, tidak mesti harus satu bahasa
	// PENTING! Sesuaikan Key Sender dengan Receiver agar pesan diterima oleh Receiver
	
	q, err := ch.QueueDeclare(
		"message-from-golang", //name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Gagal mendeklarasikan Queue")

	msgs, err := ch.Consume(
		q.Name, // Queue
		"",	// Consumer
		true, // Auto-ack
		false, // Exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	failOnError(err, "Gagal mendaftarkan consumer")

	forever := make (chan bool)
	go func(){
		for d := range msgs{
			log.Printf("Received a Message From Golang: %s", d.Body)
		}
	}()

	log.Printf("[*] waiting for message. To exit press CTRL+C")
	<-forever
}