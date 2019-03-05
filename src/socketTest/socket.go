package socketTest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"src/github.com/gorilla/websocket"
	"time"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// "Шина" событий, регистрация клиентов и рассылка сообщений идет отсюда
type Bus struct {
	Register  chan *websocket.Conn
	broadcast chan []byte
	clients   map[*websocket.Conn]bool
}

func (b *Bus) Run() {
	for {
		select {
		case message := <-b.broadcast:
			// каждому зарегистрированному клиенту шлем сообщение
			for client := range b.clients {
				w, err := client.NextWriter(websocket.TextMessage)
				if err != nil {
					// если достучаться до клиента не удалось, то удаляем его
					delete(b.clients, client)
					continue
				}

				_, _ = w.Write(message)
			}
		case client := <-b.Register:
			// регистрируем клиентов в мапе клиентов
			log.Println("User registered")
			b.clients[client] = true
		}
	}
}

func NewBus() *Bus {
	return &Bus{
		Register:  make(chan *websocket.Conn),
		broadcast: make(chan []byte),
		clients:   make(map[*websocket.Conn]bool),
	}
}

func RunJoker(b *Bus) {
	for {
		// каждые 5 секунд ходим за шутками
		<-time.After(5 * time.Second)
		log.Println("Its joke time!")
		b.broadcast <- getJoke()
	}
}

type Joke struct {
	ID   uint32 `json:"id"`
	Joke string `json:"joke"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value Joke   `json:"value"`
}

func getJoke() []byte {
	c := http.Client{}
	resp, err := c.Get("http://api.icndb.com/jokes/random?limitTo=[nerdy]")
	if err != nil {
		return []byte("jokes API not responding")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	joke := JokeResponse{}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return []byte("Joke error")
	}

	return []byte(joke.Value.Joke)
}
