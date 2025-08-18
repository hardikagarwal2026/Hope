package api

import (
	"context"
	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/ride"
	"hope/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RideHandler struct {
	rideService service.RideService
	pb.UnimplementedRideServiceServer
}

func NewRideHandler(rideService service.RideService) *RideHandler {
	return &RideHandler{rideService: rideService}
}


func toOfferPB(o *db.RideOffer) *pb.RideOffer {
	if o == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !o.Time.IsZero() {
		ts = timestamppb.New(o.Time)
	}
	return &pb.RideOffer{
		Id:       o.ID,
		DriverId: o.DriverID,
		FromGeo:  o.FromGeo,
		ToGeo:    o.ToGeo,
		Fare:     o.Fare,
		Time:     ts,
		Seats:    int32(o.Seats),
		Status:   o.Status,
	}
}
func toRequestPB(r *db.RideRequest) *pb.RideRequest {
	if r == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !r.Time.IsZero() {
		ts = timestamppb.New(r.Time)
	}
	return &pb.RideRequest{
		Id:      r.ID,
		UserId:  r.UserID,
		FromGeo: r.FromGeo,
		ToGeo:   r.ToGeo,
		Time:    ts,
		Seats:   int32(r.Seats),
		Status:  r.Status,
	}
}

func (h *RideHandler) CreateOffer(ctx context.Context, req *pb.CreateOfferRequest) (*pb.CreateOfferResponse, error) {
	if req == nil || req.GetFromGeo() == "" || req.GetToGeo() == "" || req.GetSeats() <= 0 || req.GetTime() == nil {
		return nil, status.Error(codes.InvalidArgument, "from_geo, to_geo, time, seats are required")
	}

	driverID, ok := middleware.UserIDFromContext(ctx)
	if !ok || driverID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	offer := &db.RideOffer{
		ID:       uuid.New().String(),
		DriverID: driverID,
		FromGeo:  req.GetFromGeo(),
		ToGeo:    req.GetToGeo(),
		Fare:     req.GetFare(),
		Time:     req.GetTime().AsTime(),
		Seats:    int(req.GetSeats()),
		Status:   "active",
	}

	if err := h.rideService.CreateOffer(ctx, offer); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create offer failed: %v", err)
	}
	return &pb.CreateOfferResponse{Offer: toOfferPB(offer)}, nil
}

func (h *RideHandler) GetOffer(ctx context.Context, req *pb.GetOfferRequest) (*pb.GetOfferResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	o, err := h.rideService.GetOfferByID(ctx, req.GetId())
	if err != nil || o == nil || o.ID == "" {
		return nil, status.Error(codes.NotFound, "offer not found")
	}
	return &pb.GetOfferResponse{Offer: toOfferPB(o)}, nil
}

func (h *RideHandler) UpdateOffer(ctx context.Context, req *pb.UpdateOfferRequest) (*pb.UpdateOfferResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	callerID, _ := middleware.UserIDFromContext(ctx)
	current, err := h.rideService.GetOfferByID(ctx, req.GetId())
	if err != nil || current == nil || current.ID == "" {
		return nil, status.Error(codes.NotFound, "offer not found")
	}
	if callerID != "" && current.DriverID != callerID {
		return nil, status.Error(codes.PermissionDenied, "not your offer")
	}

	upd := &db.RideOffer{
		ID:     req.GetId(),
		Fare:   req.GetFare(),
		Seats:  int(req.GetSeats()),
		Status: req.GetStatus(),
	}
	if err := h.rideService.UpdateOffer(ctx, upd); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "update failed: %v", err)
	}
	cur, _ := h.rideService.GetOfferByID(ctx, req.GetId())
	return &pb.UpdateOfferResponse{Offer: toOfferPB(cur)}, nil
}

