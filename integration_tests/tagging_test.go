package integration

import (
	"context"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestPutObjectTaggingWithAWSSDK(t *testing.T) {
	dialer := &net.Dialer{}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if strings.HasPrefix(addr, "storage.googleapis.com") {
				return dialer.DialContext(ctx, network, "localhost:8081")
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}

	creds := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     getAWSAccessKey(),
			SecretAccessKey: getAWSSecretKey(),
			Source:          "test-env",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = &http.Client{Transport: transport}
		o.BaseEndpoint = aws.String("http://storage.googleapis.com")
	})

	bucketName := getTestBucket()
	objectKey := getTestPrefix() + "test-object-key"

	// Ensure object exists
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   strings.NewReader("Temporary content for tagging test"),
	})
	if err != nil {
		t.Fatalf("Failed to create object for tagging test: %v", err)
	}

	input := &s3.PutObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Tagging: &types.Tagging{
			TagSet: []types.Tag{
				{
					Key:   aws.String("Project"),
					Value: aws.String("Vis"),
				},
				{
					Key:   aws.String("Stage"),
					Value: aws.String("Archived"),
				},
			},
		},
	}

	t.Logf("Sending PutObjectTagging via standard AWS S3 SDK Go...")
	_, err = client.PutObjectTagging(context.TODO(), input)
	if err != nil {
		t.Fatalf("Failed to execute PutObjectTagging: %v", err)
	}

	t.Logf("Sending GetObjectTagging via standard AWS S3 SDK Go...")
	getOut, err := client.GetObjectTagging(context.TODO(), &s3.GetObjectTaggingInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		t.Fatalf("Failed to execute GetObjectTagging: %v", err)
	}
	if len(getOut.TagSet) == 0 {
		t.Fatalf("GetObjectTagging returned 0 tags, expected at least 1")
	}
	t.Logf("GetObjectTagging succeeded (Retrieved %d tags)!", len(getOut.TagSet))
}
