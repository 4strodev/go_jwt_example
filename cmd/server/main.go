package main

import (
	// Initalizing project variables
	_ "github.com/4strodev/jwt/internal/services/init"
	"github.com/4strodev/jwt/internal/server"
	"log"
)

func main() {
	log.Fatal(server.App.Listen(":3000"))
}
