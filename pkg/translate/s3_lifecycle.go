package translate

import (
	"encoding/xml"
)

// LifecycleConfiguration represents the top-level S3 XML tag
type LifecycleConfiguration struct {
	XMLName xml.Name `xml:"LifecycleConfiguration"`
	Rules   []Rule   `xml:"Rule"`
}

// Rule represents a single lifecycle rule
type Rule struct {
	ID                             string                          `xml:"ID,omitempty"`
	Priority                       int                             `xml:"Priority,omitempty"`
	Filter                         *Filter                         `xml:"Filter,omitempty"`
	Status                         string                          `xml:"Status"` // Enabled or Disabled
	Transitions                    []Transition                    `xml:"Transition,omitempty"`
	Expiration                     *Expiration                     `xml:"Expiration,omitempty"`
	NoncurrentVersionTransitions   []NoncurrentVersionTransition   `xml:"NoncurrentVersionTransition,omitempty"`
	NoncurrentVersionExpirations   *NoncurrentVersionExpiration     `xml:"NoncurrentVersionExpiration,omitempty"`
	AbortIncompleteMultipartUpload *AbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload,omitempty"`
}

// Filter can define a prefix, a tag, or a combination of them via logical operators
type Filter struct {
	Prefix                  *string                  `xml:"Prefix,omitempty"`
	Tag                     *Tag                     `xml:"Tag,omitempty"`
	And                     *AndOperator             `xml:"And,omitempty"`
	ObjectSizeGreaterThan   *int64                   `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan      *int64                   `xml:"ObjectSizeLessThan,omitempty"`
}

// AndOperator represents the <And> logical tag
type AndOperator struct {
	Prefix                *string `xml:"Prefix,omitempty"`
	Tags                  []Tag   `xml:"Tag,omitempty"`
	ObjectSizeGreaterThan *int64  `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64  `xml:"ObjectSizeLessThan,omitempty"`
}

// Tag represents an object tag
type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

// Transition represents a storage class migration
type Transition struct {
	Days         *int    `xml:"Days,omitempty"`
	Date         *string `xml:"Date,omitempty"` // yyyy-mm-ddThh:mm:ss.000Z
	StorageClass string  `xml:"StorageClass"`
}

// Expiration represents object deletion conditions
type Expiration struct {
	Days                       *int    `xml:"Days,omitempty"`
	Date                       *string `xml:"Date,omitempty"`
	ExpiredObjectDeleteMarker *bool   `xml:"ExpiredObjectDeleteMarker,omitempty"`
}

// NoncurrentVersionTransition represents versioning transitions
type NoncurrentVersionTransition struct {
	NoncurrentDays *int   `xml:"NoncurrentDays,omitempty"`
	StorageClass   string `xml:"StorageClass"`
}

// NoncurrentVersionExpiration represents versioning deletions
type NoncurrentVersionExpiration struct {
	NoncurrentDays *int `xml:"NoncurrentDays,omitempty"`
}

// AbortIncompleteMultipartUpload represents cleanup for aborted uploads
type AbortIncompleteMultipartUpload struct {
	DaysAfterInitiation *int `xml:"DaysAfterInitiation,omitempty"`
}
