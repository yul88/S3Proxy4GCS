package translate

import "encoding/xml"

// BucketLoggingStatus represents the top-level S3 XML tag
type BucketLoggingStatus struct {
	XMLName        xml.Name        `xml:"BucketLoggingStatus"`
	LoggingEnabled *LoggingEnabled `xml:"LoggingEnabled,omitempty"`
}

// LoggingEnabled represents the active logging configuration
type LoggingEnabled struct {
	TargetBucket string `xml:"TargetBucket"`
	TargetPrefix string `xml:"TargetPrefix"`
	// TargetGrants is omitted/ignored as GCS uses IAM for log delivery
}
