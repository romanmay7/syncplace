package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/romanmay7/syncplace/filemanager"
	"github.com/romanmay7/syncplace/wsocket"
	"github.com/rs/cors"
)

// JSON writer
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// API function signature
type apiFunc func(http.ResponseWriter, *http.Request) error

// API Error Type
type ApiError struct {
	Error string `json:"error"`
}

// Convert API function type to http.HandlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//Handle the Error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

//API Server Definition ------------------------------------------------------------
//==================================================================================

type SyncPlaceAPIServer struct {
	listenAddr   string
	store        Storage
	filemgr      filemanager.FileManager
	wsockHandler *wsocket.WsHandler
	wsockHub     *wsocket.Hub
}

func NewAPIServer(listenAddr string, store Storage, filemanager filemanager.FileManager,
	wsHandler *wsocket.WsHandler, wsHub *wsocket.Hub) *SyncPlaceAPIServer {
	return &SyncPlaceAPIServer{
		listenAddr:   listenAddr,
		store:        store,
		filemgr:      filemanager,
		wsockHandler: wsHandler,
		wsockHub:     wsHub,
	}
}

func (s *SyncPlaceAPIServer) Run() {
	router := mux.NewRouter()

	//FileServer for Serving files uploaded by users to Local Server Storage in "./uploads" folder
	fs := http.FileServer(http.Dir("./uploads/"))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	router.HandleFunc("/api/upload", s.filemgr.UploadFile)
	router.HandleFunc("/api/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/api/user", makeHTTPHandleFunc(s.handleUserAccount))
	router.HandleFunc("/api/user/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetUserAccountByID), s.store))

	router.HandleFunc("/api/saveBoard/{roomId}", checkIfBoardRecordExistAndSave(makeHTTPHandleFunc(s.handleSaveBoardState), s.store, s.wsockHub))
	router.HandleFunc("/api/clearBoard/{roomId}/{userName}", clearBoardState(makeHTTPHandleFunc(s.handleClearBoardState), s.wsockHub))

	router.HandleFunc("/ws/createRoom", s.wsockHandler.CreateRoom)
	router.HandleFunc("/ws/joinRoom/{roomId}", checkBoardData(s.wsockHandler.JoinRoom, s.store, s.wsockHub))
	router.HandleFunc("/ws/getRooms", s.wsockHandler.GetRooms)
	router.HandleFunc("/ws/getClients/{roomId}", s.wsockHandler.GetClients)

	log.Println("JSON API Server is running on port", s.listenAddr)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(s.listenAddr, handler)
}

//API Handlers ------------------------------------------------------------

func (s *SyncPlaceAPIServer) handleUserAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		fmt.Println("Method : ", r.Method)
		return s.handleGetUserAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUserAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUserAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET /users
func (s *SyncPlaceAPIServer) handleGetUserAccounts(w http.ResponseWriter, r *http.Request) error {
	//userAccnt := NewUserAccount("Roman", "roman@gmail.com", "1q2w3e4R")
	accounts, err := s.store.GetUserAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// GET /user/id
func (s *SyncPlaceAPIServer) handleGetUserAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetUserAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// POST /user
func (s *SyncPlaceAPIServer) handleCreateUserAccount(w http.ResponseWriter, r *http.Request) error {
	newUserReq := new(CreateUserAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&newUserReq); err != nil {
		return err
	}

	userAccnt, err := NewUserAccount(newUserReq.UserName, newUserReq.Email, newUserReq.Password)
	if err != nil {
		return err
	}

	//Store it in DB
	if err := s.store.CreateUserAccount(userAccnt); err != nil {
		return err
	}

	/*tokenString, err := createJWT(userAccnt)
	if err != nil {
		return err
	}*/

	//fmt.Println("JWT token: ", tokenString)

	return WriteJSON(w, http.StatusOK, userAccnt)
}

// DELETE /user/id
func (s *SyncPlaceAPIServer) handleDeleteUserAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteUserAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// getID Helper function
// ------------------------------------------------------------------------------------------------
func getID(r *http.Request) (int, error) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return id, fmt.Errorf("Invalid ID given %s", idString)
	}

	return id, nil
}

//--------------------------------------------------------------------------------------------------

// POST /login
func (s *SyncPlaceAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetUserAccountByUserName(req.UserName)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", acc)

	if !acc.ValidatePassword(req.Password) {
		return fmt.Errorf("Authentication Error")
	}

	tokenString, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		UserName: req.UserName,
		Token:    tokenString,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// ----------------AUTHENTICATION(JWT)--------------------------------------------------------
//===================================================================================================

// In general we should read our secret from our ENVIRONMENT
// EXAMPLE : export JWT_SECRET = syncplace999
const JWT_SECRET = "SyncPlace999"

func createJWT(account *UserAccount) (string, error) {
	//Create Claims
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"username":  account.UserName,
	}

	//secret := os.Getenv("JWT_SECRET")
	secret := JWT_SECRET
	fmt.Println(secret)
	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println("returning JWT token")

	secretBytes := []byte(secret)
	tokenString, err := tokenData.SignedString(secretBytes)
	if err != nil {
		fmt.Println("Error while signing the token")
		fmt.Println(err)
	}
	return tokenString, err

}

