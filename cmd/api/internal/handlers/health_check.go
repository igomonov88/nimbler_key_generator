package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	pb "github.com/igomonov88/nimbler_key_generator/proto"

	"github.com/igomonov88/nimbler_key_generator/internal/platform/database"
)

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.CheckHealth")
	defer span.End()

	if err := database.StatusCheck(ctx, s.DB); err != nil {
		return &pb.HealthCheckResponse{}, status.Error(http.StatusInternalServerError, "database is not ready")
	}

	return &pb.HealthCheckResponse{Version: "develop"}, nil
}
