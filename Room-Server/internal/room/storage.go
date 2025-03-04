package room

import (
	"sync"
	"time"
)

type Room struct {
	ID        string
	Password  string
	Members   map[string]*Member
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Member struct {
	ID         string
	LastActive time.Time
	ClientAddr string
}

type RoomStorage struct {
	rooms sync.Map
	mu    sync.RWMutex
}

func NewRoomStorage() *RoomStorage {
	return &RoomStorage{
		rooms: sync.Map{},
	}
}

func (s *RoomStorage) CreateRoom(roomID, password string) *Room {
	room := &Room{
		ID:        roomID,
		Password:  password,
		Members:   make(map[string]*Member),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}

	s.rooms.Store(roomID, room)
	return room
}

func (s *RoomStorage) GetRoom(roomID string) (*Room, bool) {
	value, ok := s.rooms.Load(roomID)
	if !ok {
		return nil, false
	}
	return value.(*Room), true
}

func (s *RoomStorage) DeleteRoom(roomID string) {
	s.rooms.Delete(roomID)
}

func (s *RoomStorage) RefreshRoom(roomID string) {
	if room, ok := s.GetRoom(roomID); ok {
		room.ExpiresAt = time.Now().Add(2 * time.Hour)
	}
}

func (s *RoomStorage) ListRoom() []string {
	roomIDs := make([]string, 0)

	s.rooms.Range(func(key, _ any) bool { // 这里需要 return true
		roomIDs = append(roomIDs, key.(string))
		return true // 继续遍历
	})

	return roomIDs
}
