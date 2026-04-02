package translate

import (
	"encoding/xml"
	"testing"
)

func TestTranslateS3ToGCSTagging(t *testing.T) {
	xmlInput := `<?xml version="1.0" encoding="UTF-8"?>
<Tagging xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <TagSet>
        <Tag>
            <Key>Project</Key>
            <Value>Foo</Value>
        </Tag>
        <Tag>
            <Key>Environment</Key>
            <Value>Dev</Value>
        </Tag>
    </TagSet>
</Tagging>`

	var s3Cfg Tagging
	err := xml.Unmarshal([]byte(xmlInput), &s3Cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	existingMeta := map[string]string{
		"s3tag-OldKey": "OldValue",
		"other-meta":   "Preserved",
	}

	updateMeta := TranslateS3ToGCSTagging(s3Cfg, existingMeta)

	if updateMeta["s3tag-Project"] != "Foo" {
		t.Errorf("Expected s3tag-Project to be 'Foo', got '%s'", updateMeta["s3tag-Project"])
	}

	if updateMeta["s3tag-Environment"] != "Dev" {
		t.Errorf("Expected s3tag-Environment to be 'Dev', got '%s'", updateMeta["s3tag-Environment"])
	}

	if updateMeta["s3tag-OldKey"] != "" {
		t.Errorf("Expected s3tag-OldKey to be cleared (\"\"), got '%s'", updateMeta["s3tag-OldKey"])
	}

	if val, ok := updateMeta["other-meta"]; ok {
		t.Errorf("Expected other-meta to be omitted from update (unchanged), got '%s'", val)
	}
}
