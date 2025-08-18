package api

import (
	"context"
	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/location"
	"hope/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LocationHandler struct {
	locationService service.LocationService
	pb.UnimplementedLocationServiceServer
}

func NewLocationHandler(locationService service.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

func toLocationPB(l *db.UserLocation) *pb.UserLocation {
	if l == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !l.UpdatedAt.IsZero() {
		ts = timestamppb.New(l.UpdatedAt)
	}
	return &pb.UserLocation{
		UserId:    l.UserID,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Geohash:   l.Geohash,
		UpdatedAt: ts,
	}
}


func (h *LocationHandler) UpsertLocation(ctx context.Context, req *pb.UpsertLocationRequest) (*pb.UpsertLocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	loc := db.UserLocation{
		UserID:    userID, 
		Latitude:  req.GetLatitude(),
		Longitude: req.GetLongitude(),
		Geohash:   req.GetGeohash(), 
	}

	if err := h.locationService.UpsertLocation(ctx, &loc); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "upsert failed: %v", err)
	}

	return &pb.UpsertLocationResponse{Location: toLocationPB(&loc)}, nil
}

func (h *LocationHandler) GetLocationByUser(ctx context.Context, req *pb.GetLocationByUserRequest) (*pb.GetLocationByUserResponse, error) {
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	l, err := h.locationService.GetLocationByUser(ctx, req.GetUserId())
	if err != nil || l == nil {
		return nil, status.Error(codes.NotFound, "location not found")
	}
	return &pb.GetLocationByUserResponse{Location: toLocationPB(l)}, nil
}

func (h *LocationHandler) ListNearby(ctx context.Context, req *pb.ListNearbyRequest) (*pb.ListNearbyResponse, error) {
	if req == nil || req.GetGeohashPrefix() == "" {
		return nil, status.Error(codes.InvalidArgument, "geohash_prefix required")
	}
	locs, err := h.locationService.ListNearby(ctx, req.GetGeohashPrefix(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.UserLocation, 0, len(locs))
	for i := range locs {
		out = append(out, toLocationPB(&locs[i]))
	}
	return &pb.ListNearbyResponse{Locations: out}, nil
}


func (h *LocationHandler) DeleteMyLocation(ctx context.Context, _ *pb.DeleteMyLocationRequest) (*pb.DeleteMyLocationResponse, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	if err := h.locationService.DeleteLocation(ctx, userID); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteMyLocationResponse{Success: true}, nil
}
