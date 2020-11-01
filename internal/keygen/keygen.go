package keygen

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ErrFailedToGenerateKey is used when we can't generate key.
var ErrFailedToGenerateKey = errors.New("Failed to generate key")

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

func Generate(ctx context.Context) (string, error) {
	ctx, span := trace.StartSpan(ctx, "internal.keygen.Generate")
	defer span.End()

	// randomly generate string and apply time.Now parameter there
	rs := randString(8) + time.Now().String()

	// generate hash value from randomly generated string
	b := md5.Sum([]byte(rs))

	// get first 6 params from generated value
	s := fmt.Sprintf("%x", b[0:6])
	bs := []byte(s)
	if len(bs) == 0 {
		return "", ErrFailedToGenerateKey
	}

	return string(bs[0:6]), nil
}


