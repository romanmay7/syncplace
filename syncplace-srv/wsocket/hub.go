package wsocket

type BoardElement struct {
	ID           string      `json:"id"`
	RoughElement interface{} `json:"roughElement"`
	Type         string      `json:"type"`
	X1           int         `json:"x1"`
	Y1           int         `json:"y1"`
	X2           int         `json:"x2"`
	Y2           int         `json:"y2"`
}

type ChatMessage struct {
	MsgID     string `json:"msgId"`
	RoomID    string `json:"roomId"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	FilePath  string `json:"filePath"`
}

type Room struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Clients      map[string]*Client `json:"clients"`
	Elements     []interface{}      `json:"elements"`
	ChatMessages []ChatMessage      `json:"chatMessages"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	PrivateMsg chan *Message
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		PrivateMsg: make(chan *Message),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					//Broadcast a message saying that the client has left the room

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {

				for _, cl := range h.Rooms[m.RoomID].Clients {
					if cl.Username != m.Username {
						cl.Message <- m
					}
				}
			}
		case m := <-h.PrivateMsg:
			if _, ok := h.Rooms[m.RoomID]; ok {

				for _, cl := range h.Rooms[m.RoomID].Clients {
					if cl.Username == m.Username {
						cl.Message <- m
					}
				}
			}
		}
	}
}
