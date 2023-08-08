package handlers

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	ps "services/api/payment_service"
	psClient "services/cmd/payment_service/client"
	"services/pkg/kafka"
	"time"
)

func Start() {
	for {
		client, conn, err := psClient.CreatePaymentServiceClient()
		if err != nil {
			return
		}
		defer conn.Close()

		idForCheck, err := consumeMessage("transfer_request_topic")
		if err != nil {
			log.Printf("failed to get additional data from Kafka: %v", err)
			continue
		}

		fullRequest := &ps.PaymentTransferRequest{}
		err = json.Unmarshal(idForCheck.Value, &fullRequest)
		if err != nil {
			log.Println("error unmarshal:", err)
			continue
		}
		log.Println(fullRequest.RequestId)

		responseCheckId, err := client.CheckIdRepeatition(context.Background(), fullRequest)
		if err != nil {
			log.Printf("problem on the server side(1): %v\n", err)
			continue
		}

		if err := produceMessage("check_result_topic", responseCheckId.AdditionalData); err != nil {
			log.Printf("failed to send additional data to Kafka (check_result_topic): %v", err)
			continue
		}

		if !responseCheckId.Success {
			log.Print(responseCheckId.ErrorMessage)
			continue
		}

		responseRusultApplication, err := client.RegistrationApplication(context.Background(), fullRequest)
		if err != nil {
			log.Printf("problem on the server side(2): %v\n", err)
			continue
		}

		if err := produceMessage("check_status_topic", responseRusultApplication.AdditionalData); err != nil {
			log.Printf("failed to send additional data to Kafka (check_status_topic): %v", err)
			continue
		}

		time.Sleep(time.Second) // добавить задержку и контролировать нагрузку на сервер?
	}
}

func produceMessage(topic string, data *anypb.Any) error {
	kafkaClient, err := kafka.NewKafka()
	if err != nil {
		return err
	}

	return kafkaClient.ProduceMessage(topic, data)
}

func consumeMessage(topic string) (*anypb.Any, error) {
	kafkaClient, err := kafka.NewKafka()
	if err != nil {
		return nil, err
	}

	return kafkaClient.ConsumeMessage(topic)
}
