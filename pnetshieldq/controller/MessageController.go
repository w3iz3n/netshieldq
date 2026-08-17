package controller

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"pnetshieldq/dao"
	"pnetshieldq/entity"
	"strconv"
	"time"
)

type MessageController struct {
	MessageDao    *dao.MessageDao
	FriendshipDao *dao.FriendshipDao
	mqttClients   mqtt.Client
	UserDao       *dao.UserDao
	wsClients     map[string]*websocket.Conn
}

func NewMessageController(messageDao *dao.MessageDao, friendshipDao *dao.FriendshipDao, mqttClients mqtt.Client, userDao *dao.UserDao) *MessageController {
	return &MessageController{
		MessageDao:    messageDao,
		FriendshipDao: friendshipDao,
		mqttClients:   mqttClients,
		UserDao:       userDao,
		wsClients:     make(map[string]*websocket.Conn),
	}
}
func (mc *MessageController) GetMessageStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := mc.MessageDao.CountMessageTypesByHour()
	if err != nil {
		http.Error(w, "Failed to get message statistics", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

type MessageRequest struct {
	ID           int64  `json:"id"`
	SenderID     uint   `json:"senderId"`
	ReceiverID   uint   `json:"receiverId"`
	Text         string `json:"text"`
	Timestamp    string `json:"timestamp"`
	IsHistorical bool   `json:"isHistorical"`
	MessageType  string `json:"messageType"`
}

// @Summary Send a message
// @Description Send a message to a user
// @Accept  json
// @Produce  json
// @Param messageRequest body controller.MessageRequest true "Message request"
// @Success 201 {object} entity.Message
// @Router /message/send [post]
func (mc *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	var messageRequest MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	timestamp, err := time.Parse(time.RFC3339, messageRequest.Timestamp)
	if err != nil {
		http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
		return
	}

	message := entity.Message{
		SenderID:    messageRequest.SenderID,
		ReceiverID:  messageRequest.ReceiverID,
		Content:     messageRequest.Text,
		Timestamp:   timestamp,
		MessageType: messageRequest.MessageType,
		Status:      "unread",
		Filename:    "",
	}

	if err := mc.MessageDao.Create(&message); err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(struct {
		SenderID     uint   `json:"senderId"`
		ReceiverID   uint   `json:"receiverId"`
		Content      string `json:"content"`
		Timestamp    string `json:"timestamp"`
		MessageType  string `json:"message-type"`
		IsHistorical bool   `json:"isHistorical"`
	}{
		SenderID:     message.SenderID,
		ReceiverID:   message.ReceiverID,
		Content:      message.Content,
		Timestamp:    message.Timestamp.Format(time.RFC3339),
		MessageType:  message.MessageType,
		IsHistorical: false,
	})

	mqttClient := mc.mqttClients
	mqttClient.Publish("user/"+strconv.Itoa(int(message.ReceiverID))+"/"+strconv.Itoa(int(message.SenderID)), 1, true, string(payload))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

//函数接收不需要了 因为在发送的时候存了一遍 数据库 然后如果要发往前端的话在mqtt的回调函数会调用websocket分发到对应的前端 如果对方在线的话

// @Summary Get historical messages between two users
// @Description Get all messages between two users
// @Accept  json
// @Produce  json
// @Param userID path int true "User ID"
// @Param friendID path int true "Friend ID"
// @Success 200 {array} map[string]interface{}
// @Router /message/history/{userID}/{friendID} [get]
func (mc *MessageController) GetHistoricalMessages(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID和朋友ID
	userID, err := strconv.Atoi(mux.Vars(r)["userID"])
	friendID, err := strconv.Atoi(mux.Vars(r)["friendID"])
	fmt.Println(userID, friendID)
	if err != nil {
		http.Error(w, "Invalid user ID or friend ID", http.StatusBadRequest)
		return
	}

	// 从数据库中获取历史消息
	messages, err := mc.MessageDao.GetMessagesBetweenUsers(uint(userID), uint(friendID))
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	// 将每个消息转换为指定的结构体格式，并将所有消息放入一个数组中
	var payload []map[string]interface{}
	for _, message := range messages {
		payload = append(payload, map[string]interface{}{
			"MessageId":    message.MessageID,
			"SenderID":     message.SenderID,
			"ReceiverID":   message.ReceiverID,
			"Content":      message.Content,
			"Timestamp":    message.Timestamp.Format(time.RFC3339),
			"MessageType":  message.MessageType,
			"IsHistorical": true, // 历史消息
		})
	}

	// 将消息数组返回给客户端
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
