package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"hope/db"
	"hope/repository"

	"github.com/google/uuid"
)

var (
	errMatchNotFound = errors.New("match not found")
	errForbidden     = errors.New("forbidden")
)

type MatchService interface {
	RequestToJoin(ctx context.Context, match *db.Match) error
	AcceptRideRequest(ctx context.Context, driverID, requestID string) (*db.Match, error)
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

func (s matchService) AcceptRideRequest(ctx context.Context, driverID, requestID string) (*db.Match, error) {
	driverID = strings.TrimSpace(driverID)
	requestID = strings.TrimSpace(requestID)
	if driverID == "" || requestID == "" {
		return nil, errors.New("missing driver or request")
	}
	req, err := s.riderequestrepo.FindByID(ctx, requestID)
	if err != nil || req == nil || req.ID == "" {
		return nil, errors.New("ride request not found")
	}
	if req.Status != "active" {
		return nil, errors.New("request not active")
	}
	if req.UserID == driverID {
		return nil, errors.New("cannot accept own request")
	}

	offer := &db.RideOffer{
		ID:       uuid.New().String(),
		DriverID: driverID,
		FromGeo:  req.FromGeo,
		ToGeo:    req.ToGeo,
		Fare:     0,
		Time:     req.Time,
		Seats:    max(1, req.Seats),
		Status:   "matched",
	}

	match := &db.Match{
		ID:        uuid.New().String(),
		RiderID:   req.UserID,
		DriverID:  driverID,
		RideID:    offer.ID,
		Status:    "accepted",
		CreatedAt: time.Now().UTC(),
	}

	if err := s.rideofferepo.Create(ctx, offer); err != nil {
		return nil, err
	}
	if err := s.matchrepo.Create(ctx, match); err != nil {
		return nil, err
	}
	if err := s.riderequestrepo.UpdateStatus(ctx, req.ID, "matched"); err != nil {
		return nil, err
	}
	return match, nil
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
