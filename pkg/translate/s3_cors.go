package translate

import "encoding/xml"

// CORSConfiguration represents the top-level S3 XML tag
type CORSConfiguration struct {
	XMLName   xml.Name   `xml:"CORSConfiguration"`
	CORSRules []CORSRule `xml:"CORSRule"`
}

// CORSRule represents a single CORS rule in S3
type CORSRule struct {
	ID             *string  `xml:"ID,omitempty"`
	AllowedHeaders []string `xml:"AllowedHeader,omitempty"`
	AllowedMethods []string `xml:"AllowedMethod"`
	AllowedOrigins []string `xml:"AllowedOrigin"`
	ExposeHeaders  []string `xml:"ExposeHeader,omitempty"`
	MaxAgeSeconds  *int     `xml:"MaxAgeSeconds,omitempty"`
}
