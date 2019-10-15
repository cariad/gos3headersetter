package gos3headersetter

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func compare(current string, new string) (bool, string) {
	if new == "" {
		// There is no new value to apply, so no change from the current.
		return false, current
	}
	return (new != current), new
}

func endsWith(s string, suffix string) bool {
	ls := strings.ToLower(s)
	lsuffix := strings.ToLower(suffix)
	return strings.HasSuffix(ls, lsuffix)
}

func makeClient() (*s3.S3, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}
