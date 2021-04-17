// Copyright (c) 2020 SoftServe Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package shared

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type MQTTClient interface {
	SendSample(payload []byte) error
	SendResults(r *Response) error
	Subscribe(topic string, callback ResultHandler)
	Close()
}

type client struct {
	client      mqtt.Client
	topic       string
	resultTopic string
	qos         byte
}

type ResultHandler func([]byte)

func NewMQTTClient(config *MQTTConfig) MQTTClient {
	op := mqtt.NewClientOptions().
		SetClientID(config.Client).
		AddBroker(fmt.Sprintf("tcp://%s", config.Broker)).
		SetUsername(config.Username).
		SetPassword(config.Password).
		SetAutoReconnect(true).
		SetMaxReconnectInterval(1 * time.Second).
		SetPingTimeout(1 * time.Second)

	op.OnConnectionLost = func(c mqtt.Client, err error) {
		logrus.Error("Lost connection to mqtt broker", err)
	}

	cl := &client{
		client:      mqtt.NewClient(op),
		qos:         byte(config.QOS),
		topic:       config.RequestsTopic,
		resultTopic: config.ResultTopic,
	}

	token := cl.client.Connect()
	token.WaitTimeout(2 * time.Second)

	if !cl.client.IsConnected() {
		logrus.Fatal(token.Error())
	}

	return cl
}

func (c *client) SendSample(payload []byte) error {
	b := make([]byte, 8+len(payload))
	binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))

	copy(b[8:], payload)

	token := c.client.Publish(c.topic, c.qos, false, b)
	if token.WaitTimeout(2*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *client) SendResults(r *Response) error {
	if len(c.resultTopic) == 0 {
		return nil
	}

	data, err := r.ToJson()
	if err != nil {
		return err
	}

	token := c.client.Publish(c.resultTopic, c.qos, false, data)
	if token.WaitTimeout(2*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *client) Subscribe(topic string, callback ResultHandler) {
	c.client.Subscribe(topic, c.qos, func(c mqtt.Client, message mqtt.Message) {
		callback(message.Payload())
		message.Ack()
	})
}

func (c *client) Close() {
	logrus.Info("Closing broker connection")
	c.client.Disconnect(2000)
}
