package iface

import (
	"context"
	"csrf/ddd/usecase"
	"csrf/pb"
)

type CsrfHandler struct {
	CSRFTokenUseCate *usecase.CSRFTokenUseCase
	pb.UnimplementedCsrfServiceServer
}

func (h *CsrfHandler) GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	token, err := h.CSRFTokenUseCate.Generate(req.SessionId)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateTokenResponse{CsrfToken: token.Token}, nil
}

func (h *CsrfHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	err := h.CSRFTokenUseCate.Validate(req.CsrfToken, req.SessionId)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	}
	return &pb.ValidateTokenResponse{IsValid: true}, nil
}
