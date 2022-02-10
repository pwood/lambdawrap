package impl

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

// S3 is a structure to provide an S3Fetcher used by the S3Fetch wrap in lambdawrap. It is a very simplistic
// implementation. Clients must provide an initialised AWS S3 client.
//
//   cfg, err := config.LoadDefaultConfig(ctx)
//   if err != nil ...
//
//   client := s3.NewFromConfig(cfg)
//
//   s3 := S3{S3Client: client}
type S3 struct {
	S3Client s3.Client
}

// Fetch is to be passed into the S3Fetch wrap, it will fetch the exact version of the object located in the bucket
// from the S3Entity.
func (s *S3) Fetch(ctx context.Context, e events.S3Entity) (io.ReadCloser, error) {
	req := &s3.GetObjectInput{
		Bucket:    &e.Bucket.Name,
		Key:       &e.Object.Key,
		VersionId: &e.Object.VersionID,
	}

	out, err := s.S3Client.GetObject(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s3 fetch error: %w", err)
	}

	return out.Body, nil
}