// ----------------------------------------------------------------------------------------
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission Denied"})
}

// ----------------------------------------------------------------------------------------
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Auth MiddleWare")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetUserAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if account.UserName != claims["username"] {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Invalid Token"})
			return
		}

		handlerFunc(w, r)
	}
}

// ----------------------------------------------------------------------------------------
func validateJWT(tokenString string) (*jwt.Token, error) {
	//secret:=os.Getenv("JWT_SECRET")
	secret := JWT_SECRET
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

// ----------------------------------------------------------------------------------------
func checkBoardData(handlerFunc http.HandlerFunc, s Storage, hub *wsocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling checkBoardData MiddleWare")

		vars := mux.Vars(r)
		roomID := vars["roomId"]

		fmt.Println("Checking data for Room: " + roomID)
		//If there is still no Clients already connected to the Room
		if len(hub.Rooms[roomID].Clients) == 0 {

			exist, err := s.CheckIfBoardRecordExist(roomID)

			if err != nil {
				fmt.Print(err)
			}

			if exist {
				fmt.Println("There is already some data for Room: " + roomID)
				//Get elements data(Board State) from DB
				//var elements []interface{}
				elements, err := s.GetBoardState(roomID)

				if err != nil {
					WriteJSON(w, http.StatusForbidden, ApiError{Error: "Invalid Board State"})
					fmt.Println(err)
					return
				}

				//check if there is some data we got from DB
				if len(elements) > 0 {
					fmt.Println("Retrived data for the Room: " + roomID)
					hub.Rooms[roomID].Elements = elements
				}

			}

			//Get Chat Messages data  from DB
			fmt.Println("Check for CHAT Messages")
			chatMessages, err := s.GetRoomChatMessages(roomID)
			//fmt.Println(chatMessages)

			if err != nil {
				WriteJSON(w, http.StatusForbidden, ApiError{Error: "Invalid Chat Messages Data"})
				fmt.Println(err)
				return
			}

			if len(chatMessages) > 0 {
				hub.Rooms[roomID].ChatMessages = chatMessages
			}
		}

		handlerFunc(w, r)
	}
}

// ----------------------------------------------------------------------------------------

// GET /saveBoard
func (s *SyncPlaceAPIServer) handleSaveBoardState(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return WriteJSON(w, http.StatusOK, "Board State Saved")
}

// -  -  - -  -  -  -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  -
func checkIfBoardRecordExistAndSave(handlerFunc http.HandlerFunc, s Storage, hub *wsocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling checkIfBoardRecordExistAndSave MiddleWare")

		vars := mux.Vars(r)
		roomID := vars["roomId"]

		var exist bool
		//get elements data(Board State) from DB
		exist, err := s.CheckIfBoardRecordExist(roomID)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "DB Query Error"})
			fmt.Print(err)
			return
		}

		if exist {

			s.UpdateBoardStateRecord(roomID, hub.Rooms[roomID].Elements)

		} else {
			s.CreateBoardStateRecord(roomID, hub.Rooms[roomID].Name, hub.Rooms[roomID].Elements)
		}

		// Save Chat Messages Data
		// (Insert every new chat message)
		if len(hub.Rooms[roomID].ChatMessages) != 0 {
			for _, msg := range hub.Rooms[roomID].ChatMessages {

				err := s.AddNewChatMessage(msg.RoomID, msg.MsgID, msg.Timestamp, msg.Content, msg.Sender, msg.FilePath)
				if err != nil {
					WriteJSON(w, http.StatusForbidden, ApiError{Error: "DB Error: Storing New Chat Message "})
					fmt.Print(err)
					return
				}
			}
		}

		handlerFunc(w, r)
	}
}

// ----------------------------------------------------------------------------------------

// GET /clearBoard
func (s *SyncPlaceAPIServer) handleClearBoardState(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return WriteJSON(w, http.StatusOK, "Board State has been cleared")
}

// -  -  - -  -  -  -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  - -  -  -

func clearBoardState(handlerFunc http.HandlerFunc, hub *wsocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling clearBoardState MiddleWare")

		vars := mux.Vars(r)
		roomID := vars["roomId"]
		username := vars["userName"]

		//Clear all board data for this Room
		hub.Rooms[roomID].Elements = []interface{}{}

		m := &wsocket.Message{
			Kind:         wsocket.KindBoardStateUpdate,
			Content:      "Update The Board State",
			Elements:     hub.Rooms[roomID].Elements,
			ChatMessages: hub.Rooms[roomID].ChatMessages,
			RoomID:       roomID,
			Username:     username,
		}

		//update all connected Clients with new(cleared) board state
		hub.Broadcast <- m

		handlerFunc(w, r)

	}
}
