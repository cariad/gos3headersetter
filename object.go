package gos3headersetter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Object describes an S3 object. Use NewObject() to create new instances.
type Object struct {
	Bucket string
	Key    string

	newHeaders map[string]string
}

// NewObject returns a instance of Object.
func NewObject(bucket string, key string) Object {
	return Object{
		Bucket:     bucket,
		Key:        key,
		newHeaders: make(map[string]string),
	}
}

func (o Object) String() string {
	return fmt.Sprintf("s3://%s/%s", o.Bucket, o.Key)
}

// Apply applies the rules to this S3 object.
func (o Object) Apply(rules []Rule) error {
	for _, rule := range rules {
		o.queueRuleEffect(rule)
	}

	client, err := makeClient()
	if err != nil {
		return err
	}

	in := o.makeHeadObjectInput()
	head, err := client.HeadObject(in)
	if err != nil {
		return err
	}

	copy := o.makeCopyObjectInput(head.Metadata)
	hasChanges := o.updateCopyObjectInput(head, copy)
	if !hasChanges {
		o.log("no changes to apply")
		return nil
	}

	o.log("Applying changes...")
	_, err = client.CopyObject(copy)
	return err
}

func (o Object) makeHeadObjectInput() *s3.HeadObjectInput {
	return &s3.HeadObjectInput{
		Bucket: &o.Bucket,
		Key:    &o.Key,
	}
}

func (o Object) makeCopyObjectInput(metadata map[string]*string) *s3.CopyObjectInput {
	return &s3.CopyObjectInput{
		Bucket:     &o.Bucket,
		Key:        &o.Key,
		CopySource: aws.String(o.Bucket + "/" + o.Key),

		// Perform a fake replacement of the metadata, otherwise AWS will
		// reject the copy because nothing has changed. It doesn't notice
		// that the "ContentType" and/or "CacheControl" have changed.
		Metadata:          metadata,
		MetadataDirective: aws.String("REPLACE"),
	}
}

func (o Object) log(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	info("%v: %s", o, msg)
}

func (o Object) calculateNewValue(header string, currentPtr *string) (bool, string) {
	new := o.newHeaders[header]

	current := ""

	if currentPtr != nil {
		current = *currentPtr
	}

	if new == "" {
		o.log("the rules do not specify a change for %s", header)
		return false, current
	} else if new == current {
		o.log("%s will remain \"%s\"", header, current)
		return false, current
	}

	if current == "" {
		o.log("%s will be set to \"%s\"", header, new)
	} else {
		o.log("%s will be updated from \"%s\" to \"%s\"", header, current, new)
	}
	return true, new
}

func (o Object) updateCopyObjectInput(head *s3.HeadObjectOutput, in *s3.CopyObjectInput) bool {
	ccChange, ccNew := o.calculateNewValue("Cache-Control", head.CacheControl)
	in.CacheControl = &ccNew

	ctChange, ctNew := o.calculateNewValue("Content-Type", head.ContentType)
	in.ContentType = &ctNew

	return ccChange || ctChange
}

func (o *Object) queueRuleEffect(rule Rule) {
	for _, when := range rule.When {
		if endsWith(o.Key, when.Extension) {
			o.newHeaders[rule.Header] = when.Then
			return
		}
	}
	o.newHeaders[rule.Header] = rule.Else
}
