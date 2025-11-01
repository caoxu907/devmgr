// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"devmgr/internal/consts"
)

type (
	IMQClient interface {
		// Init initializes the producer and starts the consume loop.
		Start() error
		// Publish sends a message synchronously.
		Publish(topic string, key []byte, payload []byte) error
		// Register subscription, execute all registered subscriptions after connecting to Broker.
		RegisterSubscription(topic string, handler consts.MessageHandler)
		// Close shuts down the producer and consumer gracefully.
		Close()
	}
)

var (
	localMQClient IMQClient
)

func MQClient() IMQClient {
	if localMQClient == nil {
		panic("implement not found for interface IMQClient, forgot register?")
	}
	return localMQClient
}

func RegisterMQClient(i IMQClient) {
	localMQClient = i
}
