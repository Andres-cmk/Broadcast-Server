package main

import (
	"broadcastServer/clients"
	"broadcastServer/server"
	"flag"
	"log"
)

func main() {

	operation := flag.String("operation", "", "start, connect")
	port := flag.String("port", "8080", "server port")
	username := flag.String("username", "You", "your username")

	flag.Parse()

	switch *operation {
	case "start":
		log.Println("Iniciando servidor en el puerto", *port)
		err := server.StartServer(port)
		if err != nil {
			log.Fatal("Error iniciando el servidor:", err)
		}
	case "connect":
		log.Println("Conectando al servidor como", *username)
		clients.ConnectToServer(port, username)
	default:
		log.Println("Uso: go run main.go -operation=start|connect -port=8080 [-username=TuNombre]")
	}

}
