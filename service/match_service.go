package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"hope/db"
	"hope/repository"
)

var (
	errMatchNotFound = errors.New("match not found")
	errForbidden     = errors.New("forbidden")
)

type MatchService interface {
	RequestToJoin(ctx context.Context, match *db.Match) error
	AcceptRequest(ctx context.Context, callerID, matchID string) error
	RejectRequest(ctx context.Context, callerID, matchID string) error
	CompleteMatch(ctx context.Context, matchID string) error
	GetMatchByID(ctx context.Context, matchID string) (*db.Match, error)
	ListMatchesByRide(ctx context.Context, rideID string) ([]db.Match, error)
	ListMatchesByRider(ctx context.Context, riderID string) ([]db.Match, error)
}

type matchService struct {
	matchrepo       repository.MatchRepository
	rideofferepo    repository.RideOfferRepository
	riderequestrepo repository.RideRequestRepository
}

func NewMatchService(matchrepo repository.MatchRepository, rideofferepo repository.RideOfferRepository, riderequestrepo repository.RideRequestRepository) MatchService {
	return &matchService{matchrepo: matchrepo, rideofferepo: rideofferepo, riderequestrepo: riderequestrepo}
}

func (s matchService) RequestToJoin(ctx context.Context, match *db.Match) error {
	if match == nil {
		return errMissingFields
	}
	
	match.RiderID = strings.TrimSpace(match.RiderID)
	match.RideID = strings.TrimSpace(match.RideID)

	if match.RiderID == "" || match.RideID == "" {
		return errMissingFields
	}

	offer, err := s.rideofferepo.FindByID(ctx, match.RideID)
	if err != nil || offer == nil || offer.ID == "" {
		return errors.New("ride offer not found")
	}
	match.DriverID = offer.DriverID
	if match.DriverID == "" {
		return errors.New("offer has no driver")
	}
	match.Status = "requested"
	if match.CreatedAt.IsZero() {
		match.CreatedAt = time.Now().UTC()
	}

	return s.matchrepo.Create(ctx, match)
}

func (s matchService) AcceptRequest(ctx context.Context, callerID, matchID string) error {
	callerID = strings.TrimSpace(callerID)
	matchID = strings.TrimSpace(matchID)
	if callerID == "" || matchID == "" {
		return errors.New("missing caller or match")
	}

	m, err := s.matchrepo.FindByID(ctx, matchID)
	if err != nil || m == nil || m.ID == "" {
		return errMatchNotFound
	}
	if m.DriverID != callerID {
		return errForbidden
	}
	if m.Status != "requested" {
		return errors.New("invalid state transition")
	}

	return s.matchrepo.UpdateStatus(ctx, matchID, "accepted")
}

func (s matchService) RejectRequest(ctx context.Context, callerID, matchID string) error {
	callerID = strings.TrimSpace(callerID)
	matchID = strings.TrimSpace(matchID)
	if callerID == "" || matchID == "" {
		return errors.New("missing caller or match")
	}
	m, err := s.matchrepo.FindByID(ctx, matchID)
	if err != nil || m == nil || m.ID == "" {
		return errMatchNotFound
	}
	if m.DriverID != callerID {
		return errForbidden
	}
	if m.Status != "requested" {
		return errors.New("invalid state transition")
	}

	return s.matchrepo.UpdateStatus(ctx, matchID, "rejected")
}

func (s matchService) CompleteMatch(ctx context.Context, matchID string) error {
	return s.matchrepo.UpdateStatus(ctx, strings.TrimSpace(matchID), "completed")
}

func (s matchService) GetMatchByID(ctx context.Context, matchID string) (*db.Match, error) {
	m, err := s.matchrepo.FindByID(ctx, strings.TrimSpace(matchID))
	if err != nil || m == nil || m.ID == "" {
		return nil, errMatchNotFound
	}
	return m, nil
}

func (s matchService) ListMatchesByRide(ctx context.Context, rideID string) ([]db.Match, error) {
	return s.matchrepo.FindByRideID(ctx, strings.TrimSpace(rideID))
}

func (s matchService) ListMatchesByRider(ctx context.Context, riderID string) ([]db.Match, error) {
	return s.matchrepo.FindByRiderID(ctx, strings.TrimSpace(riderID))
}
