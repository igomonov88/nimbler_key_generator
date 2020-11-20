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

// GetKey used to get key to used in generated url path.
func (s *Server) GetKey(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	ctx, span := trace.StartSpan(ctx, "handlers.CheckHealth")
	defer span.End()

	// Try to get generated key from database
	key, err := storage.GetKey(ctx, s.DB)
	if err != nil {
		switch err {
		case storage.ErrNoKeys:

			// If we can't get generated key from database due to no rows result,
			// which can happened when all quota generated keys are used, we will
			// generate key manually.
			key, err := keygen.Generate(ctx)
			if err != nil {
				return &pb.GetKeyResponse{}, status.Error(http.StatusInternalServerError, err.Error())
			}

			// Add generated key to database
			if err := storage.AddKey(ctx, s.DB, key, true); err != nil {
				return &pb.GetKeyResponse{}, status.Error(http.StatusInternalServerError, err.Error())
			}

			return &pb.GetKeyResponse{Key: key}, nil
		default:
			return &pb.GetKeyResponse{}, status.Error(http.StatusInternalServerError, err.Error())
		}
	}

	return &pb.GetKeyResponse{Key: key}, nil
}

