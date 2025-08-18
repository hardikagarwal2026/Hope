package api

import (
	"context"
	"hope/db"
	"hope/middleware"
	pb "hope/proto/v1/review"
	"hope/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReviewHandler struct {
	reviewService service.ReviewService
	pb.UnimplementedReviewServiceServer
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func toReviewPB(r *db.Review) *pb.Review {
	if r == nil {
		return nil
	}
	var ts *timestamppb.Timestamp
	if !r.CreatedAt.IsZero() {
		ts = timestamppb.New(r.CreatedAt)
	}
	return &pb.Review{
		Id:         r.ID,
		RideId:     r.RideID,
		FromUserId: r.FromUserID,
		ToUserId:   r.ToUserID,
		Score:      int32(r.Score),
		Comment:    r.Comment,
		CreatedAt:  ts,
	}
}

func (h *ReviewHandler) SubmitReview(ctx context.Context, req *pb.SubmitReviewRequest) (*pb.SubmitReviewResponse, error) {
	if req == nil || req.GetRideId() == "" || req.GetToUserId() == "" || req.GetScore() == 0 {
		return nil, status.Error(codes.InvalidArgument, "ride_id, to_user_id, score are required")
	}

	callerID, ok := middleware.UserIDFromContext(ctx)
	if !ok || callerID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}

	if req.GetScore() < 1 || req.GetScore() > 5 {
		return nil, status.Error(codes.InvalidArgument, "score must be 1..5")
	}
	if callerID == req.GetToUserId() {
		return nil, status.Error(codes.InvalidArgument, "cannot review yourself")
	}

	r := &db.Review{
		ID:         uuid.New().String(),
		RideID:     req.GetRideId(),
		FromUserID: callerID, 
		ToUserID:   req.GetToUserId(),
		Score:      int(req.GetScore()),
		Comment:    req.GetComment(),
	}

	if err := h.reviewService.SubmitReview(ctx, r); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "submit failed: %v", err)
	}
	return &pb.SubmitReviewResponse{Review: toReviewPB(r)}, nil
}

func (h *ReviewHandler) ListReviewsByUser(ctx context.Context, req *pb.ListReviewsByUserRequest) (*pb.ListReviewsByUserResponse, error) {
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	revs, err := h.reviewService.ListReviewsByUser(ctx, req.GetUserId(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.Review, 0, len(revs))
	for i := range revs {
		out = append(out, toReviewPB(&revs[i]))
	}
	return &pb.ListReviewsByUserResponse{Reviews: out}, nil
}

func (h *ReviewHandler) ListMyReviews(ctx context.Context, req *pb.ListMyReviewsRequest) (*pb.ListMyReviewsResponse, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	revs, err := h.reviewService.ListReviewsByUser(ctx, userID, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.Review, 0, len(revs))
	for i := range revs {
		out = append(out, toReviewPB(&revs[i]))
	}
	return &pb.ListMyReviewsResponse{Reviews: out}, nil
}

func (h *ReviewHandler) ListReviewsByRide(ctx context.Context, req *pb.ListReviewsByRideRequest) (*pb.ListReviewsByRideResponse, error) {
	if req == nil || req.GetRideId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ride_id is required")
	}
	revs, err := h.reviewService.ListReviewsByRide(ctx, req.GetRideId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.Review, 0, len(revs))
	for i := range revs {
		out = append(out, toReviewPB(&revs[i]))
	}
	return &pb.ListReviewsByRideResponse{Reviews: out}, nil
}

func (h *ReviewHandler) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	if req == nil || req.GetReviewId() == "" {
		return nil, status.Error(codes.InvalidArgument, "review_id is required")
	}

	callerID, ok := middleware.UserIDFromContext(ctx)
	if !ok || callerID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth")
	}
	reviews, err := h.reviewService.ListReviewsByUser(ctx, callerID, 1000)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "load failed: %v", err)
	}
	owned := false
	for i := range reviews {
		if reviews[i].ID == req.GetReviewId() {
			owned = true
			break
		}
	}
	if !owned {
		return nil, status.Error(codes.PermissionDenied, "not allowed to delete this review")
	}

	if err := h.reviewService.DeleteReview(ctx, req.GetReviewId()); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.DeleteReviewResponse{Success: true}, nil
}
