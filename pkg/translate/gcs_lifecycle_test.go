package translate

import (
	"encoding/xml"
	"testing"
)

func TestTranslateS3ToGCS(t *testing.T) {
	xmlInput := `
	<LifecycleConfiguration>
		<Rule>
			<ID>TestRule</ID>
			<Status>Enabled</Status>
			<Filter>
				<Prefix>logs/</Prefix>
			</Filter>
			<Transition>
				<Days>30</Days>
				<StorageClass>GLACIER</StorageClass>
			</Transition>
			<Expiration>
				<Days>365</Days>
			</Expiration>
		</Rule>
	</LifecycleConfiguration>
	`

	var s3Cfg LifecycleConfiguration
	if err := xml.Unmarshal([]byte(xmlInput), &s3Cfg); err != nil {
		t.Fatalf("Failed to unmarshal XML input: %v", err)
	}

	gcsJSON, err := TranslateS3ToGCS(&s3Cfg)
	if err != nil {
		t.Fatalf("Failed to translate S3 to GCS: %v", err)
	}

	t.Logf("Generated GCS JSON:\n%s", string(gcsJSON))

	// Simple assertions
	if len(gcsJSON) == 0 {
		t.Error("Generated JSON is empty")
	}

	// Double check some keywords are present
	jsonStr := string(gcsJSON)
	if !contains(jsonStr, `"Delete"`) {
		t.Error(`Expected Action type "Delete"`)
	}
	if !contains(jsonStr, `"SetStorageClass"`) {
		t.Error(`Expected Action type "SetStorageClass"`)
	}
	if !contains(jsonStr, `"COLDLINE"`) {
		t.Error(`Expected GCS StorageClass mapping "COLDLINE" for GLACIER`)
	}
}

func contains(s, substr string) bool {
	// Simple string contains for testing without heavy deps
	return len(s) >= len(substr) && func() bool {
		for i := 0; i < len(s)-len(substr)+1; i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}
func TestTranslateS3ToGCS_PrefixFilter(t *testing.T) {
	xmlInput := `
	<LifecycleConfiguration>
		<Rule>
			<ID>PrefixRule</ID>
			<Status>Enabled</Status>
			<Filter>
				<Prefix>images/</Prefix>
			</Filter>
			<Expiration>
				<Days>90</Days>
			</Expiration>
		</Rule>
	</LifecycleConfiguration>
	`

	var s3Cfg LifecycleConfiguration
	if err := xml.Unmarshal([]byte(xmlInput), &s3Cfg); err != nil {
		t.Fatalf("Failed to unmarshal XML input: %v", err)
	}

	gcsJSON, err := TranslateS3ToGCS(&s3Cfg)
	if err != nil {
		t.Fatalf("Failed to translate S3 to GCS: %v", err)
	}

	jsonStr := string(gcsJSON)
	if !contains(jsonStr, `"matchesPrefix"`) {
		t.Error(`Expected "matchesPrefix" condition`)
	}
	if !contains(jsonStr, `["images/"]`) && !contains(jsonStr, `"images/"`) { // It should be an array of strings
		t.Error(`Expected "images/" inside matchesPrefix array`)
	}
}
