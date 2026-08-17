package utils

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func CreateMQTTClient(mqttBroker string, clientID string, username string, password string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(mqttBroker).
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).SetCleanSession(false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Failed to connect MQTT client: %v", token.Error())
		return nil, token.Error()
	}
	return client, nil
}
