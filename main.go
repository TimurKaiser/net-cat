package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type client struct {
	conn   net.Conn
	name   string
	writer *bufio.Writer
}

var (
	clients    []*client
	clientsMux sync.Mutex
	messages   []string
	historyLen = 10
)

func main() {
	args := os.Args[1:]
	port := 8989

	if len(args) == 0 {
		CreateTCPServer(port)
		return
	} else if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		CreateTCPServer(port)
		return
	}
}

func CreateTCPServer(value int) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: value})
	if err != nil {
		log.Fatal("Error creating the server:", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port : %d\n", value)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting a connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func printFileContents(writer *bufio.Writer, filename string) error {
	// Ouvrir le fichier texte
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Créer un scanner pour lire le fichier ligne par ligne
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := writer.WriteString(scanner.Text() + "\n") // Écrire chaque ligne dans le writer fourni
		if err != nil {
			return err
		}
	}

	// Vérifier les erreurs lors de la lecture du fichier
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func handleConnection(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Afficher le message de bienvenue
	writer.WriteString("Welcome to TCP-Chat!\n")

	// Afficher le logo
	err := printFileContents(writer, "logo.txt")
	if err != nil {
		log.Printf("Error displaying logo: %v\n", err)
		conn.Close()
		return
	}

	// Demander le nom de l'utilisateur
	writer.WriteString("[ENTER YOUR NAME]:")
	writer.Flush()

	// Lire le nom de l'utilisateur depuis la connexion
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading name: %v\n", err)
		conn.Close()
		return
	}

	name = strings.TrimSpace(name)

	newClient := &client{
		conn:   conn,
		name:   name,
		writer: writer,
	}

	clientsMux.Lock()
	clients = append(clients, newClient)
	clientsMux.Unlock()

	// Envoyer un message pour annoncer l'arrivée de l'utilisateur à tous les clients existants
	message := fmt.Sprintf("[%s] %s joined the chat!\n", time.Now().Format("2006-01-02 15:04:05"), name)
	clientsMux.Lock()
	for _, c := range clients {
		c.writer.WriteString(message)
		c.writer.Flush()
	}
	clientsMux.Unlock()

	// Envoyer l'historique des messages au nouveau client
	sendHistory(newClient)

	// Gérer les messages entrants de l'utilisateur
	for {
		// Effacer la ligne précédente dans le terminal du client
		writer.WriteString("\033[1A\033[K")

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n%s Has left the chat\n", name)
			break
		}

		message = fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, message)
		message = strings.TrimSpace(message)

		fmt.Print(message + "\n")

		clientsMux.Lock()
		for _, c := range clients {
			c.writer.WriteString(message + "\n")
			c.writer.Flush()
		}
		clientsMux.Unlock()

		// Ajouter le message à l'historique
		addMessageToHistory(message)
	}

	removeClient(newClient)
}

func sendHistory(c *client) {
	// Envoyer les messages dans l'ordre correcte
	for _, msg := range messages {
		c.writer.WriteString(msg + "\n")
		c.writer.Flush()
	}
}

func addMessageToHistory(message string) {
	// Verrouiller pour éviter les accès concurrents à l'historique des messages
	clientsMux.Lock()
	defer clientsMux.Unlock()

	// Ajouter le message à l'historique
	messages = append(messages, message)

	// Limiter la taille de l'historique
	if len(messages) > historyLen {
		messages = messages[len(messages)-historyLen:]
	}
}

func removeClient(clientToRemove *client) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	for i, c := range clients {
		if c == clientToRemove {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}
