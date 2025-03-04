package client

import (
	"fmt"
	"net"
	"sync"

	"github.com/CHENG/Room-Online/Room-Server/internal/pkg/utils"
)

// 封装的一个Conn结构体，包含了唯一的ClientID
type ClientConn struct {
	net.Conn
	ClientID string
}

// Client 结构体，表示一个客户端
type Client struct {
	*ClientConn
	Nickname string
	RoomID   string
}

func (c *Client) String() string {
	if c.ClientConn == nil {
		return fmt.Sprintf("Client{Nickname: %s, RoomID: %s, ClientConn: nil}", c.Nickname, c.RoomID)
	}
	return fmt.Sprintf("Client{Nickname: %s, RoomID: %s, ClientID: %s, RemoteAddr: %s}",
		c.Nickname, c.RoomID, c.ClientID, c.RemoteAddr().String())
}

type ClientStorage struct {
	clients sync.Map
	mu      sync.RWMutex
}

func NewClientStorage() *ClientStorage {
	return &ClientStorage{
		clients: sync.Map{},
	}
}

func (s *ClientStorage) Register(conn *ClientConn, nickname string) *Client {

	if nickname == "" {
		nickname = utils.GenerateRandomNickname()
	}

	client := &Client{
		ClientConn: conn, // 使用 ClientConn
		Nickname:   nickname,
		RoomID:     "",
	}

	s.clients.Store(conn.ClientID, client)
	return client
}

func (s *ClientStorage) Logout(conn *ClientConn) bool {
	if _, ok := s.clients.Load(conn.ClientID); !ok {
		return false // key 不存在，删除失败
	}

	// key 存在，删除并返回 true
	s.clients.Delete(conn.ClientID)
	return true
}

func (s *ClientStorage) ExitRoom(conn *ClientConn) {
	if client, ok := s.clients.Load(conn.ClientID); ok {
		client.(*Client).RoomID = ""
	}
}

func (s *ClientStorage) GetClient(conn *ClientConn) (*Client, bool) {
	value, ok := s.clients.Load(conn.ClientID)
	if !ok {
		return nil, false
	}
	return value.(*Client), true
}
