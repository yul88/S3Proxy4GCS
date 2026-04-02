package translate

import (
	"encoding/xml"
	"testing"
)

func TestTranslateS3ToGCSWebsite(t *testing.T) {
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<WebsiteConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <IndexDocument>
        <Suffix>index.html</Suffix>
    </IndexDocument>
    <ErrorDocument>
        <Key>error.html</Key>
    </ErrorDocument>
</WebsiteConfiguration>`

	var s3Cfg WebsiteConfiguration
	err := xml.Unmarshal([]byte(xmlInput), &s3Cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	gcsWebsite := TranslateS3ToGCSWebsite(s3Cfg)

	if gcsWebsite == nil {
		t.Fatalf("Expected non-nil GCS Website settings")
	}

	if gcsWebsite.MainPageSuffix != "index.html" {
		t.Errorf("Expected MainPageSuffix 'index.html', got '%s'", gcsWebsite.MainPageSuffix)
	}

	if gcsWebsite.NotFoundPage != "error.html" {
		t.Errorf("Expected NotFoundPage 'error.html', got '%s'", gcsWebsite.NotFoundPage)
	}
}

func TestTranslateS3ToGCSWebsitePartial(t *testing.T) {
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<WebsiteConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <IndexDocument>
        <Suffix>index.html</Suffix>
    </IndexDocument>
</WebsiteConfiguration>`

	var s3Cfg WebsiteConfiguration
	err := xml.Unmarshal([]byte(xmlInput), &s3Cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	gcsWebsite := TranslateS3ToGCSWebsite(s3Cfg)

	if gcsWebsite == nil {
		t.Fatalf("Expected non-nil GCS Website settings")
	}

	if gcsWebsite.MainPageSuffix != "index.html" {
		t.Errorf("Expected MainPageSuffix 'index.html', got '%s'", gcsWebsite.MainPageSuffix)
	}

	if gcsWebsite.NotFoundPage != "" {
		t.Errorf("Expected empty NotFoundPage, got '%s'", gcsWebsite.NotFoundPage)
	}
}
