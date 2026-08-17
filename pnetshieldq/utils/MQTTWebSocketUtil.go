package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MQTTWebSocketUtil struct {
	wsClients map[string]*websocket.Conn
}

func NewMQTTWebSocketUtil() *MQTTWebSocketUtil {
	return &MQTTWebSocketUtil{
		wsClients: make(map[string]*websocket.Conn),
	}
}

// 分发
func (util *MQTTWebSocketUtil) HandleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	topicParts := strings.Split(msg.Topic(), "/")
	receiverID, _ := strconv.Atoi(topicParts[1])
	//这里调用了websocket传输mqtt收到的消息
	wsConn, ok := util.wsClients[strconv.Itoa(receiverID)]
	if ok {
		wsConn.WriteMessage(websocket.TextMessage, msg.Payload())
	}
}

// user发起连接后 提取jwt管理一对一websocket
// 登陆后 前端收到jwt 就可以发起websocket连接
func (util *MQTTWebSocketUtil) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	var token string
	token = r.URL.Query().Get("token")
	userID, err := GetUserIDFromTokenString(token)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	util.wsClients[strconv.Itoa(userID)] = wsConn
	go func() {
		for {
			messageType, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			if messageType == websocket.TextMessage && string(message) == "pong" {
				log.Println("Received pong message")
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := wsConn.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()
}
