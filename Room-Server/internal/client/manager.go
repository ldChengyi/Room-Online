package client

import "errors"

type ClientManager struct {
	Storage *ClientStorage
}

var (
	ErrNickNameExists       = errors.New("[Log Error] nickname has already exists")
	ErrLogoutClientNotFount = errors.New("[Log Error] logout failed: Client not found")
)

func NewClientManager(storage *ClientStorage) *ClientManager {
	return &ClientManager{
		Storage: storage,
	}
}

func (cm *ClientManager) Register(conn *ClientConn, nickname string) *Client {
	return cm.Storage.Register(conn, nickname)
}

func (cm *ClientManager) Logout(conn *ClientConn) error {
	if !cm.Storage.Logout(conn) {
		return ErrLogoutClientNotFount
	}
	return nil
}

func (cm *ClientManager) ExitRoom(conn *ClientConn) {
	cm.Storage.ExitRoom(conn)
}

func (cm *ClientManager) GetClient(conn *ClientConn) (*Client, bool) {
	return cm.Storage.GetClient(conn)
}
