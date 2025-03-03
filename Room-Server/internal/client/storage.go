package client

import (
	"errors"
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

type ClientStorage struct {
	clients sync.Map
	mu      sync.RWMutex
}

var (
	ErrRoomExists = errors.New("[Log Error] nickname has already exists")
)

func NewClientStorage() *ClientStorage {
	return &ClientStorage{
		clients: sync.Map{},
	}
}

func (s *ClientStorage) Register(clientConn *ClientConn, nickname string) *Client {

	if nickname == "" {
		nickname = utils.GenerateRandomNickname()
	}

	client := &Client{
		ClientConn: clientConn, // 使用 ClientConn
		Nickname:   nickname,
		RoomID:     "",
	}

	s.clients.Store(clientConn.ClientID, client)
	return client
}

func (s *ClientStorage) Logout(clientID string) {
	s.clients.Delete(clientID)
}

func (s *ClientStorage) GetClient(clientID string) (*Client, bool) {
	value, ok := s.clients.Load(clientID)
	if !ok {
		return nil, false
	}
	return value.(*Client), true
}
