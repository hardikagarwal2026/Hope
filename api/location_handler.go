package api

import (
	"context"
	"hope/db"
	pb "hope/proto/v1/location"
	"hope/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LocationHandler struct{
	locationService service.LocationService
	pb.UnimplementedLocationServiceServer
}

func NewLocationHandler(locationService service.LocationService)*LocationHandler {
	return &LocationHandler{locationService: locationService}
}

func toLocationPB(l *db.UserLocation)*pb.UserLocation {
	if l == nil {
		return nil
	}
	
	var ts *timestamppb.Timestamp
	if !l.Updatedat.IsZero(){
		ts = timestamppb.New(l.Updatedat)
	}
	
	return &pb.UserLocation{
		UserId: l.UserID,
		Latitude: l.Latitude,
		Longitude: l.Longitude,
		Geohash: l.Geohash,
		UpdatedAt: ts,
	}
}

func(h *LocationHandler)UpsertLocation(ctx context.Context, req *pb.UpsertLocationRequest)(*pb.UpsertLocationResponse, error){
	if req == nil || req.GetLocation() == nil {
		return nil, status.Error(codes.InvalidArgument, "location not provided")
	}

	in := req.GetLocation()
	l := db.UserLocation{
		UserID: in.GetUserId(),
		Latitude: in.GetLatitude(),
		Longitude: in.GetLongitude(),
		Geohash: in.GetGeohash(),	
	}

	err := h.locationService.UpsertLocation(ctx, &l)
	if err != nil{
		return nil, status.Errorf(codes.InvalidArgument, "upsert failed:%v", err)
	}

	return &pb.UpsertLocationResponse{
		Location: toLocationPB(&l),
	}, nil
}


func(h *LocationHandler)GetLocationByUser(ctx context.Context, req *pb.GetLocationByUserRequest)(*pb.GetLocationByUserResponse, error){
	if req == nil || req.GetUserId() == ""{
		return nil, status.Error(codes.InvalidArgument, "user id required")
	}

	l, err := h.locationService.GetLocationByUser(ctx, req.GetUserId())
	if err != nil || l == nil{
		return nil, status.Error(codes.NotFound, "location not found")
	}
	return &pb.GetLocationByUserResponse{
		Location: toLocationPB(l),
	}, nil
}


func(h *LocationHandler)ListNearby(ctx context.Context, req *pb.ListNearbyRequest)(*pb.ListNearbyResponse, error){
	if req == nil || req.GetGeohashPrefix() == ""{
		return nil, status.Error(codes.InvalidArgument, "geohash required")
	}

	loc , err := h.locationService.ListNearby(ctx, req.GetGeohashPrefix(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.UserLocation, 0, len(loc))
	for i := range loc {
		out = append(out, toLocationPB(&loc[i]))
	}

	return &pb.ListNearbyResponse{
		Locations: out,
	}, nil
}

func(h *LocationHandler)DeleteLocation(ctx context.Context, req *pb.DeleteLocationRequest)(*pb.DeleteLocationResponse, error){
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id required")
	}

	err := h.locationService.DeleteLocation(ctx, req.GetUserId())
	if err != nil{
		return nil, status.Errorf(codes.Internal, "delete failed:%v", err)
	}

	return &pb.DeleteLocationResponse{
		Success: true,
	}, nil
}