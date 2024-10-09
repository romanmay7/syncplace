package wsocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	//"github.com/tidwall/gjson"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Kind     string        `json:"kind"`
	Content  string        `json:"content"`
	Elements []interface{} `json:"elements"`
	Element  interface{}   `json:"element"`
	RoomID   string        `json:"roomId"`
	Username string        `json:"username"`
}

/*
const (
	KindBoardStateUpdate = iota + 1
	KindElementUpdate
	KindNewUserConnected
	KindUserDisconnected
)*/

const (
	KindBoardStateUpdate = "1"
	KindElementUpdate    = "2"
	KindNewUserConnected = "3"
	KindUserDisconnected = "4"
)

// WRITE TO CLIENT's WEBSOCKET **********************************************************************

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		//Pick up the Message from Client's Message channel
		message, ok := <-c.Message
		if !ok {
			fmt.Println("Some problem has occured")
			return
		}
		fmt.Println("Sending message to user:" + c.Username)
		//Send Message to Client's Websocket (to UI)

		err := c.Conn.WriteJSON(message)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Socket is Open")
		}
	}
}

// READ FROM CLIENT's WEBSOCKET **********************************************************************

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		fmt.Println("GOT MESSAGE FROM UI")
		//Reading Messages sent to Client's Websocket (from UI) ----------------------------
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				fmt.Printf("error: %v", err)
			}
			break
		}

		//Deserialize and parse the Data ---------------------------------------------------

		var data Message

		deserialize_err := json.Unmarshal(m, &data)
		if deserialize_err != nil {
			fmt.Println("Problem with Deserialization")
			return
		}

		//fmt.Println(data)

		kind := data.Kind
		message := data.Content
		roomId := data.RoomID
		elementData := data.Element

		fmt.Println("kind:" + string(kind))
		fmt.Println("content:" + message)
		fmt.Println("roomId:" + roomId)
		fmt.Println(elementData)
		//Handle Different Types of Messages accordingly ---------------------------------------

		if kind == KindBoardStateUpdate {
			//UPDATE elements state in ROOM
		} else if kind == KindElementUpdate {

			fmt.Println("UPDATE elements state in ROOM")

			found := false

			for i, element := range hub.Rooms[roomId].Elements {
				// Perform type assertion (if necessary) to access element.id
				if id, ok := element.(map[string]interface{})["id"]; ok {
					if id == elementData.(map[string]interface{})["id"] {
						hub.Rooms[roomId].Elements[i] = elementData
						found = true
						break
					}
				} else {
					// Handle non-map elements gracefully (optional)
					fmt.Printf("Warning: element %d is not a map[string]interface{}\n", i)
				}
			}
			fmt.Println("ADD the new Element to the Storage")
			// If element not found, append it
			if !found {
				hub.Rooms[roomId].Elements = append(hub.Rooms[roomId].Elements, elementData)
				fmt.Println("Element added (not found in existing list)")
			}

			msg := &Message{
				Kind:     KindElementUpdate,
				Content:  message,
				Element:  elementData,
				RoomID:   c.RoomID,
				Username: c.Username,
			}
			//BROADCAST updated Element data to all other clients in same Room
			hub.Broadcast <- msg

		}

	}
}
