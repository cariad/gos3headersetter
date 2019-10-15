package gos3headersetter

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestObjectString(t *testing.T) {
	o := NewObject("bucket", "key")
	assertEqualString(t, "string", o.String(), "s3://bucket/key")
}

func TestMakeHeadObjectInput(t *testing.T) {
	o := NewObject("bucket", "key")
	in := o.makeHeadObjectInput()
	assertEqualString(t, "HeadObjectInput bucket", *in.Bucket, "bucket")
	assertEqualString(t, "HeadObjectInput key", *in.Key, "key")
}

func TestMakeCopyObjectInput(t *testing.T) {
	o := NewObject("bucket", "key")
	m := make(map[string]*string)
	m["x"] = aws.String("y")

	in := o.makeCopyObjectInput(m)
	assertEqualString(t, "CopyObjectInput bucket", *in.Bucket, "bucket")
	assertEqualString(t, "CopyObjectInput key", *in.Key, "key")
	assertEqualString(t, "CopyObjectInput source", *in.CopySource, "bucket/key")

	assertEqualString(t, "CopyObjectInput metadata", *in.Metadata["x"], "y")

	assertEqualString(t, "CopyObjectInput directive", *in.MetadataDirective, "REPLACE")

}
func RunQueueRuleEffectTest(t *testing.T, key string, header string, expect string) {
	o := NewObject("bucket", key)

	rule := Rule{
		Header: header,
		When: []When{
			When{
				Extension: ".html",
				Then:      "for-html",
			},
			When{
				Extension: ".css",
				Then:      "for-css",
			},
		},
		Else: "for-else",
	}

	o.queueRuleEffect(rule)
	assertEqualString(t, "header", o.newHeaders[header], expect)
}

func TestQueueRuleEffectForCacheControlWhenFoundFirst(t *testing.T) {
	RunQueueRuleEffectTest(t, "i.html", "Cache-Control", "for-html")
}

func TestQueueRuleEffectForCacheControlWhenFoundSecond(t *testing.T) {
	RunQueueRuleEffectTest(t, "i.css", "Cache-Control", "for-css")
}

func TestQueueRuleEffectForCacheControlElse(t *testing.T) {
	RunQueueRuleEffectTest(t, "i.png", "Cache-Control", "for-else")
}

func TestUpdateCopyObjectInput(t *testing.T) {
	o := NewObject("bucket", "key")

	o.newHeaders["Cache-Control"] = "new-cachecontrol"
	o.newHeaders["Content-Type"] = "new-contenttype"

	head := &s3.HeadObjectOutput{
		CacheControl: aws.String("current-cachecontrol"),
		ContentType:  aws.String("current-contenttype"),
	}

	in := &s3.CopyObjectInput{}

	changes := o.updateCopyObjectInput(head, in)
	assertEqualBool(t, "changes", changes, true)
	assertEqualString(t, "cache control", *in.CacheControl, "new-cachecontrol")
	assertEqualString(t, "content type", *in.ContentType, "new-contenttype")
}

func TestCalculateNewValueWithNilCurrentAndEmptyNew(t *testing.T) {
	o := NewObject("", "")
	change, new := o.calculateNewValue("Cache-Control", nil)
	assertEqualBool(t, "change", change, false)
	assertEqualString(t, "new value", new, "")
}

func TestCalculateNewValueWithNilCurrentAndRealNew(t *testing.T) {
	o := NewObject("", "")
	o.newHeaders["Cache-Control"] = "new-cachecontrol"
	change, new := o.calculateNewValue("Cache-Control", nil)
	assertEqualBool(t, "change", change, true)
	assertEqualString(t, "new value", new, "new-cachecontrol")
}

func TestCalculateNewValueWithEmptyCurrentAndEmptyNew(t *testing.T) {
	o := NewObject("", "")
	current := ""
	change, new := o.calculateNewValue("Cache-Control", &current)
	assertEqualBool(t, "change", change, false)
	assertEqualString(t, "new value", new, "")
}

func TestCalculateNewValueWithEmptyCurrentAndRealNew(t *testing.T) {
	o := NewObject("", "")
	o.newHeaders["Cache-Control"] = "new-cachecontrol"
	current := ""
	change, new := o.calculateNewValue("Cache-Control", &current)
	assertEqualBool(t, "change", change, true)
	assertEqualString(t, "new value", new, "new-cachecontrol")
}

func TestCalculateNewValueWithDifferentCurrentAndNew(t *testing.T) {
	o := NewObject("", "")
	o.newHeaders["Cache-Control"] = "new-cachecontrol"
	current := "current-cachecontrol"
	change, new := o.calculateNewValue("Cache-Control", &current)
	assertEqualBool(t, "change", change, true)
	assertEqualString(t, "new value", new, "new-cachecontrol")
}

func TestCalculateNewValueWithSameCurrentAndNew(t *testing.T) {
	o := NewObject("", "")
	o.newHeaders["Cache-Control"] = "x"
	current := "x"
	change, new := o.calculateNewValue("Cache-Control", &current)
	assertEqualBool(t, "change", change, false)
	assertEqualString(t, "new value", new, "x")
}
