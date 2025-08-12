package api

import (
	"context"
	"hope/db"
	pb "hope/proto/v1/match"
	"hope/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type MatchHandler struct{
	matchService service.MatchService
	pb.UnimplementedMatchServiceServer
}

func NewMatchHandler(matchService service.MatchService) *MatchHandler{
	return &MatchHandler{matchService: matchService}
}

func toMatchPB(m *db.Match)*pb.Match{
	if m == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !m.CreatedAt.IsZero() {
		ts = timestamppb.New(m.CreatedAt)
	}
	return &pb.Match{
		Id: m.ID,
		RiderId: m.RiderID,
		DriverId: m.DriverID,
		RideId: m.RideId,
		Status: m.Status,
		CreatedAt: ts,

	}
}

func(h *MatchHandler)RequestToJoin(ctx context.Context, req *pb.RequestToJoinRequest)(*pb.RequestToJoinResponse, error){
	if req == nil || req.GetMatch() == nil {
		return nil, status.Error(codes.InvalidArgument, "Match Required")
	}
	
	in := req.GetMatch()
	if in.GetRideId() == "" || in.GetRiderId() == "" || in.GetDriverId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id, rider_id, and driver_id are required")
	}
	match := &db.Match{
		RiderID:  in.GetRiderId(),
		DriverID: in.GetDriverId(),
		RideId:   in.GetRideId(),
	}
	if err := h.matchService.RequestToJoin(ctx, match); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "request failed: %v", err)
	}
	return &pb.RequestToJoinResponse{Match: toMatchPB(match)}, nil
	
}


func (h *MatchHandler) AcceptRequest(ctx context.Context, req *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
	if req == nil || req.GetMatchId() == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}
	if err := h.matchService.AcceptRequest(ctx, req.GetMatchId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "accept failed: %v", err)
	}
	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}
	return &pb.AcceptRequestResponse{Match: toMatchPB(m)}, nil
}


func (h *MatchHandler) RejectRequest(ctx context.Context, req *pb.RejectRequestRequest) (*pb.RejectRequestResponse, error) {
	if req == nil || req.GetMatchId() == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}
	if err := h.matchService.RejectRequest(ctx, req.GetMatchId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "reject failed: %v", err)
	}
	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}
	return &pb.RejectRequestResponse{Match: toMatchPB(m)}, nil
}


func (h *MatchHandler) CompleteMatch(ctx context.Context, req *pb.CompleteMatchRequest) (*pb.CompleteMatchResponse, error) {
	if req == nil || req.GetMatchId() == "" {
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
	if req == nil || req.GetMatchId() == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}
	m, err := h.matchService.GetMatchByID(ctx, req.GetMatchId())
	if err != nil || m == nil || m.ID == "" {
		return nil, status.Error(codes.NotFound, "match not found")
	}
	return &pb.GetMatchResponse{Match: toMatchPB(m)}, nil
}


func (h *MatchHandler) ListMatchesByRide(ctx context.Context, req *pb.ListMatchesByRideRequest) (*pb.ListMatchesByRideResponse, error) {
	if req == nil || req.GetRideId() == "" {
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
	if req == nil || req.GetRiderId() == "" {
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
