package main

import (
	"fmt"

	"github.com/CHENG/Room-Online/Room-Server/internal/core"
	
)

func main() {
	srv := core.NewServer()
	fmt.Println("[SERVER INFO] Server is starting...")
	srv.Run()
}
