package translate

import (
	"log/slog"

	"cloud.google.com/go/storage"
)

// TranslateS3ToGCSLogging maps S3 BucketLoggingStatus to GCS BucketLogging
func TranslateS3ToGCSLogging(s3Cfg BucketLoggingStatus) *storage.BucketLogging {
	if s3Cfg.LoggingEnabled == nil {
		return nil // Disabled or unset
	}

	slog.Info("Translating S3 Logging", "targetBucket", s3Cfg.LoggingEnabled.TargetBucket, "targetPrefix", s3Cfg.LoggingEnabled.TargetPrefix)

	return &storage.BucketLogging{
		LogBucket:       s3Cfg.LoggingEnabled.TargetBucket,
		LogObjectPrefix: s3Cfg.LoggingEnabled.TargetPrefix,
	}
}

// TranslateGCSToS3Logging converts GCS BucketLogging to S3 BucketLoggingStatus XML
func TranslateGCSToS3Logging(gcsLogging *storage.BucketLogging) *BucketLoggingStatus {
	if gcsLogging == nil || gcsLogging.LogBucket == "" {
		return &BucketLoggingStatus{}
	}

	return &BucketLoggingStatus{
		LoggingEnabled: &LoggingEnabled{
			TargetBucket: gcsLogging.LogBucket,
			TargetPrefix: gcsLogging.LogObjectPrefix,
		},
	}
}
