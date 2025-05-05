package controllers

import (
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/services"
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/pb"
)

type CollectorHandler struct {
	pb.UnimplementedCollectorServiceServer
	service services.CollectorServiceInterface
}

func NewCollectorHandler(service services.CollectorServiceInterface) *CollectorHandler {
	return &CollectorHandler{
		service: service,
	}
}
