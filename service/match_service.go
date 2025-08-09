package service

import (
	"context"
	"errors"
	"hope/db"
	"hope/repository"
	"time"
)

// methods related to match service
type MatchService interface{
	RequestToJoin(ctx context.Context, match *db.Match) error
    AcceptRequest(ctx context.Context, matchID string) error
    RejectRequest(ctx context.Context, matchID string) error
    CompleteMatch(ctx context.Context, matchID string) error
    GetMatchByID(ctx context.Context, matchID string) (*db.Match, error)
    ListMatchesByRide(ctx context.Context, rideID string) ([]db.Match, error)
    ListMatchesByRider(ctx context.Context, riderID string) ([]db.Match, error)
}

// struct implementing interface, and have all the dependencies matcgrepo, rideoffer, riderequest repo
type matchService struct{
	matchrepo 			repository.MatchRepository
	rideofferepo 		repository.RideOfferRepository
	riderequestrepo 	repository.RideRequestRepository
}

//constructor to call newmatchservice
func NewMatchService(matchrepo  repository.MatchRepository, rideofferepo  repository.RideOfferRepository, riderequestrepo repository.RideRequestRepository) MatchService {
	return &matchService{matchrepo: matchrepo, rideofferepo: rideofferepo, riderequestrepo: riderequestrepo}
}


func(s matchService) RequestToJoin(ctx context.Context, match *db.Match) error {
	if match.RiderID == "" || match.DriverID == "" || match.RideId == "" {
		return errors.New("missing required fields")
	}
	match.Status = "requested"
	match.CreatedAt = time.Now() //time when match requested
	return s.matchrepo.Create(ctx,match)
}

func(s matchService) AcceptRequest(ctx context.Context, matchID string) error {
	return s.matchrepo.UpdateStatus(ctx, matchID, "accepted")	
}

func(s matchService) RejectRequest(ctx context.Context, matchID string) error {
	return s.matchrepo.UpdateStatus(ctx, matchID, "rejected")
}

func(s matchService) CompleteMatch(ctx context.Context, matchID string) error {
	return s.matchrepo.UpdateStatus(ctx, matchID, "completed")
}


func(s matchService) GetMatchByID(ctx context.Context, matchID string) (*db.Match, error){
    match, err:= s.matchrepo.FindByID(ctx, matchID)
	if err != nil || match == nil || match.ID == "" {
        return nil, errors.New("match not found")
    }
	return match, nil
}


func(s matchService) ListMatchesByRide(ctx context.Context, rideID string) ([]db.Match, error) {
	return s.matchrepo.FindByRideID(ctx, rideID)
}

func(s matchService) ListMatchesByRider(ctx context.Context, riderID string) ([]db.Match, error){
	return s.matchrepo.FindByRiderID(ctx, riderID)
}







