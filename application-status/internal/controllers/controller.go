package controllers

import (
	"github.com/KlassnayaAfrodita/github-user-score/application-status/internal/services"
	api "github.com/KlassnayaAfrodita/github-user-score/application-status/pkg/pb/application-status/api"
)

type ApplicationStatusController struct {
	api.UnimplementedApplicationStatusServiceServer
	service services.ApplicationStatusClientInterface
}

func NewApplicationStatusControlles(service services.ApplicationStatusClientInterface) *ApplicationStatusController {
	return &ApplicationStatusController{service: service}
}
