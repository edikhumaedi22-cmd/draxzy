package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Player merepresentasikan pemain dalam game
type Player struct {
	ID       string
	Name     string
	Hero     string
	Level    int
	HP       float64
	MaxHP    float64
	MP       float64
	MaxMP    float64
	Position Position
	Conn     *websocket.Conn
}

// Position posisi player di map
type Position struct {
	X float64
	Y float64
	Z float64
}

// BattleRoom room pertempuran
type BattleRoom struct {
	ID      string
	Players map[string]*Player
	MaxPlayers int
}

var rooms = make(map[string]*BattleRoom)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	player := &Player{
		ID:    "player_1",
		Name:  "Player1",
		Hero:  "Li Bai",
		Level: 1,
		HP:    3000,
		MaxHP: 3000,
		MP:    500,
		MaxMP: 500,
		Conn:  conn,
	}

	log.Printf("Player %s connected with hero %s", player.Name, player.Hero)

	// Game loop untuk player ini
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received from %s: %s", player.Name, string(message))
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("HOK Game Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
