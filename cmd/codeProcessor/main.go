package main

import (
	"Code-compilation-system/codeProcessor/code"
	"Code-compilation-system/config"
	"Code-compilation-system/repository/rabbit_mq"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	log.Print("CodeProcessor start")
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	amqpURL := fmt.Sprintf("amqp://guest:guest@%s:%d", cfg.RabbitMQ.HostName, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("connecting to rabbitMQ:  %s", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %s", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to queueDeclare: %s", err.Error())
	}
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("failed to set сh.QoS: %s", err.Error())
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to create consumer")
	}

	log.Print("CodeProcessor waiting tasks")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range messages {
			log.Printf("Message received: %s", d.Body)
			var tempMessage rabbit_mq.TaskMessageRMQ
			if err := json.Unmarshal(d.Body, &tempMessage); err != nil {
				log.Printf("Failed to unmarshal message: %s", err.Error())
				d.Nack(false, false)
				continue
			}

			log.Printf("Task in process id: %s, translator: %s, code: %s", tempMessage.TaskID, tempMessage.Translator, tempMessage.Code)

			result, err := code.RunInDocker(tempMessage.Translator, tempMessage.Code)
			var status string
			if err != nil {
				status = "error"
				result = err.Error()
				log.Printf("RunInDocker error for task %s: %s", tempMessage.TaskID, err.Error())
				sendCommit(tempMessage.TaskID, err.Error(), "error")
			} else {
				status = "ready"
			}

			if err := sendCommit(tempMessage.TaskID, result, status); err != nil {
				log.Printf("Error send commit task id: %s, %s", tempMessage.TaskID, err.Error())
				d.Nack(false, true)
				continue
			}

			log.Printf("Task ready id: %s, status: %s, result: %s", tempMessage.TaskID, status, result)
			d.Ack(false)
		}

	}()

	<-sigCh

	log.Println("Shutdown codeProcessor")
	if err := ch.Close(); err != nil {
		log.Printf("error closing channel: %s", err.Error())
	}
	if err := conn.Close(); err != nil {
		log.Printf("error close connection: %s", err.Error())
	}
	log.Println("CodeProcessor stopped")

	log.Println("App exit")
	os.Exit(0)

}

func sendCommit(taskID, result, status string) error {
	myMap := map[string]string{
		"task_id": taskID,
		"result":  result,
		"status":  status,
	}

	body, err := json.Marshal(myMap)
	if err != nil {
		return fmt.Errorf("marshal map: %s", err.Error())
	}

	resp, err := http.Post("http://app:8080/commit", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("post \"http://app:8080/commit\" %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("\"http://app:8080/commit\" returned status %d", resp.StatusCode)
	}
	return nil

}
