package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"chat.app/backend/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var Mongo *mongo.Client

// we declare a mongo variable which is a pointer to mongo.Client type
// which is a struct
// this is because this stores the info about the connection
var MessagesCollection *mongo.Collection

// then i declare another pointer to Collection struct
// then i use the client map to check the online status
// kind of like lazy update
// then the brodcast channel stores the message that has to
// reach evry client and keep a buffer for that to avoid any
// sort of delay
// we have the register channel for user to register
// on the channel as multiple users might want to register
// so having a channel is better than a siple list of
// array because the array will be shared now if 4 people
// register on this and since we are using diff
// threads for each of them they might try to acces the
// data simultaneously causing race
// channel is the standard way of preventing that
// unregister is the complementary of regiter
type Hub struct {
	RoomID     string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// we now write a function that creates the new
// hub for creating a new hub each hub is like a room
// since we can have multiple hubs it's better to write
// function to initialize it
// why not use a contructor directly in the struct
// this is the constructor
func NewHub(roomID string) *Hub {
	return &Hub{
		RoomID:     roomID,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// this is the run function the job of this is to
// run to for a particular hub manage the entry and
// exit of any thread which is essentially one client
func (h *Hub) Run() {
	for {
		//an infinite loop which is supposed to
		//run as long as the connection is alive as it
		//manages any new upcoming client
		//and since we are not creating a new memory
		//frequntly we won't run out of stack
		//now we have to make use of the chnnel i made
		//in the hub to send data
		//we use the select statement for this
		//so the thin is that we might need
		//any data from any channel so if we get
		//any data from any channel we must use that channel
		//rather than waiting forever
		//it choose whichever is ready simultneously
		//and random pick if multiple ready
		//in case of none ready run default we don't have here so
		//do nothing and one iteration can run only one case

		select {
		case client := <-h.Register:
			//use the client var and mark the map
			h.Clients[client] = true
		case client := <-h.Unregister:
			//if the client is in the map delete it from the map
			//and then close the channel by calling close as
			//there is no point in sending the channel of this client
			//as the client is disconnected not closing it
			//will simply result in useless data sending on offline
			//client won't throw error or blcok just inefficeint
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case messageBytes := <-h.Broadcast:
			// 1. Unmarshal
			var msg models.Message
			//messageBytes stores the data in JSON
			//format form the channel and then we
			//unmarshel it in the struct type do the msg
			//from the models package
			if err := json.Unmarshal(messageBytes, &msg); err != nil {
				//check for error
				log.Println("Error unmarshaling message:", err)
				continue
			}

			// 2. Set server-side timestamp and roomID
			//if the time stamp is 0 set it server side time
			//or use the client side timestamp only
			if msg.Timestamp.IsZero() {
				msg.Timestamp = time.Now()
			}
			//if the roomId of that message is null
			//we give it the hub's room id
			//i don't understand what are we trying to
			//acheive here??
			// if msg.RoomID == "" {
			// 	msg.RoomID = h.RoomID
			// }
			//so overwrite always for saftey
			//always do server side id upadtes
			msg.RoomID = h.RoomID

			// 3. Save to MongoDB
			//the msg of message models is saved to mongo
			go h.saveMessage(msg)

			// 4. Marshal back to JSON
			//conver back to JSON
			finalBytes, err := json.Marshal(msg)
			if err != nil {
				//check error
				log.Println("Error marshaling message:", err)
				continue
			}

			// 5. Broadcast to clients
			//this message has to be shown to all the clients
			//now we have to send the message from the
			//server to all the clients
			//loop to all the clients in the map
			for client := range h.Clients {
				select {
				//the final bytes is the client struct
				//we send that to the channel of client
				//that sends the message
				case client.Send <- finalBytes:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

// now the function to save the Message struct to the hub
// why as we got the message from client then we save the
// message from the client to the hub then from the hub
// brodcast to all other clients including the sender
// this is begin defined as a method of the hub struct
func (h *Hub) saveMessage(msg models.Message) {
	//if messageCollection what is that where is it defined??
	//the variable collections we defined above
	if MessagesCollection == nil {
		return
	}

	//go way of creating a timer for something to try this
	//if i get a response in 5 sec good other wise cancel;
	//the cancel function allows for cancelling before the timer
	//why would i need that what if i got the response and updated in less thatn
	//5 sec?

	//ok but why can't we create sperate threads for this updae
	//then we won't have to worry about this stopping
	//the program and we can wait for infintie time??
	//because too many go routines and too much memory
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel both as varible and function??
	//call when the function ends
	//why this well if the 5 sec is not compelete but the
	//prgram crashes so end it manually
	defer cancel()

	//now use that ctx timer on this mongo db upload
	//try to upload to the mongo db within 5 sec
	//otherwise cacel and throw error
	_, err := MessagesCollection.InsertOne(ctx, msg)
	if err != nil {
		log.Println("Error saving message to DB:", err)
	}
}

// hub is one room this manages the multiple rooms
type HubManager struct {
	//the map of hub ids
	Hubs map[string]*Hub
	//why a mutex i thought we were using the the advanced
	//inbuilt channels for this race condition?
	//because it's simply more eaiser with mutex here
	//as we have multiple reads and writes to the hub
	//not owned by single go routine
	//didn't get much but ok
	//RWMutex to make sure that we can do multiple reads
	mu sync.RWMutex
}

// the constructor in disguise for the HubManager
func NewHubManager() *HubManager {
	return &HubManager{
		Hubs: make(map[string]*Hub),
	}
}

func (hm *HubManager) GetHub(roomID string) *Hub {
	//another method get the rood id hub from the map in hubManager
	//lock it to read from multiple clients
	hm.mu.RLock()
	//try to find the id
	hub, ok := hm.Hubs[roomID]
	//unlock when reading ends
	hm.mu.RUnlock()

	if ok {
		//if found in read;
		return hub
	}
	//otherwise write

	//now write lock
	hm.mu.Lock()
	defer hm.mu.Unlock()
	//check again if anything created it in between
	if hub, ok := hm.Hubs[roomID]; ok {
		//if the hub is present then return the pointer
		return hub
	}

	//otherwise create a new room and return that
	hub = NewHub(roomID)
	//we writing here but we never blocked other from
	//reading or writing why??
	hm.Hubs[roomID] = hub
	//fire the thread
	go hub.Run()
	return hub
}

// now what do we do here??
func (hm *HubManager) ListRooms() []string {
	//lock the read lock
	hm.mu.RLock()
	//defer
	defer hm.mu.RUnlock()

	//make a list of rooms id
	rooms := make([]string, 0, len(hm.Hubs))
	for roomID := range hm.Hubs {
		rooms = append(rooms, roomID)
	}
	//return that list;
	//but only one return point what is the use of defer??
	return rooms
}

// CloseHub closes a hub and disconnects all its clients
func (hm *HubManager) CloseHub(roomID string) {
	//close the hud id
	//we are going to write add the write lock
	hm.mu.Lock()
	//defer
	defer hm.mu.Unlock()

	//if doesn't exist can't delete
	//how is this possible that one id can be closed
	//multiple times??i means once closed why will you call
	//close again??
	hub, exists := hm.Hubs[roomID]
	if !exists {
		return
	}

	// Close all client connections in this hub
	for client := range hub.Clients {
		//close the channel
		close(client.Send)
		//close websocket
		client.Conn.Close()
		delete(hub.Clients, client)
	}

	// Remove hub from manager
	delete(hm.Hubs, roomID)
	log.Printf("Hub for room '%s' has been closed", roomID)
}
