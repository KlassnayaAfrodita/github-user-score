package controllers

import (
	"github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/services"
	pb "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/pkg/pb/scoring-manager/api"
)

type ScoringManagerHandler struct {
	pb.UnimplementedScoringManagerServiceServer
	service services.ScoringManagerServiceInterface
}

func NewCollectorHandler(service services.ScoringManagerServiceInterface) *ScoringManagerHandler {
	return &ScoringManagerHandler{
		service: service,
	}
}
