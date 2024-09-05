package handler

import (
	"context"
	"yourproject/application"
	"yourproject/domain"
	pb "yourproject/proto"
)

type CSRFHandler struct {
	service domain.CSRFService
	pb.UnimplementedCSRFServiceServer
}

func NewCSRFHandler(service domain.CSRFService) *CSRFHandler {
	return &CSRFHandler{service: service}
}

func (h *CSRFHandler) GenerateCSRFToken(ctx context.Context, req *pb.GenerateCSRFTokenRequest) (*pb.GenerateCSRFTokenResponse, error) {
	token, err := h.service.GenerateToken()
	if err != nil {
		return nil, err
	}
	return &pb.GenerateCSRFTokenResponse{Token: token.Token}, nil
}

func (h *CSRFHandler) ValidateCSRFToken(ctx context.Context, req *pb.ValidateCSRFTokenRequest) (*pb.ValidateCSRFTokenResponse, error) {
	err := h.service.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateCSRFTokenResponse{Valid: false}, err
	}
	return &pb.ValidateCSRFTokenResponse{Valid: true}, nil
}
