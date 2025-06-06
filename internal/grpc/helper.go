package grpc

import (
	pb "mercury/api/proto/crypto"
)

// ResponseWithBase is an interface that all gRPC responses with Base field must implement
type ResponseWithBase interface {
	GetBase() *pb.BaseResponse
}

// serviceWrapper is a generic function that wraps service calls and handles errors consistently
func serviceWrapper[T ResponseWithBase](resp T, err error) (T, error) {
	if err != nil {
		if base := resp.GetBase(); base != nil {
			base.Desc = base.Desc + ": " + err.Error()
		}
		return resp, nil
	}
	return resp, nil
}
