package s3

import (
	"context"
	"regexp"

	"github.com/minio/minio-go/v7"
	mcred "github.com/minio/minio-go/v7/pkg/credentials"

	"pkg/errors"
)

type (
	S3Client interface {
		GetListOfObjectKeys(
			ctx context.Context, bucketName string, options *minio.ListObjectsOptions, regularPrefix string,
		) ([]string, error)
		GetObject(ctx context.Context, bucketName, objectPath string, options *minio.GetObjectOptions) (*minio.Object, error)

		ping(ctx context.Context, bucketName string) error

		Get() *minio.Client
	}

	s3client struct {
		client *minio.Client
	}
)

func NewS3Client(endpoint, accessKeyID, secretAccessKey string) (S3Client, error) {
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:              mcred.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:             true,
		Transport:          nil,
		Trace:              nil,
		Region:             "",
		BucketLookup:       0,
		CustomRegionViaURL: nil,
		TrailingHeaders:    false,
		CustomMD5:          nil,
		CustomSHA256:       nil,
		MaxRetries:         0,
	})
	if err != nil {
		return nil, errors.Default.Wrap(err)
	}

	return &s3client{
		client: s3Client,
	}, nil
}

func (c *s3client) Get() *minio.Client {
	return c.client
}

// GetListOfObjectKeys returns list of all objects in bucket (you can specify folder by Prefix option)
func (c *s3client) GetListOfObjectKeys(
	ctx context.Context,
	bucketName string,
	options *minio.ListObjectsOptions,
	regularPrefix string,
) ([]string, error) {
	if options == nil {
		options = &minio.ListObjectsOptions{
			WithVersions: false,
			WithMetadata: false,
			Prefix:       "",
			Recursive:    false,
			MaxKeys:      0,
			StartAfter:   "",
			UseV1:        false,
		}
	}

	if err := c.ping(ctx, bucketName); err != nil {
		return nil, err
	}

	objChan := c.client.ListObjects(ctx, bucketName, *options)

	var objs []string
	if len(regularPrefix) != 0 {
		for v := range objChan {
			if validatePrefix(regularPrefix, v.Key) {
				objs = append(objs, v.Key)
			}
		}
	} else {
		for v := range objChan {
			objs = append(objs, v.Key)
		}
	}

	return objs, nil
}

func validatePrefix(regular, prefix string) bool {
	re := regexp.MustCompile(regular)
	return re.MatchString(prefix)
}

// GetObject Do not forget close object
//
// object.Close()
func (c *s3client) GetObject(ctx context.Context, bucketName, objectPath string, options *minio.GetObjectOptions) (*minio.Object, error) {
	if options == nil {
		options = &minio.GetObjectOptions{
			ServerSideEncryption: nil,
			VersionID:            "",
			PartNumber:           0,
			Checksum:             false,
			Internal: minio.AdvancedGetOptions{
				ReplicationDeleteMarker:           false,
				IsReplicationReadyForDeleteMarker: false,
				ReplicationProxyRequest:           "",
			},
		}
	}

	if err := c.ping(ctx, bucketName); err != nil {
		return nil, err
	}

	object, err := c.client.GetObject(ctx, bucketName, objectPath, *options)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (c *s3client) ping(ctx context.Context, bucketName string) error {
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return errors.New("Access Denied: " + err.Error())
	}

	if !exists {
		return errors.New("Bucket not exists")
	}

	return nil
}
