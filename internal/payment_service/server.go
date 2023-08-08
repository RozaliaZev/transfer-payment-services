package payment_service

import (
	"context"
	"math/rand"
	ps "services/api/payment_service"
	"services/pkg/db"
	"time"
	"google.golang.org/protobuf/types/known/anypb"
)

type PaymentServiceServer struct {
	ps.PaymentServiceServer
}

func (p PaymentServiceServer) CheckIdRepeatition(ctx context.Context, req *ps.PaymentTransferRequest) (*ps.PaymentTransferResponse, error) {

	checkResult, err := db.CheckIdRepeatition(req.SenderId, req.RequestId)
	if err != nil {
		return &ps.PaymentTransferResponse{
			Success:      false,
			ErrorMessage: "",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
				Value:   []byte("error on the database side"),
			},
		}, err
	}

	if checkResult {
		return &ps.PaymentTransferResponse{
			Success:        true,
			ErrorMessage:   "",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.BoolValue",
				Value:   []byte{0x01},
			},
		}, nil
	}

	return &ps.PaymentTransferResponse{
		Success:      false,
		ErrorMessage: "the id is not unique",
		AdditionalData: &anypb.Any{
			TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
			Value:   []byte("the id is not unique"),
		},
	}, nil
}

func (p PaymentServiceServer) RegistrationApplication(ctx context.Context, req *ps.PaymentTransferRequest) (*ps.PaymentTransferResponse, error) {
	checkResult, err := db.CheckAndChangeBalance(req.GetSenderId(), req.GetAmount())
	if err != nil {
		return &ps.PaymentTransferResponse{
			Success:      false,
			ErrorMessage: "error on the database side",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
				Value:   []byte("error on the database side"),
			},
		}, err
	}

	if !checkResult {
		return &ps.PaymentTransferResponse{
			Success:      false,
			ErrorMessage: "there are not enough funds on the balance",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
				Value:   []byte("there are not enough funds on the balance"),
			},
		}, nil
	}

	trensferPayment := &db.TransferPayment{
		SenderId:  req.GetSenderId(),
		RequestId: req.RequestId,
		Amount:    req.Amount,
	}

	err = db.AddRequestTransferPayment(trensferPayment)
	if err != nil {
		return &ps.PaymentTransferResponse{
			Success:      false,
			ErrorMessage: "error on the database side",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
				Value:   []byte("error on the database side"),
			},
		}, err
	}

	//после записи заявки в базу случайно выбрать статус для этой заявки успех / не успех
	timer := time.NewTimer(30 * time.Second)
	resultChan := make(chan string)
	go listenTimerEvent(timer, resultChan)

	requestStatus := <-resultChan

	if requestStatus == "not successful" {
		err = db.SetStatusRequestTransferPayment(trensferPayment)
		if err != nil {
			return &ps.PaymentTransferResponse{
				Success:      false,
				ErrorMessage: "error on the database side",
				AdditionalData: &anypb.Any{
					TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
					Value:   []byte("error on the database side"),
				},
			}, err
		}

		return &ps.PaymentTransferResponse{
			Success:      true,
			ErrorMessage: "",
			AdditionalData: &anypb.Any{
				TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
				Value:   []byte("unfortunately, the status of your application is not successful, your balance has remained unchanged"),
			},
		}, nil
	}

	return &ps.PaymentTransferResponse{
		Success:      true,
		ErrorMessage: "",
		AdditionalData: &anypb.Any{
			TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
			Value:   []byte("the status of your application is successful"),
		},
	}, nil

}

//функции для выбора статуса заявки
func listenTimerEvent(timer *time.Timer, resultChan chan<- string) {
	<-timer.C
	result := randomResult()
	resultChan <- result
}

func randomResult() string {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) != 0 {
		return "not successful"
	}

	return "successful"
}
