package handlers

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"

	"github.com/igomonov88/nimbler_key_generator/internal/keygen"
	"github.com/igomonov88/nimbler_key_generator/internal/storage"
	pb "github.com/igomonov88/nimbler_key_generator/proto"
)

// GenerateKeys used by cron job to generate keys in count from incoming request and put
// generated keys to database.
func (s *Server) GenerateKeys(ctx context.Context, req *pb.GenerateKeysRequest) (*pb.GenerateKeysResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.GenerateKeys")
	defer span.End()

	// Create a slice of incoming req.Count capacity
	keys := make([]string, 0, req.Count)

	// Generate req.Count keys
	for i := 0; i <= int(req.Count); i++ {
		key, err := keygen.Generate(ctx)
		if err != nil {
			return &pb.GenerateKeysResponse{}, status.Error(http.StatusInternalServerError, err.Error())
		}
		keys = append(keys, key)
	}

	// Add generate keys to database
	if err := storage.AddGeneratedKeys(ctx, s.DB, keys); err != nil {
		return &pb.GenerateKeysResponse{}, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &pb.GenerateKeysResponse{}, nil
}
