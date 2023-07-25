package s3

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Action interface {
	PutObject(bucket, key string, data []byte) error
	GetObject(bucket, key string) ([]byte, error)
	ListObjects(bucket, prefix string) ([]string, error)
	DeleteObject(bucket, key string) error
	BucketExists(bucket string) (bool, error)
}

// S3Wrapper represents the S3 client and its configurations.
type S3Wrapper struct {
	region string
	svc    *s3.S3
}

// NewS3Wrapper creates a new instance of the S3Wrapper.
func NewS3Wrapper(region string) (*S3Wrapper, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return &S3Wrapper{
		region: region,
		svc:    s3.New(sess),
	}, nil
}

// PutObject uploads an object to the specified bucket and key.
func (w *S3Wrapper) PutObject(bucket, key string, data []byte) error {
	_, err := w.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}

// GetObject retrieves an object from the specified bucket and key.
func (w *S3Wrapper) GetObject(bucket, key string) ([]byte, error) {
	result, err := w.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return ioutil.ReadAll(result.Body)
}

// ListObjects lists objects in the specified bucket with the given prefix.
func (w *S3Wrapper) ListObjects(bucket, prefix string) ([]string, error) {
	var objects []string
	err := w.svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			objects = append(objects, *obj.Key)
		}
		return !lastPage
	})
	if err != nil {
		return nil, err
	}

	return objects, nil
}

// DeleteObject deletes an object from the specified bucket and key.
func (w *S3Wrapper) DeleteObject(bucket, key string) error {
	_, err := w.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// BucketExists checks if the specified bucket exists.
func (w *S3Wrapper) BucketExists(bucket string) (bool, error) {
	_, err := w.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		if strings.Contains(err.Error(), "status code: 404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
