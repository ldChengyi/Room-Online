package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/CHENG/Room-Online/Room-server/internal/room"
)

// TCPServer 结构体，表示一个 TCP 服务器
type TCPServer struct {
	port     int               // 服务器监听的端口
	roomMgr  *room.RoomManager // 房间管理器，负责管理房间的逻辑
	listener net.Listener      // TCP 监听器
}

// NewTCPServer 创建一个新的 TCP 服务器实例
func NewTCPServer(port int, rm *room.RoomManager) *TCPServer {
	return &TCPServer{
		port:    port, // 设置监听端口
		roomMgr: rm,   // 绑定房间管理器
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
		// 接受新的 TCP 连接
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return // 如果监听器已关闭，退出循环
		}

		// 为每个连接启动一个新的 goroutine 进行处理
		go s.handleConn(conn) // 这里可以优化，使用连接池
	}
}

// handleConn 处理单个 TCP 连接
func (s *TCPServer) handleConn(conn net.Conn) {
	defer conn.Close() // 确保连接在处理完成后关闭

	// 创建一个 TCP 处理器，负责解析和执行客户端的命令
	handler := NewTCPHandler(s.roomMgr)

	// 处理连接中的数据
	handler.HandleConnection(conn)
}
