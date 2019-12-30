package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s",msg,err)
	}
}

func GetMessage(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
        var tmpl = template.Must(template.New("inbox").ParseFiles("views/inbox.html"))

        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var from = r.FormValue("from")
		var message = r.Form.Get("message")
		var data = map[string]string{"from": from, "message": message}
		
		SendMessage(from,message)

        if err := tmpl.Execute(w, data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    http.Error(w, "", http.StatusBadRequest)
}


func SendMessage(from string, message string){
	// Open Connection
	koneksi, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err,"Tidak bisa terkoneksi")
	defer koneksi.Close()

	// Open Channel
	ch, err := koneksi.Channel()
	failOnError(err,"Tidak bisa membuka channel")
	defer ch.Close()
	
	// Pesan yang dikirimkan 
	// Notes: Pesan masih dideklare di hardcode
	q, err := ch.QueueDeclare(
		"message-from-golang",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err,"Gagal mendeklarasikan Queue")

	// messagenya:
	// pesan := "This Message Send Via GOLANG"
	pesan := "From: "+from+", Message: "+message 

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(pesan),
		},
	)
	failOnError(err, "gagal mengirim pesan")
}

func HandleIndex(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
        var tmpl = template.Must(template.New("form").ParseFiles("views/index.html"))
        var err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    http.Error(w, "", http.StatusBadRequest)
}

func HandleInbox(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
        var tmpl = template.Must(template.New("inbox").ParseFiles("views/inbox.html"))
        var err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    http.Error(w, "", http.StatusBadRequest)
}

func main()  {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/vendor/bootstrap"))))

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/inbox",HandleInbox)
	http.HandleFunc("/SendMessage", GetMessage)

	fmt.Println("Menjalankan Server Pada Port 8080...")
	http.ListenAndServe(":8080",nil)
}