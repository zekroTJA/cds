package stores

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

type S3 struct {
	client   *minio.Client
	bucket   string
	basePath string
}

var _ Store = (*S3)(nil)

func NewS3(endpoint, accessKey, secretKey, region, bucket, basePath string, secure bool) (t *S3, err error) {
	t = &S3{}

	t.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	t.bucket = bucket
	t.basePath = basePath

	return t, nil
}

func (t *S3) Get(pth string) (rc io.ReadCloser, mimeType string, lastModified time.Time, err error) {
	pth = path.Join(t.basePath, pth)
	obj, err := t.client.GetObject(context.Background(), t.bucket, pth, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", time.Time{}, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, "", time.Time{}, err
	}

	mimeType = stat.ContentType
	lastModified = stat.LastModified

	return obj, mimeType, lastModified, err
}

func (t *S3) List(path string) ([]PathEntry, error) {
	cInfo := t.client.ListObjects(context.Background(), t.bucket, minio.ListObjectsOptions{
		Prefix: path,
		UseV1:  false,
	})

	var res []PathEntry
	for i := range cInfo {
		if i.Err != nil {
			return nil, i.Err
		}

		name := strings.TrimPrefix(i.Key, path)
		if name == "" {
			continue
		}

		res = append(res, PathEntry{
			IsDir: strings.HasSuffix(name, "/"),
			Name:  name,
		})
	}

	return res, nil
}

func (t *S3) IsNotExist(err error) bool {
	return minio.ToErrorResponse(err).StatusCode == http.StatusNotFound
}
