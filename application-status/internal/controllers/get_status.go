package controllers

import (
	"context"
	api "github.com/KlassnayaAfrodita/github-user-score/application-status/pkg/pb/application-status/api"
)

func (controller *ApplicationStatusController) GetScore(ctx context.Context, request *api.GetScoreRequest) (*api.GetScoreResponse, error) {
	resp, err := controller.service.GetStatus(ctx, request.GetApplicationId())
	return &api.GetScoreResponse{
		Scoring: resp.Scoring,
		Status:  api.ScoringStatus(resp.Status),
	}, err
}
