package v1

import (
	"encoding/json"
	"fmt"
	"gameservice/internal/entity"
	"gameservice/internal/usecase"
	"gameservice/pkg/logger"
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
	u usecase.Game
}

func InitWSRoutes(e *echo.Echo, gameUseCase usecase.Game) {
	r := &gameRoutes{u: gameUseCase}
	e.GET("/ws", r.handleWS)
}

func (r *gameRoutes) handleWS(c echo.Context) error {
	// Поднимаем WebSocket соединение.
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.L.Error("Error upgrading connection:" + err.Error())
		return err
	}

	// Получаем player_id из query параметра
	playerID := c.QueryParam("player_id")
	if playerID == "" {
		logger.L.Error("Player ID is required")
		return nil
	}

	// Создаём игрока
	player := &Player{
		conn: conn,
		send: make(chan []byte),
		id:   playerID,
	}

	// Находим или создаём комнату
	room := roomManager.FindOrCreateRoom(player)

	logger.L.Info("Room ID: " + room.id)

	// Запускаем горутины для чтения и записи
	go player.writePump() // для отправки сообщений (writePump),
	// go player.readPump()  // для чтения сообщений (readPump).

	// Если комната полная, начинаем игру
	if room.full {
		r.startGame(room)
	} else {
		logger.L.Info("Player is waiting for an opponent...")
	}

	return nil
}

// Игрок
type Player struct {
	conn *websocket.Conn // WebSocket соединение игрока
	send chan []byte     // Канал для отправки сообщений игроку
	id   string          // Уникальный идентификатор игрока
	room *Room           // Ссылка на комнату, в которой находится игрок
}

// Комната для игры
type Room struct {
	id      string
	players [2]*Player // максимум два игрока
	full    bool       // флаг "заполнена ли комната"
	mutex   sync.Mutex // защита данных комнаты от гонок (конкурентного доступа)
}

// Менеджер комнат
type RoomManager struct {
	rooms map[string]*Room // все активные комнаты
	mutex sync.Mutex       // защита доступа к картe комнат
} // Идея: менеджер отвечает за создание и поиск комнат.

// Маппинг для всех комнат
var roomManager = RoomManager{
	rooms: make(map[string]*Room),
} // Важно: разрешаем соединение с любым источником (на продакшене это небезопасно).

// Для WebSocket соединений
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Генерация уникального ID для комнаты
func generateRoomID() string {
	return uuid.New().String()
}

// Метод для поиска или создания комнаты
func (rm *RoomManager) FindOrCreateRoom(p *Player) *Room {
	// Блокируем доступ к комнатам (rm.mutex.Lock()).
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Ищем первую неполную комнату
	for _, room := range rm.rooms {
		// Если найдена:
		if !room.full {
			room.players[1] = p // Ставим нового игрока во вторую позицию.
			room.full = true    // Помечаем комнату как полную.
			p.room = room       // привязваем комнату к игроку
			return room
		}
	}

	// Если не нашли — создаём новую
	id := generateRoomID()
	newRoom := &Room{
		id:      id,
		players: [2]*Player{p, nil}, // Создаём с одним игроком.
		full:    false,
	}
	rm.rooms[id] = newRoom // Сохраняем в мапу.
	p.room = newRoom       // привязваем комнату к игроку
	return newRoom
}

// Метод для начала игры
func (r *gameRoutes) startGame(room *Room) {

	gameSettings, _ := r.u.GetGameSettings()
	// time.Sleep(5 * time.Second)
	// Блокируем доступ к комнате
	room.mutex.Lock()
	defer room.mutex.Unlock()

	// Получаем игроков.
	p1 := room.players[0]
	p2 := room.players[1]

	logger.L.Info("Игра игроков " + fmt.Sprint(room.players[0]) + fmt.Sprint(room.players[1]) + " начнется после небольшого перерыва")

	// Бросаем кубики (rand.Intn(7) генерирует 0–6).
	flag := false
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	roll1 := random.Intn(7)
	roll2 := random.Intn(7)
	for !flag {
		if roll1 != roll2 {
			flag = true
		}
	}

	var gameResult entity.GameResult

	gameResult.LoseAmount = gameSettings.LoseAmount
	gameResult.WinAmount = gameSettings.WinAmount
	gameResult.GameTime = time.Now()

	// Формируем JSON с результатом:
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

	////////////
	room.full = false
	room.players[0] = nil
	room.players[1] = nil
	roomManager.mutex.Lock()
	delete(roomManager.rooms, room.id)
	roomManager.mutex.Unlock()
	////////////

	// Отправляем результат обеим сторонам
	msg, _ := json.Marshal(result)
	p1.send <- msg
	p2.send <- msg

}

// Отправка сообщений игроку
func (p *Player) writePump() {
	// Ждём сообщения в канале send.
	for msg := range p.send {
		// При получении отправляем через WebSocket клиенту.
		logger.L.Info("writePump send message to client ")
		err := p.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.L.Error("Write error:" + err.Error())
			break
		}
	}
}
