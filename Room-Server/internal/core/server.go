package core

import (
	"fmt"

	"github.com/CHENG/Room-Online/Room-Server/internal/client"
	"github.com/CHENG/Room-Online/Room-Server/internal/room"
	"github.com/CHENG/Room-Online/Room-Server/internal/server"
)

type Server struct {
	tcpSrv *server.TCPServer
}

func NewServer() *Server {
	roomStorage := room.NewRoomStorage()
	roomManager := room.NewRoomManager(roomStorage)

	clientStorage := client.NewClientStorage()
	clientManager := client.NewClientManager(clientStorage)

	tcpServer := server.NewTCPServer(1204, roomManager, clientManager)
	return &Server{
		tcpSrv: tcpServer,
	}
}

func (s *Server) Run() error {
	ctx := HandleSignals() // 信号处理上下文
	go s.tcpSrv.Start(ctx)
	fmt.Println("[SERVER INFO] Server has already started!")

	<-ctx.Done()
	fmt.Println("[Server Info] Server is stopping!")
	return s.tcpSrv.Stop()
}
