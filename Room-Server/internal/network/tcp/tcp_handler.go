package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/CHENG/Room-Online/Room-server/internal/room"
)

type TCPHandler struct {
	RoomManager *room.RoomManager
}

func NewTCPHandler(rm *room.RoomManager) *TCPHandler {
	return &TCPHandler{RoomManager: rm}
}

func (h *TCPHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		response := h.handleCommand(cmd, conn.RemoteAddr())
		conn.Write(append([]byte(response), '\n'))
	}
}

func (h *TCPHandler) handleCommand(cmd string, addr net.Addr) string {
	parts := strings.Split(cmd, " ")
	if len(parts) < 1 {
		return "ERROR Invalid command"
	}

	switch strings.ToUpper(parts[0]) {
	case "CREATE":
		return h.handleCreate(parts[1:], addr)
	case "JOIN":
		return h.handleJoin(parts[1:], addr)
	case "LIST":
		return h.handleList()
	default:
		return "ERROR Unknown command"
	}
}

func (h *TCPHandler) handleCreate(args []string, addr net.Addr) string {
	if len(args) < 2 {
		return "ERROR Usage: CREATE <room_id> <password>"
	}

	roomID := args[0]
	password := args[1]

	_, err := h.RoomManager.CreateRoom(roomID, password, addr)
	if err != nil {
		return fmt.Sprintf("ERROR Create failed: %v", err)
	}

	return fmt.Sprintf("OK Room %s created", roomID)
}

func (h *TCPHandler) handleJoin(args []string, addr net.Addr) string {
	// 类似CREATE的实现
	if len(args) < 2 {
		return "ERROR Usage: JOIN <room_id> <password>"
	}

	roomID := args[0]
	password := args[1]
	clientID := args[2]

	err := h.RoomManager.JoinRoom(roomID, password, clientID, addr)
	if err != nil {
		return fmt.Sprintf("ERROR Join failed: %v", err)
	}

	return fmt.Sprintf("OK you have joind room: %s", roomID)
}

func (h *TCPHandler) handleList() string {
	// 返回房间列表
	return h.RoomManager.ListRoom()
}
