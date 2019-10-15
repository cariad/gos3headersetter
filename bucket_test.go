package gos3headersetter

import (
	"testing"
)

func TestBucketSelectionAsString(t *testing.T) {
	b := Bucket{Bucket: "bucket"}
	assertEqualString(t, "selection", b.SelectionAsString(), "all objects in bucket \"bucket\"")
}

func TestBucketSelectionAsStringWithPrefix(t *testing.T) {
	b := Bucket{Bucket: "bucket", KeyPrefix: "prefix"}
	assertEqualString(t, "selection", b.SelectionAsString(), "objects in bucket \"bucket\" with key prefix \"prefix\"")
}

func TestMakeListObjectsV2Input(t *testing.T) {
	b := Bucket{Bucket: "bucket"}
	in := b.makeListObjectsV2Input()
	assertEqualString(t, "bucket", *in.Bucket, "bucket")
	assertNilStringPtr(t, "prefix", in.Prefix)
}

func TestMakeListObjectsV2InputWithPrefix(t *testing.T) {
	b := Bucket{Bucket: "bucket", KeyPrefix: "prefix"}
	in := b.makeListObjectsV2Input()
	assertEqualString(t, "bucket", *in.Bucket, "bucket")
	assertEqualString(t, "prefix", *in.Prefix, "prefix")
}
