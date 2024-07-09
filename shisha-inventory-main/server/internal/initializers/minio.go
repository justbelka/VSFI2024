package initializers

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
)

var MinioClient *minio.Client

func InitMinIO(ctx context.Context, endpoint string, accessKey string, secretKey string) error {
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	err = MinioClient.MakeBucket(ctx, "premium-images", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(ctx, "premium-images")
		if errBucketExists == nil && exists {
			zlog.Print("We already own premium-images")
		} else {
			return err
		}
	} else {
		zlog.Print("Successfully created premium-images")
	}
	premium_policy := "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":[\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::premium-images/*\"]}]}"
	err = MinioClient.SetBucketPolicy(ctx, "premium-images", premium_policy)
	if err != nil {
		zlog.Print("Error set public policy")
		log.Fatalf("Error setting public policy: %v", err)
	}

	err = MinioClient.MakeBucket(ctx, "user-images", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(ctx, "user-images")
		if errBucketExists == nil && exists {
			zlog.Print("We already own user-images")
		} else {
			return err
		}
	} else {
		zlog.Print("Successfully created user-images")
	}
	user_policy := "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":[\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::user-images/*\"]}]}"
	err = MinioClient.SetBucketPolicy(ctx, "user-images", user_policy)
	if err != nil {
		zlog.Print("Error set public policy")
		return errors.Wrap(err, "Error setting public policy: ")
	}
	return nil
}
