//go:build wireinject
// +build wireinject
package di

import (
	"hope/api"
	"hope/config"
	"hope/repository"
	"hope/service"
	"github.com/google/wire"
)

type Handlers struct {
	AuthHandler     *api.AuthHandler
	ChatHandler     *api.ChatHandler
	LocationHandler *api.LocationHandler
	MatchHandler    *api.MatchHandler
	ReviewHandler   *api.ReviewHandler
	RideHandler     *api.RideHandler
	UserHandler     *api.UserHandler
}

// Provider Set
var ProviderSetService = wire.NewSet(
	config.GetAllowedDomains,
	config.InitDatabase,
	config.GetJWTSecret,
	config.GetDatabaseConfig,
	config.ProvideGoogleClientID,

	repository.NewUserRepository,
	repository.NewRideRequestRepository,
	repository.NewrideOfferRepository,
	repository.NewUserLocationRepository,
	repository.NewMatchRepository,
	repository.NewChatMessageRepository,
	repository.NewReviewRepository,

	service.NewAuthService,
	service.NewUserService,
	service.NewRideService,
	service.NewMatchService,
	service.NewChatService,
	service.NewReviewService,
	service.NewLocationService,

	api.NewAuthHandler,
	api.NewChatHandler,
	api.NewLocationHandler,
	api.NewMatchHandler,
	api.NewReviewHandler,
	api.NewRideHandler,
	api.NewUserHandler,

	wire.Struct(new(Handlers), "*"),
)

func InitApp() (*Handlers, error) {
    wire.Build(ProviderSetService)
    return nil, nil
}
