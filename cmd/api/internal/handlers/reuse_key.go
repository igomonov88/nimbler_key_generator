package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_key_generator/internal/storage"
	pb "github.com/igomonov88/nimbler_key_generator/proto"
)

// ReuseKeys used by cron job and only for the keys which time to leave is expired.
func (s *Server) ReuseKeys(ctx context.Context, req *pb.ReuseKeysRequest) (*pb.ReuseKeysResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.ReuseKeys")
	defer span.End()

	// Add keys to database
	if err := storage.ReuseKeys(ctx, s.DB, req.GetKeys()); err != nil {
		return &pb.ReuseKeysResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.ReuseKeysResponse{}, nil
}

