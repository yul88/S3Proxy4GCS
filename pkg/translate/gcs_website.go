package translate

import (
	"log/slog"

	"cloud.google.com/go/storage"
)

// TranslateS3ToGCSWebsite maps S3 WebsiteConfiguration to GCS BucketWebsite
func TranslateS3ToGCSWebsite(s3Cfg WebsiteConfiguration) *storage.BucketWebsite {
	slog.Info("Translating S3 Website Configuration")

	gcsWebsite := &storage.BucketWebsite{}

	if s3Cfg.IndexDocument != nil {
		gcsWebsite.MainPageSuffix = s3Cfg.IndexDocument.Suffix
		slog.Debug("Website MainPageSuffix", "suffix", gcsWebsite.MainPageSuffix)
	}

	if s3Cfg.ErrorDocument != nil {
		gcsWebsite.NotFoundPage = s3Cfg.ErrorDocument.Key
		slog.Debug("Website NotFoundPage", "page", gcsWebsite.NotFoundPage)
	}

	return gcsWebsite
}
