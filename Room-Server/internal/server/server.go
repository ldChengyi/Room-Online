package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/CHENG/Room-Online/Room-Server/internal/client"
	"github.com/CHENG/Room-Online/Room-Server/internal/room"
	"github.com/google/uuid"
)

// TCPServer 结构体，表示一个 TCP 服务器
type TCPServer struct {
	port      int               // 服务器监听的端口
	roomMgr   *room.RoomManager // 房间管理器，负责管理房间的逻辑
	clientMgr *client.ClientManager
	listener  net.Listener // TCP 监听器
}

// NewTCPServer 创建一个新的 TCP 服务器实例
func NewTCPServer(port int, rm *room.RoomManager, cm *client.ClientManager) *TCPServer {
	return &TCPServer{
		port:      port, // 设置监听端口
		roomMgr:   rm,   // 绑定房间管理器
		clientMgr: cm,   // 绑定客户端管理器
	}
}

// Start 启动 TCP 服务器
func (s *TCPServer) Start(ctx context.Context) error {
	fmt.Println("[LOG INFO] TCPServer is starting...")

	// 创建一个监听配置，设置 KeepAlive 3 分钟，保持长连接
	lc := net.ListenConfig{
		KeepAlive: 3 * time.Minute,
	}

	// 监听指定端口
	listener, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err // 监听失败返回错误
	}

	s.listener = listener

	// 开启一个 goroutine 处理连接
	go s.acceptLoop(ctx, listener)

	fmt.Println("[LOG INFO] TCPServer has already started!")
	return nil // 成功启动服务器
}

// Stop 关闭 TCP 服务器，停止接受新的连接
func (s *TCPServer) Stop() error {
	if s.listener == nil {
		return errors.New("listener is not initialized")
	}
	fmt.Println("[LOG INFO] TCPServer has stopped!")
	return s.listener.Close()
}

// acceptLoop 持续接受新连接
func (s *TCPServer) acceptLoop(ctx context.Context, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		// 生成唯一的 ClientID
		clientID := uuid.New().String()

		// 创建 ClientConn 实例
		clientConn := &client.ClientConn{
			Conn:     conn,
			ClientID: clientID,
			// Nickname 可以在后续的握手或协议中设置
		}

		go s.handleConn(clientConn)
	}
}

func (s *TCPServer) handleConn(conn *client.ClientConn) {
	defer conn.Close()

	handler := NewTCPHandler(s.roomMgr, s.clientMgr)

	handler.HandleConnection(conn)
}
