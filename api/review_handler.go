package api

import (
	"context"
	"hope/db"
	pb "hope/proto/v1/review"
	"hope/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReviewHandler struct{
	reviewService service.ReviewService
	pb.UnimplementedReviewServiceServer
}


func NewReviewHandler(reviewService service.ReviewService)*ReviewHandler{
	return &ReviewHandler{reviewService: reviewService}
}

func toReviewPB(r *db.Review)*pb.Review{
	if r == nil{
		return nil
	}

	//variable of type timestamppb.Timestamp
	// use for server-side timestamp
	var ts *timestamppb.Timestamp
	if !r.CreatedAt.IsZero() {
		ts = timestamppb.New(r.CreatedAt)
	}

	return &pb.Review{
		Id: r.ID,
		RideId: r.RideID,
		FromUserId: r.FromUserID,
		ToUserId: r.ToUserID,
		Score: int32(r.Score),
		Comment: r.Comment,
		CreatedAt: ts,
	}
}


func(h *ReviewHandler)SubmitReview(ctx context.Context, req *pb.SubmitReviewRequest)(*pb.SubmitReviewResponse, error){
	if req == nil || req.GetReview() == nil {
		return nil, status.Error(codes.InvalidArgument, "review is required")
	}

	in := req.GetReview()
	r := &db.Review{
		ID: in.GetId(),
		RideID: in.GetRideId(),
		FromUserID: in.GetFromUserId(),
		ToUserID: in.GetToUserId(),
		Score: int(in.GetScore()),
		Comment: in.GetComment(),
	}
	if err := h.reviewService.SubmitReview(ctx, r);err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "submit failed:%v", err)
	}

	return &pb.SubmitReviewResponse{Review: toReviewPB(r)}, nil
}


func(h *ReviewHandler) ListReviewsByUser(ctx context.Context, req *pb.ListReviewsByUserRequest)(*pb.ListReviewsByUserResponse, error){
	if req == nil || req.GetUserId() == ""{
		return nil, status.Error(codes.InvalidArgument, "user id required")
	}

	revs, err := h.reviewService.ListReviewsByUser(ctx, req.GetUserId(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed:%v", err)
	}
	out := make([]*pb.Review, 0, len(revs))
	for i := range revs {
		out = append(out, toReviewPB(&revs[i]))
	}

	return &pb.ListReviewsByUserResponse{
		Reviews: out,
	}, nil

}

func(h *ReviewHandler)ListReviewsByRide(ctx context.Context,req *pb.ListReviewsByRideRequest)(*pb.ListReviewsByRideResponse, error){
	if req == nil || req.GetRideId() == ""{
		return  nil, status.Error(codes.InvalidArgument, "ride id required")
	}

	revs, err := h.reviewService.ListReviewsByRide(ctx, req.GetRideId())
	if err != nil{
		return nil, status.Errorf(codes.Internal, "failed to load revie:%v", err)
	}

	out := make([]*pb.Review, 0, len(revs))
	for i := range revs {
		out = append(out, toReviewPB(&revs[i]))
	}

	return &pb.ListReviewsByRideResponse{
		Reviews: out,
	}, nil
}


func(h *ReviewHandler)DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest)(*pb.DeleteReviewResponse, error){
	if req == nil || req.GetReviewId() == "" {
		return nil, status.Error(codes.InvalidArgument, "review id required")
	}

	err := h.reviewService.DeleteReview(ctx, req.GetReviewId())
	if err != nil{
		return nil, status.Errorf(codes.Internal, "delete failed:%v", err)
	}
	return &pb.DeleteReviewResponse{
		Success: true,
	}, nil
}