package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/Baral-Chief-of-Compliance/gameserver/types"
	"github.com/gorilla/websocket"
)

const wsServerEndpoint = "ws://localhost:40000/ws"

type Login struct {
	ClientId int    `json:"clientID"`
	Username string `json:"username"`
}

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func newGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		conn:     conn,
		clientID: rand.Intn(math.MaxInt),
		username: username,
	}
}

func (c *GameClient) login() error {
	b, err := json.Marshal(
		types.Login{
			ClientID: c.clientID,
			Username: c.username,
		})
	if err != nil {
		return err
	}
	msg := types.WSMessage{
		Type: "login",
		Data: b,
	}
	return c.conn.WriteJSON(msg)
}

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := newGameClient(conn, "test")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}

	go func() {
		var msg types.WSMessage
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("WS read error: ", err)
				continue
			}
			fmt.Println("got message from the server: ", msg)
		}
	}()

	for {
		x := rand.Intn(1000)
		y := rand.Intn(1000)
		state := types.PlayerState{
			Health:   100,
			Position: types.Position{X: x, Y: y},
		}
		b, err := json.Marshal(state)

		if err != nil {
			log.Fatal(err)
		}

		msg := types.WSMessage{
			Type: "playerState",
			Data: b,
		}
		if err := conn.WriteJSON(msg); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
