package wsocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// JSON writer
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type WsHandler struct {
	hub *Hub
}

func NewHandler(h *Hub) *WsHandler {
	return &WsHandler{
		hub: h,
	}
}

// ----------------------------------------------------------------------------
type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *WsHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {

	var req CreateRoomRequest
	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Create Room Request")
	// Check if Room already has been created
	_, exists := h.hub.Rooms[req.ID]
	if exists {
		fmt.Println("Room already has been created")
		// Send a response
		w.WriteHeader(http.StatusCreated) // Or another appropriate status code
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Room already has been created"})

	} else {
		fmt.Println("Creating New Room . . .")

		// request handling logic
		h.hub.Rooms[req.ID] = &Room{
			ID:      req.ID,
			Name:    req.Name,
			Clients: make(map[string]*Client),
		}

		// Send a response
		w.WriteHeader(http.StatusCreated) // Or another appropriate status code
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Room created successfully"})
	}

}

// ----------------------------------------------------------------------------
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ----------------------------------------------------------------------------

func (h *WsHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	roomID := vars["roomId"]
	username := r.URL.Query().Get("username")

	fmt.Println(roomID + " : " + username)

	//Get Client's Info
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       username,
		RoomID:   roomID,
		Username: username,
	}

	//if len(h.hub.Rooms[roomID].Clients) == 0:
	//get elements data(board state) from DB
	//   var elements:= callback(roomID)

	//Create MESSAGES
	m := &Message{
		Kind:     KindNewUserConnected,
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	m2 := &Message{
		Kind:         KindBoardStateUpdate,
		Content:      "Update The Board State",
		Elements:     h.hub.Rooms[roomID].Elements,
		ChatMessages: h.hub.Rooms[roomID].ChatMessages,
		RoomID:       roomID,
		Username:     username,
	}

	//Register a new Client through the Register channel
	h.hub.Register <- cl
	//Broadcast that Message (that new User has joined the Room) to all Users

	h.hub.Broadcast <- m

	if len(h.hub.Rooms[roomID].Elements) > 0 || len(h.hub.Rooms[roomID].ChatMessages) > 0 {
		if len(h.hub.Rooms[roomID].Clients) > 1 {
			fmt.Println(" Multiple Connected Clients | Recieving Room Data")
			//Send message to newly joined client to make him update his Elements Board State and get Room's Chat Messages
			cl.Message <- m2
			//Else in case this is the first Client who joined we'll make him update  his Elements Board State and get Room's Chat Messages
		} else {
			fmt.Println("First Client Joined | Recieving Room Data")
			h.hub.PrivateMsg <- m2
		}

	}

	//Launch Client's <READ/WRITE from/to SOCKET> loop functions as two separate Goroutines
	go cl.writeMessage()
	cl.readMessage(h.hub)

}

// ----------------------------------------------------------------------------
type RoomsResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *WsHandler) GetRooms(w http.ResponseWriter, r *http.Request) {

	rooms := make([]RoomsResult, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomsResult{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	WriteJSON(w, http.StatusOK, rooms)
}

type ClientsResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ClientsRequest struct {
	RoomID string `json:"roomId"`
}

// ----------------------------------------------------------------------------
func (h *WsHandler) GetClients(w http.ResponseWriter, r *http.Request) {

	var req ClientsRequest

	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var clients []ClientsResult

	if _, ok := h.hub.Rooms[req.RoomID]; !ok {
		clients = make([]ClientsResult, 0)
		WriteJSON(w, http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[req.RoomID].Clients {
		clients = append(clients, ClientsResult{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	WriteJSON(w, http.StatusOK, clients)

}
