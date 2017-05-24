package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

/**********************/
/* Structures         */
/**********************/

type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Kind  string `json:"kind"`
	Name  string `json:"name"`
	Color string `json:"color"`
	State string `json:"state"` // "alive" ou "dead"
	Body  []Pos  `json:"body"`

	WS   *websocket.Conn `json:"-"`
	Used bool            `json:"-"`
}

type Update struct {
	Kind   string  `json:"kind"`
	Snakes []Snake `json:"snakes"`
}

// sent when JS connected
type Init struct {
	Kind        string `json:"kind"`
	PlayersSlot []int  `json:"players_slot"`
	StateGame   string `json:"state_game"` // “waiting” or “playing” or “ended”
	MapSize     int    `json:"map_size"`
}

// allow to add kind
type KindOnly struct {
	Kind string `json:"kind"`
}

// get move
type Move struct {
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

// get nb player
type Connect struct {
	Kind string `json:"kind"`
	Slot int    `json:"slot"`
}

/********************/
/* Global variables */
/********************/

var GeneralMutex sync.Mutex // allow to lock global informations

var StateGame = Init{
	Kind:        "init",
	StateGame:   "waiting",
	MapSize:     50,
	PlayersSlot: []int{1, 2},
}

// players with position/color/name by default
var Player1 = Snake{
	Kind:  "snake",
	Name:  "p1",
	Color: "red",
	State: "alive",
	Body: []Pos{
		Pos{X: 1, Y: 3},
		Pos{X: 1, Y: 2},
		Pos{X: 1, Y: 1},
	},
}
var Player2 = Snake{
	Kind:  "snake",
	Name:  "p2",
	Color: "purple",
	State: "alive",
	Body: []Pos{
		Pos{X: 10, Y: 3},
		Pos{X: 10, Y: 2},
		Pos{X: 10, Y: 1},
	},
}

var socketList []*websocket.Conn

/**********************/
/* Functions          */
/**********************/

func main() {
	http.Handle("/", websocket.Handler(HandleClient))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func HandleClient(ws *websocket.Conn) {
	// lock mutex
	// clientlist += ws
	// unlock mutex
	if len(socketList) <= 1 {
		ws.Write(getInitMessage())
		ws.Write(getUpdateMessage())
	} else {
		Player1.WS.Write(getInitMessage())
		Player2.WS.Write(getInitMessage())
		Player1.WS.Write(getUpdateMessage())
		Player2.WS.Write(getUpdateMessage())
	}

	for {
		// 1. receive message
		var content string
		err := websocket.Message.Receive(ws, &content)
		fmt.Println("Message:", string(content))

		if err != nil {
			fmt.Println(err)
			return
		}

		// 2. find kind of message
		var k KindOnly
		err = json.Unmarshal([]byte(content), &k) // JSON text -> obj
		if err != nil {
			fmt.Println(err)
			return
		}
		kind := k.Kind

		// 3. send good function
		GeneralMutex.Lock() // lock before execute function

		if kind == "move" {
			parseMove(content)
		} else if kind == "connect" {
			setPlayer(content, ws)
			setWebsocket(ws)
		} else {
			fmt.Println("Kind unknown")
		}
		ws.Write(getUpdateMessage())

		GeneralMutex.Unlock() // unlock when it's done
	}
}

func setWebsocket(ws *websocket.Conn) {
	if len(socketList) == 0 {
		socketList = append(socketList, ws)
		Player1.WS = ws
		fmt.Println(Player1)
	} else {
		for _, socket := range socketList {
			if socket == ws {
				fmt.Println("ws already exists")
			} else {
				socketList = append(socketList, ws)
				Player2.WS = ws
				fmt.Println(Player2)
			}
		}
	}
}

func setPlayer(content string, ws *websocket.Conn) {
	// get nbPlayer
	byteJson := []byte(content)
	var snake Connect
	err := json.Unmarshal(byteJson, &snake)
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg string
	var index int
	if snake.Slot == 1 {
		Player1.Used = true
		index = 0
		msg = Player1.Name + " selected!"
	}
	if snake.Slot == 2 {
		Player2.Used = true
		index = 1
		msg = Player2.Name + " selected!"
	}
	StateGame.PlayersSlot = append(StateGame.PlayersSlot[:index], StateGame.PlayersSlot[index+1:]...)
	fmt.Println(StateGame)

	message, err := json.Marshal(msg) // transform in json
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal")
	}

	// update all clients
	if Player1.Used == true {
		Player1.WS.Write(message)
	}
	if Player2.Used == true {
		Player2.WS.Write(message)
	}
	fmt.Println(msg)
}

func parseMove(content string) {
	// get move
	byteJson := []byte(content)
	var move Move
	err := json.Unmarshal(byteJson, &move)
	if err != nil {
		fmt.Println(err)
		return
	}
	key := move.Key

	// execute move
	head := Player1.Body[0] // get head position

	if key == "left" {
		head.X -= 1
		fmt.Println("left", head)
	} else if key == "right" {
		head.X += 1
		fmt.Println("right", head)
	} else if key == "up" {
		head.Y -= 1
		fmt.Println("up", head)
	} else if key == "down" {
		head.Y += 1
		fmt.Println("down", head)
	} else {
		fmt.Println("key =", key)
	}

	var newBody []Pos                          // recreate body
	newBody = append(newBody, head)            // first put head
	newBody = append(newBody, Player1.Body...) // then old body
	newBody = newBody[0 : len(newBody)-1]      // remove queue

	Player1.Body = newBody

	// update all clients
	if Player1.Used == true {
		Player1.WS.Write(getUpdateMessage())
	}
	if Player2.Used == true {
		Player2.WS.Write(getUpdateMessage())
	}
}

// "update" in protocol
func getUpdateMessage() []byte {
	var m Update
	m.Kind = "update"
	m.Snakes = []Snake{Player1, Player2}

	message, err := json.Marshal(m) // transform in json
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal map")
	}
	return message // (Json)
}

// "init" in protocole
func getInitMessage() []byte {
	message, err := json.Marshal(StateGame) // transform in json
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal init")
	}
	return message
}
