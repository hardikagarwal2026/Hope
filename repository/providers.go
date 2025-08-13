package repository

import(
	"github.com/google/wire"
    "gorm.io/gorm"
)


//Provider functions for each repository
// Constructs from all the repo, written here, as the provider functions
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func NewRideRequestRepository(db *gorm.DB) RideRequestRepository {
	return &rideRequestRepository{db: db}
}

func NewrideOfferRepository(db *gorm.DB) RideOfferRepository {
	return &rideOfferRepository{db: db}
}

func NewUserLocationRepository(db *gorm.DB) UserLocationRepository {
    return &userLocationRepository{db: db}
}

func NewMatchRepository(db *gorm.DB) MatchRepository{
	return &matchRepository{db: db}
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

//when multiple constrcutors depend on each
// we let wire.NewSet to figure out during compile time
//wires all the dependent repos
var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewRideRequestRepository,
	NewrideOfferRepository,
	NewUserLocationRepository,
	NewMatchRepository,
	NewChatMessageRepository,
	NewReviewRepository,
)