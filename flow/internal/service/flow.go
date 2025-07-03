package service

import (
	"context"

	pb "flow/api/flow/v1"
)

type FlowService struct {
	pb.UnimplementedFlowServer
}

func NewFlowService() *FlowService {
	return &FlowService{}
}

func (s *FlowService) CreateFlow(ctx context.Context, req *pb.CreateFlowRequest) (*pb.CreateFlowReply, error) {
	return &pb.CreateFlowReply{}, nil
}
func (s *FlowService) UpdateFlow(ctx context.Context, req *pb.UpdateFlowRequest) (*pb.UpdateFlowReply, error) {
	return &pb.UpdateFlowReply{}, nil
}
func (s *FlowService) DeleteFlow(ctx context.Context, req *pb.DeleteFlowRequest) (*pb.DeleteFlowReply, error) {
	return &pb.DeleteFlowReply{}, nil
}
func (s *FlowService) GetFlow(ctx context.Context, req *pb.GetFlowRequest) (*pb.GetFlowReply, error) {
	return &pb.GetFlowReply{}, nil
}
func (s *FlowService) ListFlow(ctx context.Context, req *pb.ListFlowRequest) (*pb.ListFlowReply, error) {
	return &pb.ListFlowReply{}, nil
}
