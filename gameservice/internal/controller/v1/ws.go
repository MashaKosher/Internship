package v1

import (
	"encoding/json"
	"fmt"
	"gameservice/internal/di"
	"gameservice/internal/entity"
	randomapi "gameservice/third_party/randomApi"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type wsRoutes struct {
	u di.GameService
	l di.LoggerType
	c di.ConfigType
	b di.Bus
}

const playerAmount = 2
const minDiceAmount = 0
const maxDiceAmount = 6

func InitWSRoutes(e *echo.Echo, deps di.Container) {
	r := &wsRoutes{u: deps.Services.Game, l: deps.Logger, c: deps.Config, b: deps.Bus}
	e.GET("/ws", r.handleWS)
}

func (r *wsRoutes) handleWS(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		r.l.Error("Error upgrading connection:" + err.Error())
		return err
	}

	playerID := c.QueryParam("player_id")
	if playerID == "" {
		r.l.Error("Player ID is required")
		return nil
	}

	player := &Player{
		conn: conn,
		send: make(chan []byte),
		id:   playerID,
	}

	room := roomManager.FindOrCreateRoom(player)

	r.l.Info("Room ID: " + room.id)

	go player.writePump(r.l)

	if room.full {
		r.startGame(room)
	} else {
		r.l.Info("Player is waiting for an opponent...")
	}

	return nil
}

// Игрок
type Player struct {
	conn *websocket.Conn
	send chan []byte
	id   string
	room *Room
}

// Комната для игры
type Room struct {
	id      string
	players [playerAmount]*Player
	full    bool
	mutex   sync.Mutex
}

// Менеджер комнат
type RoomManager struct {
	rooms map[string]*Room
	mutex sync.Mutex
}

// Маппинг для всех комнат
var roomManager = RoomManager{
	rooms: make(map[string]*Room),
}

// Upgrade connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func generateRoomID() string {
	return uuid.New().String()
}

func (rm *RoomManager) FindOrCreateRoom(p *Player) *Room {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for _, room := range rm.rooms {
		if !room.full {
			room.players[1] = p
			room.full = true
			p.room = room
			return room
		}
	}

	id := generateRoomID()
	newRoom := &Room{
		id:      id,
		players: [playerAmount]*Player{p, nil},
		full:    false,
	}
	rm.rooms[id] = newRoom
	p.room = newRoom
	return newRoom
}

func generateTwoRandomInt() (int, int) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(maxDiceAmount + 1), random.Intn(maxDiceAmount + 1)
}

func generateTwoRandomWithApi(c di.ConfigType, ch chan struct{}, resCh chan int) {
	client := randomapi.NewClient(c.Random.ApiKey, c.Random.ApiUrl)

	res, err := client.GenerateIntegers(playerAmount, minDiceAmount, maxDiceAmount, true)
	if err != nil {
		return
	}

	ch <- struct{}{}
	resCh <- res[0]
	resCh <- res[1]
}

// Метод для начала игры
func (r *wsRoutes) startGame(room *Room) {

	gameSettings, _ := r.u.GetGameSettings()
	// time.Sleep(5 * time.Second)
	room.mutex.Lock()
	defer room.mutex.Unlock()

	p1 := room.players[0]
	p2 := room.players[1]

	r.l.Info("Игра игроков " + fmt.Sprint(room.players[0]) + fmt.Sprint(room.players[1]) + " начнется после небольшого перерыва")

	// Бросаем кубики
	flag := false
	var roll1, roll2 int
	for !flag {

		ch := make(chan struct{})
		res := make(chan int, 2)
		go generateTwoRandomWithApi(r.c, ch, res)
		select {
		case <-time.After(time.Second * time.Duration(gameSettings.WaitingTime)):
			roll1, roll2 = generateTwoRandomInt()
			r.l.Info("Numbers generated ourselves")
		case <-ch:
			r.l.Info("Numbers generated with API")
			roll1 = <-res
			roll2 = <-res
		}

		if roll1 != roll2 {
			flag = true
		}
	}

	var gameResult entity.GameResult

	gameResult.LoseAmount = gameSettings.LoseAmount
	gameResult.WinAmount = gameSettings.WinAmount
	gameResult.GameTime = time.Now()

	result := map[string]interface{}{
		"type": "game_result",
		"p1":   roll1,
		"p2":   roll2,
	}

	var winner_id int
	var winner_res int
	var loser_id int
	var loser_res int
	var err error
	if roll1 > roll2 {
		winner_id, err = strconv.Atoi(p1.id)
		if err != nil {
			winner_id = -1
		}
		winner_res = roll1

		loser_id, err = strconv.Atoi(p2.id)
		if err != nil {
			loser_id = -1
		}
		loser_res = roll2
	} else if roll2 > roll1 {
		winner_id, err = strconv.Atoi(p2.id)
		if err != nil {
			winner_id = -1
		}
		winner_res = roll2

		loser_id, err = strconv.Atoi(p1.id)
		if err != nil {
			loser_id = -1
		}
		loser_res = roll1
	}

	result["winner"] = winner_id
	result["loser"] = loser_id

	gameResult.Winner = winner_id
	gameResult.WinnerResult = winner_res

	gameResult.Loser = loser_id
	gameResult.LoserResult = loser_res

	go r.u.SaveGame(gameResult)

	go r.b.MatchInfoProducer.SendMatchInfo(gameResult)

	////////////
	room.full = false
	room.players[0] = nil
	room.players[1] = nil
	roomManager.mutex.Lock()
	delete(roomManager.rooms, room.id)
	roomManager.mutex.Unlock()
	////////////
	msg, _ := json.Marshal(result)
	p1.send <- msg
	p2.send <- msg

}

func (p *Player) writePump(logger di.LoggerType) {
	for msg := range p.send {
		logger.Info("writePump send message to client ")
		err := p.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.Error("Write error:" + err.Error())
			break
		}
	}
}
