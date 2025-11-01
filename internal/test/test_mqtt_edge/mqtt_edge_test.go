package test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func printHelp() {
	fmt.Println("Edge Gateway MQTT Client")
	fmt.Println("Usage:")
	fmt.Println("  ./gateway [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -id string    Specify custom MQTT client ID")
	fmt.Println("  -h, --help    Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  ./gateway -id my-client-123")
	fmt.Println("  ./gateway (auto-generates UUID as client ID)")
}

func TestMain(m *testing.M) {
	flag.StringVar(&clientID, "id", "device_key_1", "Specify custom MQTT client ID")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.Usage = func() { printHelp() }
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestMqttEdge(t *testing.T) {
	// MQTT connection info
	mqttBroker := "tcp://10.17.196.182:1883"
	//mqttBroker := "tcp://192.168.80.131:1883"
	log.Printf("Connecting to MQTT broker: %s", mqttBroker)
	log.Printf("MQTT client ID: %s", clientID)

	// Create MQTT client
	mqttc := NewMqttClient(mqttBroker, clientID)

	// Register MQTT subscription
	h := NewHandler()
	subTopic := fmt.Sprintf("down/egw/%s/notify", clientID)
	handle := func(mqttClient mqtt.Client, msg mqtt.Message) {
		h.CloudMsgHandler(mqttClient, msg, clientID)
	}
	mqttc.RegisterSubscription(subTopic, handle)

	// Start MQTT client
	err := mqttc.Start()
	if err != nil {
		log.Fatalf("Failed to start MQTT client: %v", err)
	}
	log.Println("MQTT connection established successfully")

	//发送心跳
	h.StartHeartbeat(mqttc, clientID)
	defer h.StopHeartbeat()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal: %v, shutting down...", sig)

	// Cleanup resources
	h.Close()
	mqttc.Disconnect()
	log.Println("MQTT client disconnected, exiting")
}
