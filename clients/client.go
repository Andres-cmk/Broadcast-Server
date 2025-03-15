package clients

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const buffer = 1024 // 1024 bytes

type Client struct {
	name      string
	entryDate time.Time
	messages  []string
	history   map[string]time.Time
}

func ConnectToServer(portServer *string, username *string) {
	conn, err := net.Dial("tcp", "localhost:"+*portServer)
	if err != nil {
		log.Fatal("No se pudo conectar al servidor:", err)
	}
	defer conn.Close()

	// Enviar el nombre de usuario al servidor
	_, err = conn.Write([]byte(*username))
	if err != nil {
		log.Fatal("Error enviando el nombre de usuario:", err)
	}

	// Goroutine para recibir mensajes del servidor
	go func() {
		buf := make([]byte, buffer)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("Desconectado del servidor.")
				os.Exit(0)
			}
			fmt.Println(string(buf[:n]))
		}
	}()

	// Leer entrada del usuario y enviarla al servidor
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Escribe un mensaje:")
	for scanner.Scan() {
		message := scanner.Text()
		if _, err := conn.Write([]byte(message)); err != nil {
			log.Println("Error enviando mensaje:", err)
			break
		}
	}
}
