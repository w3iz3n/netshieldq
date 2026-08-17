package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"pnetshieldq/config"
	"pnetshieldq/controller"
	"pnetshieldq/dao"
	"pnetshieldq/log"
	"pnetshieldq/router"
	"pnetshieldq/utils"
)

// @version         2.0
func main() {
	log := log.GetLogger()
	totalConfig, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := config.NewDBConnection(totalConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Info("Successfully connected to database")
	}

	cmd := exec.Command("emqx", "start")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start MQTT broker: %v", err)
	} else {
		log.Info("Successfully started MQTT broker")
	}
	defer cmd.Process.Kill()

	//oss
	ossclient, _ := utils.NewOSSClient(totalConfig)
	mqttBroker := fmt.Sprintf("tcp://%s:%s", totalConfig.BrokerIP, totalConfig.BrokerPort)
	clientID := "user_" + "user11"
	mqttClient, _ := utils.CreateMQTTClient(mqttBroker, clientID, totalConfig.ACLName, totalConfig.ACLPassword)
	util := utils.NewMQTTWebSocketUtil()
	userDao := dao.NewUserDao(db)
	friendshipDao := dao.NewFriendshipDao(db)
	messageDao := dao.NewMessageDao(db)
	adminDao := dao.NewAdminDao(db)
	fileDao := dao.NewFileDao(db)
	userkeyDao := dao.NewUserKeyDAO(db)
	frienddao := dao.NewFriendCDAO(db)
	messagerecordDao := dao.NewMessageRecordDao(db)
	userAvatarDao := dao.NewUserAvatarDAO(db)
	userController := controller.NewUserController(userDao, userAvatarDao)
	friendshipController := controller.NewFriendshipController(friendshipDao, userDao, mqttClient)
	messageController := controller.NewMessageController(messageDao, friendshipDao, mqttClient, userDao)
	adminController := controller.NewAdminController(adminDao, userDao)
	fileController := controller.NewFileController(fileDao, *ossclient)
	friendcController := controller.NewFriendCController(frienddao)
	userkeycontroller := controller.NewUserKeyController(userkeyDao)
	messagerecordController := controller.NewMessageRecordController(messagerecordDao)
	logController := controller.NewLogController("./log/app.log")
	configController := controller.NewConfigController()
	useravatarcontroller := controller.NewUserAvatarController(userAvatarDao)
	serverStatusController := &controller.ServerStatusController{}

	//遗嘱信息先不处理
	//反正加载历史信息是拉数据库
	mqttClient.Subscribe("user/+/+", 1, util.HandleMQTTMessage)

	handler := router.InitRouter(userController, friendshipController, messageController, adminController, util, fileController, userkeycontroller, logController, messagerecordController, configController, friendcController, useravatarcontroller, serverStatusController)

	log.Fatal(http.ListenAndServe(":3000", handler))
}
