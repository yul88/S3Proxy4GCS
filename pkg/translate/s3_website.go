package translate

import "encoding/xml"

// WebsiteConfiguration represents the top-level S3 XML tag
type WebsiteConfiguration struct {
	XMLName        xml.Name        `xml:"WebsiteConfiguration"`
	IndexDocument  *IndexDocument  `xml:"IndexDocument,omitempty"`
	ErrorDocument  *ErrorDocument  `xml:"ErrorDocument,omitempty"`
	// RoutingRules is omitted/ignored for simplicity as GCS doesn't support it natively in the same way
}

// IndexDocument represents the main page suffix
type IndexDocument struct {
	Suffix string `xml:"Suffix"`
}

// ErrorDocument represents the not found page
type ErrorDocument struct {
	Key string `xml:"Key"`
}
