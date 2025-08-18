package api

import (
	"context"
	"strings"
	"time"

	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/match"
	"hope/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MatchHandler struct {
	matchService service.MatchService
	pb.UnimplementedMatchServiceServer
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func toMatchPB(m *db.Match) *pb.Match {
	if m == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !m.CreatedAt.IsZero() {
		ts = timestamppb.New(m.CreatedAt)
	}
	return &pb.Match{
		Id:        m.ID,
		RiderId:   m.RiderID,
		DriverId:  m.DriverID,
		RideId:    m.RideID,
		Status:    m.Status,
		CreatedAt: ts,
	}
}

func (h *MatchHandler) RequestToJoin(ctx context.Context, req *pb.RequestToJoinRequest) (*pb.RequestToJoinResponse, error) {
	if req == nil || strings.TrimSpace(req.GetRideId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id is required")
	}

	riderID, ok := middleware.UserIDFromContext(ctx)
	if !ok || riderID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	m := &db.Match{
		ID:        uuid.New().String(),
		RiderID:   riderID,
		RideID:    strings.TrimSpace(req.GetRideId()),
		Status:    "requested",
		CreatedAt: time.Now().UTC(),
	}

	if err := h.matchService.RequestToJoin(ctx, m); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "request failed: %v", err)
	}

	return &pb.RequestToJoinResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) AcceptRideRequest(ctx context.Context, req *pb.AcceptRideRequestRequest) (*pb.AcceptRideRequestResponse, error) {
	if req == nil || strings.TrimSpace(req.GetRequestId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "request_id is required")
	}
	driverID, ok := middleware.UserIDFromContext(ctx)
	if !ok || driverID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	m, err := h.matchService.AcceptRideRequest(ctx, driverID, req.GetRequestId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not created")
	}
	return &pb.AcceptRideRequestResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) AcceptRequest(ctx context.Context, req *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
	if req == nil || strings.TrimSpace(req.GetMatchId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	callerID, ok := middleware.UserIDFromContext(ctx)
	if !ok || callerID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	if err := h.matchService.AcceptRequest(ctx, callerID, req.GetMatchId()); err != nil {
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "forbidden"):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case strings.Contains(msg, "not found"):
			return nil, status.Error(codes.NotFound, err.Error())
		case strings.Contains(msg, "invalid state"):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Errorf(codes.InvalidArgument, "accept failed: %v", err)
		}
	}

	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}

	return &pb.AcceptRequestResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) RejectRequest(ctx context.Context, req *pb.RejectRequestRequest) (*pb.RejectRequestResponse, error) {
	if req == nil || strings.TrimSpace(req.GetMatchId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	callerID, ok := middleware.UserIDFromContext(ctx)
	if !ok || callerID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	if err := h.matchService.RejectRequest(ctx, callerID, req.GetMatchId()); err != nil {
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "forbidden"):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case strings.Contains(msg, "not found"):
			return nil, status.Error(codes.NotFound, err.Error())
		case strings.Contains(msg, "invalid state"):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Errorf(codes.InvalidArgument, "reject failed: %v", err)
		}
	}

	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}

	return &pb.RejectRequestResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) CompleteMatch(ctx context.Context, req *pb.CompleteMatchRequest) (*pb.CompleteMatchResponse, error) {
	if req == nil || strings.TrimSpace(req.GetMatchId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	if err := h.matchService.CompleteMatch(ctx, req.GetMatchId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "complete failed: %v", err)
	}

	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}

	return &pb.CompleteMatchResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.GetMatchResponse, error) {
	if req == nil || strings.TrimSpace(req.GetMatchId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}

	return &pb.GetMatchResponse{Match: toMatchPB(m)}, nil
}

func (h *MatchHandler) ListMatchesByRide(ctx context.Context, req *pb.ListMatchesByRideRequest) (*pb.ListMatchesByRideResponse, error) {
	if req == nil || strings.TrimSpace(req.GetRideId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id is required")
	}

	ms, err := h.matchService.ListMatchesByRide(ctx, req.GetRideId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}

	out := make([]*pb.Match, 0, len(ms))
	for i := range ms {
		out = append(out, toMatchPB(&ms[i]))
	}

	return &pb.ListMatchesByRideResponse{Matches: out}, nil
}

func (h *MatchHandler) ListMatchesByRider(ctx context.Context, req *pb.ListMatchesByRiderRequest) (*pb.ListMatchesByRiderResponse, error) {
	if req == nil || strings.TrimSpace(req.GetRiderId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "rider_id is required")
	}

	ms, err := h.matchService.ListMatchesByRider(ctx, req.GetRiderId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}

	out := make([]*pb.Match, 0, len(ms))
	for i := range ms {
		out = append(out, toMatchPB(&ms[i]))
	}

	return &pb.ListMatchesByRiderResponse{Matches: out}, nil
}

func (h *MatchHandler) ListMyMatches(ctx context.Context, _ *pb.ListMyMatchesRequest) (*pb.ListMyMatchesResponse, error) {
	riderID, ok := middleware.UserIDFromContext(ctx)
	if !ok || riderID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	ms, err := h.matchService.ListMatchesByRider(ctx, riderID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}

	out := make([]*pb.Match, 0, len(ms))
	for i := range ms {
		out = append(out, toMatchPB(&ms[i]))
	}

	return &pb.ListMyMatchesResponse{Matches: out}, nil
}