func (h *RideHandler) DeleteOffer(ctx context.Context, req *pb.DeleteOfferRequest) (*pb.DeleteOfferResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	callerID, _ := middleware.UserIDFromContext(ctx)
	cur, _ := h.rideService.GetOfferByID(ctx, req.GetId())
	if cur == nil || cur.ID == "" {
		return nil, status.Error(codes.NotFound, "offer not found")
	}
	if callerID != "" && cur.DriverID != callerID {
		return nil, status.Error(codes.PermissionDenied, "not your offer")
	}

	if err := h.rideService.DeleteOffer(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteOfferResponse{Success: true}, nil
}

func (h *RideHandler) ListNearbyOffers(ctx context.Context, req *pb.ListNearbyOffersRequest) (*pb.ListNearbyOffersResponse, error) {
	if req == nil || req.GetGeohashPrefix() == "" {
		return nil, status.Error(codes.InvalidArgument, "geohash_prefix is required")
	}
	list, err := h.rideService.ListNearbyOffers(ctx, req.GetGeohashPrefix(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.RideOffer, 0, len(list))
	for i := range list {
		out = append(out, toOfferPB(&list[i]))
	}
	return &pb.ListNearbyOffersResponse{Offers: out}, nil
}

func (h *RideHandler) ListMyOffers(ctx context.Context, req *pb.ListMyOffersRequest) (*pb.ListMyOffersResponse, error) {
	callerID, ok := middleware.UserIDFromContext(ctx)
	if !ok || callerID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	limit := int(req.GetLimit())
	offers, err := h.rideService.ListMyOffers(ctx, callerID, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.RideOffer, 0, len(offers))
	for i := range offers {
		out = append(out, toOfferPB(&offers[i]))
	}
	return &pb.ListMyOffersResponse{Offers: out}, nil
}

func (h *RideHandler) CreateRequest(ctx context.Context, req *pb.CreateRequestRequest) (*pb.CreateRequestResponse, error) {
	if req == nil || req.GetFromGeo() == "" || req.GetToGeo() == "" || req.GetSeats() <= 0 || req.GetTime() == nil {
		return nil, status.Error(codes.InvalidArgument, "from_geo, to_geo, time, seats are required")
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	r := &db.RideRequest{
		ID:      uuid.New().String(),
		UserID:  userID,
		FromGeo: req.GetFromGeo(),
		ToGeo:   req.GetToGeo(),
		Time:    req.GetTime().AsTime(),
		Seats:   int(req.GetSeats()),
		Status:  "active",
	}
	if s := req.GetStatus(); s != "" {
		r.Status = s
	}

	if err := h.rideService.CreateRequest(ctx, r); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create request failed: %v", err)
	}
	return &pb.CreateRequestResponse{Request: toRequestPB(r)}, nil
}

func (h *RideHandler) GetRequest(ctx context.Context, req *pb.GetRequestRequest) (*pb.GetRequestResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	r, err := h.rideService.GetRequestByID(ctx, req.GetId())
	if err != nil || r == nil || r.ID == "" {
		return nil, status.Error(codes.NotFound, "request not found")
	}
	return &pb.GetRequestResponse{Request: toRequestPB(r)}, nil
}

func (h *RideHandler) UpdateRequestStatus(ctx context.Context, req *pb.UpdateRequestStatusRequest) (*pb.UpdateRequestStatusResponse, error) {
	if req == nil || req.GetId() == "" || req.GetStatus() == "" {
		return nil, status.Error(codes.InvalidArgument, "id and status are required")
	}
	callerID, _ := middleware.UserIDFromContext(ctx)
	cur, _ := h.rideService.GetRequestByID(ctx, req.GetId())
	if cur == nil || cur.ID == "" {
		return nil, status.Error(codes.NotFound, "request not found")
	}
	if callerID != "" && cur.UserID != callerID {
		return nil, status.Error(codes.PermissionDenied, "not your request")
	}

	if err := h.rideService.UpdateRequestStatus(ctx, req.GetId(), req.GetStatus()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "update status failed: %v", err)
	}
	r, _ := h.rideService.GetRequestByID(ctx, req.GetId())
	return &pb.UpdateRequestStatusResponse{Request: toRequestPB(r)}, nil
}

func (h *RideHandler) DeleteRequest(ctx context.Context, req *pb.DeleteRequestRequest) (*pb.DeleteRequestResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	callerID, _ := middleware.UserIDFromContext(ctx)
	cur, _ := h.rideService.GetRequestByID(ctx, req.GetId())
	if cur == nil || cur.ID == "" {
		return nil, status.Error(codes.NotFound, "request not found")
	}
	if callerID != "" && cur.UserID != callerID {
		return nil, status.Error(codes.PermissionDenied, "not your request")
	}

	if err := h.rideService.DeleteRequest(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteRequestResponse{Success: true}, nil
}

func (h *RideHandler) ListNearbyRequests(ctx context.Context, req *pb.ListNearbyRequestsRequest) (*pb.ListNearbyRequestsResponse, error) {
	if req == nil || req.GetGeohashPrefix() == "" {
		return nil, status.Error(codes.InvalidArgument, "geohash_prefix is required")
	}
	list, err := h.rideService.ListNearbyRequests(ctx, req.GetGeohashPrefix(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.RideRequest, 0, len(list))
	for i := range list {
		out = append(out, toRequestPB(&list[i]))
	}
	return &pb.ListNearbyRequestsResponse{Requests: out}, nil
}

func (h *RideHandler) ListMyRequests(ctx context.Context, req *pb.ListMyRequestsRequest) (*pb.ListMyRequestsResponse, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	limit := int(req.GetLimit())

	requests, err := h.rideService.ListMyRequests(ctx, userID, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.RideRequest, 0, len(requests))
	for i := range requests {
		out = append(out, toRequestPB(&requests[i]))
	}

	return &pb.ListMyRequestsResponse{Requests: out}, nil
}
