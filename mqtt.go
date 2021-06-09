package utilities

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// MQTTSubscribe:  setsup a subscription to a mqtt server - hides the drudgery of mqtt commands in a closure
// for the client it supplies the configuration options and waits on a callback to receive the message payload
// on the return side it has stop call which can stop the listener task from running
func MQTTSubscribe(brokerurl, topic, connName, user, passwd string, port int, uponMessage func(topic string, msg []byte)) (stop func()) {
	var client mqtt.Client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", brokerurl, port))
	opts.SetClientID(connName)
	opts.SetUsername(user)
	opts.SetPassword(passwd)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// handling the message is to be donw by the client calling code
		uponMessage(msg.Topic(), msg.Payload())
	})
	opts.OnConnect = func(client mqtt.Client) {
		log.WithFields(log.Fields{
			"url":      brokerurl,
			"port":     port,
			"user":     user,
			"password": "*****",
		}).Info("MQTT: connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Println("MQTT: Disconnected")
	}
	client = mqtt.NewClient(opts)
	go func() { // task that listens to incoming messages
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.WithFields(log.Fields{
				"err": token.Error(),
			}).Warn("failed to connect to MQTT broker")
			return
		}
		token := client.Subscribe(topic, 1, nil)
		log.WithFields(log.Fields{
			"topic": topic,
		}).Info("MQTT:Subscribed to topic")
		token.Wait()

	}()
	return func() {
		// this would be the flushing / stopping function to close loop
		client.Disconnect(250)
	}
}
