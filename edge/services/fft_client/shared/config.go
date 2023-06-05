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

type MQTTConfig struct {
	Broker        string `short:"b" long:"broker" description:"MQTT broker location" required:"true"`
	Client        string `short:"c" long:"client" description:"MQTT client" default:"fft-client"`
	Username      string `short:"u" long:"username" description:"MQTT username" required:"true"`
	Password      string `short:"p" long:"password" description:"MQTT password" required:"true"`
	RequestsTopic string `shot:"t" long:"topic" description:"MQTT topic" default:"sound-test"`
	ResultTopic   string `long:"result_topic" description:"MQTT trigger topic"`
	QOS           int    `short:"q" long:"qos" description:"MQTT QoS" default:"2"`
}
