package server

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/CHENG/Room-Online/Room-Server/internal/client"
	"github.com/CHENG/Room-Online/Room-Server/internal/room"
)

type TCPHandler struct {
	RoomManager   *room.RoomManager
	ClientManager *client.ClientManager
}

func NewTCPHandler(rm *room.RoomManager, cm *client.ClientManager) *TCPHandler {
	return &TCPHandler{RoomManager: rm, ClientManager: cm}
}

func (h *TCPHandler) HandleConnection(conn *client.ClientConn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		response := h.handleCommand(cmd, conn)
		conn.Write(append([]byte(response), '\n'))
	}
}

func (h *TCPHandler) handleCommand(cmd string, conn *client.ClientConn) string {
	parts := strings.Split(cmd, " ")
	if len(parts) < 1 {
		return "ERROR Invalid command"
	}

	switch strings.ToUpper(parts[0]) {
	case "REGISTER":
		if len(parts) < 2 {
			return "ERROR Usage: REGISTER <NickName>"
		}
		return h.handleRegister(parts[1], conn)
	case "EXITROOM":
		return h.handleExitRoom(conn)
	case "CHECK":
		if len(parts) < 2 || parts[1] != "CLIENT" {
			return "ERROR Invalid CHECK command"
		}
		return h.handleCheckClient(conn)
	default:
		return "ERROR Unknown command"
	}
}

func (h *TCPHandler) handleRegister(nickname string, conn *client.ClientConn) string {
	// 注册昵称的处理
	client := h.ClientManager.Register(conn, nickname)
	return fmt.Sprintf("OK your client Infomation is there: %s", client.String())
}

func (h *TCPHandler) handleExitRoom(conn *client.ClientConn) string {
	// 退出房间的处理
	h.ClientManager.ExitRoom(conn)
	return "OK You have exited the room"
}

func (h *TCPHandler) handleCheckClient(conn *client.ClientConn) string {
	// 检查客户端的状态
	client, ok := h.ClientManager.GetClient(conn)
	if !ok {
		return fmt.Sprintf("ERROR Client not found")
	}
	return fmt.Sprintf("OK Client found: %s in room: %s", client.Nickname, client.RoomID)
}
