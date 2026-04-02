package translate

import (
	"encoding/xml"
	"testing"
)

func TestTranslateS3ToGCSLogging(t *testing.T) {
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus xmlns="http://doc.s3.amazonaws.com/2006-03-01/">
  <LoggingEnabled>
    <TargetBucket>my-log-bucket</TargetBucket>
    <TargetPrefix>my-log-prefix/</TargetPrefix>
  </LoggingEnabled>
</BucketLoggingStatus>`

	var s3Cfg BucketLoggingStatus
	err := xml.Unmarshal([]byte(xmlInput), &s3Cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	gcsLogging := TranslateS3ToGCSLogging(s3Cfg)

	if gcsLogging == nil {
		t.Fatalf("Expected non-nil GCS Logging")
	}

	if gcsLogging.LogBucket != "my-log-bucket" {
		t.Errorf("Expected LogBucket 'my-log-bucket', got '%s'", gcsLogging.LogBucket)
	}

	if gcsLogging.LogObjectPrefix != "my-log-prefix/" {
		t.Errorf("Expected LogObjectPrefix 'my-log-prefix/', got '%s'", gcsLogging.LogObjectPrefix)
	}
}

func TestTranslateS3ToGCSLoggingDisabled(t *testing.T) {
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus xmlns="http://doc.s3.amazonaws.com/2006-03-01/">
</BucketLoggingStatus>`

	var s3Cfg BucketLoggingStatus
	err := xml.Unmarshal([]byte(xmlInput), &s3Cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	gcsLogging := TranslateS3ToGCSLogging(s3Cfg)

	if gcsLogging != nil {
		t.Errorf("Expected nil GCS Logging for disabled status")
	}
}
