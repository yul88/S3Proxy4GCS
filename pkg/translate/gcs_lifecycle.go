package translate

import (
	"encoding/json"
	"fmt"
)

// GCSLifecycle represents the top-level GCS Lifecycle JSON
type GCSLifecycle struct {
	Rules []GCSLifecycleRule `json:"rule"`
}

// GCSLifecycleRule represents a single rule
type GCSLifecycleRule struct {
	Action    GCSLifecycleAction    `json:"action"`
	Condition GCSLifecycleCondition `json:"condition"`
}

// GCSLifecycleAction represents the action to take (e.g., Delete, SetStorageClass)
type GCSLifecycleAction struct {
	Type         string `json:"type"`                   // Delete, SetStorageClass
	StorageClass string `json:"storageClass,omitempty"` // e.g. NEARLINE, COLDLINE, ARCHIVE
}

// GCSLifecycleCondition represents the conditions when the action should trigger
type GCSLifecycleCondition struct {
	Age                 *int     `json:"age,omitempty"`
	CreatedBefore       *string  `json:"createdBefore,omitempty"` // yyyy-mm-dd
	IsLive              *bool    `json:"isLive,omitempty"`
	MatchesStorageClass []string `json:"matchesStorageClass,omitempty"`
	MatchesPrefix       []string `json:"matchesPrefix,omitempty"`
	MatchesSuffix       []string `json:"matchesSuffix,omitempty"`
	NumNewerVersions    *int     `json:"numNewerVersions,omitempty"`
}

// TranslateS3ToGCS Lifecycle takes an S3 LifecycleConfiguration and returns a GCS JSON byte slice and an error.
func TranslateS3ToGCS(s3Cfg *LifecycleConfiguration) ([]byte, error) {
	var gcsCfg GCSLifecycle

	for _, s3Rule := range s3Cfg.Rules {
		if s3Rule.Status != "Enabled" {
			continue
		}

		// Translate Expirations
		if s3Rule.Expiration != nil {
			rule := GCSLifecycleRule{
				Action: GCSLifecycleAction{Type: "Delete"},
			}
			if err := applyRuleFilter(s3Rule.Filter, &rule.Condition); err != nil {
				return nil, err
			}
			
			if s3Rule.Expiration.Days != nil {
				rule.Condition.Age = s3Rule.Expiration.Days
			}
			if s3Rule.Expiration.Date != nil {
				gcsDate := formatDateS3toGCS(*s3Rule.Expiration.Date)
				rule.Condition.CreatedBefore = &gcsDate
			}
			gcsCfg.Rules = append(gcsCfg.Rules, rule)
		}

		// Translate Transitions
		for _, trans := range s3Rule.Transitions {
			rule := GCSLifecycleRule{
				Action: GCSLifecycleAction{
					Type:         "SetStorageClass",
					StorageClass: mapStorageClass(trans.StorageClass),
				},
			}
			if err := applyRuleFilter(s3Rule.Filter, &rule.Condition); err != nil {
				return nil, err
			}

			if trans.Days != nil {
				rule.Condition.Age = trans.Days
			}
			if trans.Date != nil {
				gcsDate := formatDateS3toGCS(*trans.Date)
				rule.Condition.CreatedBefore = &gcsDate
			}
			gcsCfg.Rules = append(gcsCfg.Rules, rule)
		}

		// Translate Noncurrent Version Expirations
		if s3Rule.NoncurrentVersionExpirations != nil {
			rule := GCSLifecycleRule{
				Action: GCSLifecycleAction{Type: "Delete"},
				Condition: GCSLifecycleCondition{
					NumNewerVersions: s3Rule.NoncurrentVersionExpirations.NoncurrentDays,
					IsLive:           boolPtr(false),
				},
			}
			if err := applyRuleFilter(s3Rule.Filter, &rule.Condition); err != nil {
				return nil, err
			}
			gcsCfg.Rules = append(gcsCfg.Rules, rule)
		}
	}

	return json.MarshalIndent(gcsCfg, "", "  ")
}

func applyRuleFilter(f *Filter, c *GCSLifecycleCondition) error {
	if f == nil {
		return nil
	}

	if f.ObjectSizeGreaterThan != nil || f.ObjectSizeLessThan != nil {
		return fmt.Errorf("ObjectSize filters (ObjectSizeGreaterThan/ObjectSizeLessThan) are not supported by GCS Lifecycle translation")
	}

	if f.Tag != nil {
		return fmt.Errorf("Tag filters are not supported by GCS Lifecycle translation")
	}

	if f.Prefix != nil && *f.Prefix != "" {
		c.MatchesPrefix = []string{*f.Prefix}
	}

	if f.And != nil {
		if f.And.ObjectSizeGreaterThan != nil || f.And.ObjectSizeLessThan != nil {
			return fmt.Errorf("ObjectSize filters in And are not supported by GCS Lifecycle translation")
		}
		if len(f.And.Tags) > 0 {
			return fmt.Errorf("Tag filters in And are not supported by GCS Lifecycle translation")
		}
		if f.And.Prefix != nil && *f.And.Prefix != "" {
			c.MatchesPrefix = append(c.MatchesPrefix, *f.And.Prefix)
		}
	}

	return nil
}

func mapStorageClass(s3Class string) string {
	switch s3Class {
	case "STANDARD_IA":
		return "NEARLINE"
	case "ONEZONE_IA":
		return "NEARLINE"
	case "INTELLIGENT_TIERING":
		return "STANDARD"
	case "GLACIER", "GLACIER_IR":
		return "COLDLINE"
	case "DEEP_ARCHIVE":
		return "ARCHIVE"
	default:
		return "STANDARD"
	}
}

func formatDateS3toGCS(s3Date string) string {
	if len(s3Date) >= 10 {
		return s3Date[:10] // Take yyyy-mm-dd
	}
	return s3Date
}

func boolPtr(b bool) *bool {
	return &b
}
