package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"net/http"
	ts "services/api/transfer_service"
	tsClient "services/cmd/transfer_service/client"
	"services/pkg/kafka"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/transaction", PostRequestTransferPayment)

	return router
}

func PostRequestTransferPayment(c *gin.Context) {
	client, conn, err := createTransferServiceClient()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to connect to gRPC server")
		return
	}
	defer conn.Close()

	var requestBody struct {
		SenderId  string  `json:"senderId,omitempty"`
		RequestId string  `json:"requestId,omitempty"`
		Amount    float64 `json:"amount,omitempty"`
	}
	if err := bindJSONRequest(c, &requestBody); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	createReq := &ts.TransferRequest{
		SenderId:  requestBody.SenderId,
		RequestId: requestBody.RequestId,
		Amount:    requestBody.Amount,
	}

	sendTransferResponse, err := sendTransferRequest(client, createReq)
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("error when creating request: %v", err))
		return
	}

	jsonDataId, err := marshalToJSON(createReq)
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		return
	}

	sendTransferResponse.AdditionalData = &anypb.Any{
		TypeUrl: "http://type.googleapis.com/google.protobuf.JSONValue",
		Value:   jsonDataId,
	}

	if err := produceMessage("transfer_request_topic", sendTransferResponse.AdditionalData); err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("failed to send additional data to Kafka: %v", err))
		return
	}

	resultCheckRequestId, err := consumeMessage("check_result_topic")
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get additional data from Kafka (result of check ID): %v", err))
		return
	}

	if resultCheckRequestId.TypeUrl == "http://type.googleapis.com/google.protobuf.BoolValue" {
		handleError(c, http.StatusBadRequest, fmt.Sprint(string(resultCheckRequestId.Value)))
		return
	}

	_, err = client.ProcessTransferData(context.Background(), createReq)
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("failed to transfer data between services: %v", err))
		return
	}

	finalRequestStatus, err := consumeMessage("check_status_topic")
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get additional data from Kafka (result of check status): %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": string(finalRequestStatus.Value),
	})
}

func createTransferServiceClient() (ts.TransferServiceClient, *grpc.ClientConn, error) {
	client, conn, err := tsClient.CreateTransferServiceClient()
	if err != nil {
		return nil, nil, err
	}
	return client, conn, nil
}

func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
}

func bindJSONRequest(c *gin.Context, requestBody *struct {
	SenderId  string  `json:"senderId,omitempty"`
	RequestId string  `json:"requestId,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
}) error {
	return c.BindJSON(requestBody)
}

func sendTransferRequest(client ts.TransferServiceClient, createReq *ts.TransferRequest) (*ts.TransferResponse, error) {
	return client.SendTransferRequest(context.Background(), createReq)
}

func marshalToJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
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
