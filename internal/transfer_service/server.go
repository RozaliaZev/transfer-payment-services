package transfer_service

import (
	"context"
	"encoding/json"
	"fmt"
	ts "services/api/transfer_service"
	"google.golang.org/protobuf/types/known/anypb"
)

type TransferServiceServer struct {
	ts.TransferServiceServer
}

func (t TransferServiceServer) SendTransferRequest(ctx context.Context, req *ts.TransferRequest) (*ts.TransferResponse, error) {

	message, err := json.Marshal(req)
	if err != nil {
		return &ts.TransferResponse{
			Success:        false,
			ErrorMessage:   fmt.Sprint(err),
			AdditionalData: nil,
		}, err
	}

	return &ts.TransferResponse{
		Success:        true,
		ErrorMessage:   "",
		AdditionalData: &anypb.Any{
			TypeUrl: "http://type.googleapis.com/google.protobuf.StringValue",
			Value:   message,
		},
	}, nil
}

func (t TransferServiceServer) ProcessTransferData(ctx context.Context, req *ts.TransferRequest) (*ts.TransferResponse, error) {

	if req.AdditionalData != nil {
		return &ts.TransferResponse{
			Success:        false,
			ErrorMessage:   string(req.AdditionalData.Value),
			AdditionalData: nil,
		}, fmt.Errorf("it is impossible to continue processing the application, %v", string(req.AdditionalData.Value))
	}

	return &ts.TransferResponse{}, nil
}
