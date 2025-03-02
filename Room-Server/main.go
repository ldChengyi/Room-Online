package main

import (
	"fmt"

	"github.com/CHENG/Room-Online/Room-server/internal/core"
)

func main() {
	srv := core.NewServer()
	fmt.Println("[SERVER INFO] Server is starting...")
	srv.Run()
}
