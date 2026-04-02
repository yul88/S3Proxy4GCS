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
)

func TestListObjectVersionsWithAWSSDK(t *testing.T) {
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
		config.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
		config.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = &http.Client{Transport: transport}
		o.BaseEndpoint = aws.String("http://storage.googleapis.com")
	})

	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(getTestBucket()),
	}

	t.Logf("Sending ListObjectVersions via standard AWS S3 SDK Go...")
	_, err = client.ListObjectVersions(context.TODO(), input)
	if err != nil {
		t.Fatalf("Failed to execute ListObjectVersions: %v", err)
	}

	t.Logf("ListObjectVersions succeeded via Proxy (DryRun verified)!")
}

func TestHeadObjectVersionIdWithAWSSDK(t *testing.T) {
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
		config.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
		config.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = &http.Client{Transport: transport}
		o.BaseEndpoint = aws.String("http://storage.googleapis.com")
	})

	// Ensure object exists and capture its version
	putOut, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(getTestBucket()),
		Key:    aws.String("test-object-key"),
		Body:   strings.NewReader("Temporary content for versioning test"),
	})
	if err != nil {
		t.Fatalf("Failed to create object for versioning test: %v", err)
	}

	expectedVersion := ""
	if putOut.VersionId != nil {
		expectedVersion = *putOut.VersionId
		t.Logf("Created object with VersionId: %s", expectedVersion)
	}

	input := &s3.HeadObjectInput{
		Bucket: aws.String(getTestBucket()),
		Key:    aws.String("test-object-key"),
	}

	t.Logf("Sending HeadObject via standard AWS S3 SDK Go...")
	resp, err := client.HeadObject(context.TODO(), input)
	if err != nil {
		t.Fatalf("Failed to execute HeadObject: %v", err)
	}

	if resp.VersionId == nil {
		t.Fatalf("Expected VersionId in HeadObject response, got nil")
	}

	gotVersion := *resp.VersionId
	if expectedVersion != "" && gotVersion != expectedVersion {
		t.Fatalf("Expected VersionId to be '%s', got '%s'", expectedVersion, gotVersion)
	}

	t.Logf("HeadObject returned VersionId: %s (DryRun verified)!", gotVersion)
}
