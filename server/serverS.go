package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const buffer = 1024 // 1024 bytes

var (
	clientsL = make(map[net.Conn]string)
	mu       sync.Mutex
)

// Agrega un cliente a la lista global
func addClient(conn net.Conn, username string) {
	mu.Lock()
	defer mu.Unlock()
	clientsL[conn] = username
	log.Println(username, "se ha conectado al servidor.")
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		log.Println(clientsL[conn], "se ha desconectado.")
		delete(clientsL, conn)
		mu.Unlock()
		conn.Close()
	}()

	buf := make([]byte, buffer)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error al leer el mensaje:", err)
			return
		}

		message := string(buf[:n])
		fmt.Println(clientsL[conn], "dice:", message)

		// Reenviar mensaje a todos los clientes conectados
		broadcast(message, conn)
	}
}

func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clientsL {
		if client != sender {
			_, err := client.Write([]byte(clientsL[sender] + ": " + message))
			if err != nil {
				log.Println("Error enviando mensaje a", clientsL[client], ":", err)
			}
		}
	}
}

func StartServer(port *string) error {
	server, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		return err
	}
	defer server.Close()

	log.Println("Servidor escuchando en el puerto", *port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Error aceptando conexión:", err)
			continue
		}

		// Esperar el nombre del cliente
		buf := make([]byte, buffer)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error al recibir nombre de usuario:", err)
			conn.Close()
			continue
		}

		username := string(buf[:n])
		addClient(conn, username)

		// Manejar al cliente en una goroutine separada
		go handleClient(conn)
	}
}
