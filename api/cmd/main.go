package main

import (
	"github.com/woojiahao/pino/api/pkg/database"
	"github.com/woojiahao/pino/api/pkg/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	port = 3000
)

var (
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbAddr = getHost()
	dbPort = getPort()
)

func getPort() uint16 {
	p := os.Getenv("DB_PORT")
	if p == "" {
		return 27017
	}
	port, err := strconv.Atoi(p)

	if err != nil {
		log.Fatalf("Invalid port")
	}

	return uint16(port)
}

func getHost() string {
	h := os.Getenv("DB_HOST")
	if h == "" {
		return "localhost"
	}

	return h
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Pong!"))
}

func main() {
	// Connect to MongoDB database
	db := database.New(dbUser, dbPass, dbAddr, dbPort)
	err := db.Connect()

	if err != nil {
		log.Fatalf("error occurred, %v", err)
	}

	// Listen to program stop signals to properly close the database resources
	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		_ = <-ch
		signal.Stop(ch)
		log.Println("Shutting down server")
		db.Close()
		os.Exit(0)
	}()

	s := server.New(port)
	s.AddEndpoint(server.GET, "/", ping)
	s.Start()
}
