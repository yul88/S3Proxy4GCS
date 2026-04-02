package translate

import (
	"log/slog"
	"time"

	"cloud.google.com/go/storage"
)

// TranslateS3ToGCSCors converts S3 CORS Configuration XML to GCS storage.CORS slice
func TranslateS3ToGCSCors(s3Cfg *CORSConfiguration) []storage.CORS {
	var gcsCors []storage.CORS

	for _, rule := range s3Cfg.CORSRules {
		var maxAge time.Duration
		if rule.MaxAgeSeconds != nil {
			maxAge = time.Duration(*rule.MaxAgeSeconds) * time.Second
		}

		if len(rule.AllowedHeaders) > 0 {
			slog.Warn("S3 AllowedHeaders (Request Headers) are not natively supported by GCS CORS translation and will be ignored.")
		}

		gcsRule := storage.CORS{
			MaxAge:          maxAge,
			Methods:         rule.AllowedMethods,
			Origins:         rule.AllowedOrigins,
			ResponseHeaders: rule.ExposeHeaders,
		}

		gcsCors = append(gcsCors, gcsRule)
	}

	return gcsCors
}

// TranslateGCSToS3Cors converts GCS CORS configuration to S3 CORSConfiguration XML
func TranslateGCSToS3Cors(gcsCors []storage.CORS) *CORSConfiguration {
	if len(gcsCors) == 0 {
		return nil
	}

	s3Cfg := &CORSConfiguration{}
	for _, rule := range gcsCors {
		var maxAge *int
		if rule.MaxAge > 0 {
			seconds := int(rule.MaxAge.Seconds())
			maxAge = &seconds
		}

		s3Rule := CORSRule{
			AllowedMethods: rule.Methods,
			AllowedOrigins: rule.Origins,
			ExposeHeaders:  rule.ResponseHeaders,
			MaxAgeSeconds:  maxAge,
		}
		s3Cfg.CORSRules = append(s3Cfg.CORSRules, s3Rule)
	}

	return s3Cfg
}
