package translate

import "encoding/xml"

// Tagging represents the top-level S3 XML tag for object tagging
type Tagging struct {
	XMLName xml.Name `xml:"Tagging"`
	TagSet  []Tag    `xml:"TagSet>Tag"`
}
