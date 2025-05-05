package controllers

import (
	pb "github.com/KlassnayaAfrodita/github-user-score/collector/internal/pb/collector/api"
	"github.com/KlassnayaAfrodita/github-user-score/collector/internal/services"
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
