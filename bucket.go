package gos3headersetter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Bucket represents an S3 bucket.
type Bucket struct {
	Bucket    string
	KeyPrefix string
}

// SelectionAsString returns a human-readable summary of the objects which will
// be selected for updating.
func (b Bucket) SelectionAsString() string {
	if b.KeyPrefix == "" {
		return fmt.Sprintf("all objects in bucket \"%s\"", b.Bucket)
	}
	return fmt.Sprintf("objects in bucket \"%s\" with key prefix \"%s\"", b.Bucket, b.KeyPrefix)
}

func (b Bucket) makeListObjectsV2Input() *s3.ListObjectsV2Input {
	in := &s3.ListObjectsV2Input{Bucket: aws.String(b.Bucket)}

	if b.KeyPrefix != "" {
		in.Prefix = &b.KeyPrefix
	}

	return in
}

// Apply applies the rules to the objects within this bucket.
func (b Bucket) Apply(rules []Rule) error {
	client, err := makeClient()
	if err != nil {
		return err
	}

	in := b.makeListObjectsV2Input()

	info("Finding %s...", b.SelectionAsString())

	err = client.ListObjectsV2Pages(
		in,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, o := range page.Contents {
				object := NewObject(b.Bucket, *o.Key)
				object.Apply(rules)
			}
			return true
		})

	return nil
}
