package room

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	ErrRoomExists   = errors.New("[Log Error] room already exists")
	ErrInvalidPass  = errors.New("[Log Error] invalid password")
	ErrRoomNotFound = errors.New("[Log Error] room not found")
)

type RoomManager struct {
	storage *RoomStorage
}

func NewRoomManager(storage *RoomStorage) *RoomManager {
	return &RoomManager{
		storage: storage,
	}
}

// 创建房间
func (rm *RoomManager) CreateRoom(roomID, password string, creatorAddr net.Addr) (*Room, error) {
	if _, exists := rm.storage.GetRoom(roomID); exists {
		return nil, ErrRoomExists
	}

	room := rm.storage.CreateRoom(roomID, password)

	// 自动添加创建者
	room.Members["host"] = &Member{
		ID:         "host",
		LastActive: time.Now(),
		ClientAddr: creatorAddr.String(),
	}

	return room, nil
}

// 加入房间
func (rm *RoomManager) JoinRoom(roomID, password, clientID string, clientAddr net.Addr) error {
	room, exists := rm.storage.GetRoom(roomID)
	if !exists {
		return ErrRoomNotFound
	}

	if room.Password != password {
		return ErrInvalidPass
	}

	if _, exists := room.Members[clientID]; exists {
		return errors.New("[Log Error] client ID already exists")
	}

	room.Members[clientID] = &Member{
		ID:         clientID,
		LastActive: time.Now(),
		ClientAddr: clientAddr.String(),
	}

	rm.storage.RefreshRoom(roomID)
	return nil
}

// 查询房间
func (rm *RoomManager) ListRoom() string {
	roomIDs := rm.storage.ListRoom()

	var result strings.Builder
	for i, id := range roomIDs {
		result.WriteString(fmt.Sprintf("[Room %d] --> %s\n", i+1, id))
	}

	return result.String()
}
