package core

import (
	"fmt"

	"github.com/CHENG/Room-Online/Room-server/internal/network/tcp"
	"github.com/CHENG/Room-Online/Room-server/internal/room"
)

type Server struct {
	tcpSrv *tcp.TCPServer
}

func NewServer() *Server {
	roomStorage := room.NewMemoryStorage()
	roomManager := room.NewRoomManager(roomStorage)
	tcp := tcp.NewTCPServer(1204, roomManager)
	return &Server{
		tcpSrv: tcp,
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
