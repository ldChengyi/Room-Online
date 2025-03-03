package client

type ClientManager struct {
	Storage *ClientStorage
}

var ()

func NewClientManager(storage *ClientStorage) *ClientManager {
	return &ClientManager{
		Storage: storage,
	}
}


