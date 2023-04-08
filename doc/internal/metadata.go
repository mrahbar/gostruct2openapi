package internal

import (
	"fmt"
	"strings"
)

type StructMetadata map[string]string

func (s StructMetadata) append(key string, value string) {
	if _, exists := s[key]; exists {
		s[key] = fmt.Sprintf("%s\n%s", s[key], value)
	} else {
		s[key] = value
	}
}

func (s StructMetadata) Lookup(key string, _default string) string {
	if v, exists := s[key]; exists {
		return v
	} else {
		return _default
	}
}

const (
	metadataToken   = "@"
	DescriptionAttr = "@description"
	TitleAttr       = "@title"
)

type MetadataParser struct {
}

func NewMetadataParser() *MetadataParser {
	return &MetadataParser{}
}

func (o *MetadataParser) ParseStructDesc(desc string) StructMetadata {
	out := make(StructMetadata)
	comments := strings.Split(desc, "\n")

	for line := 0; line < len(comments); line++ {
		commentLine := strings.TrimSpace(comments[line])
		attribute := strings.Split(commentLine, " ")[0]
		value := strings.TrimSpace(commentLine[len(attribute):])

		if strings.HasPrefix(attribute, metadataToken) {
			out.append(attribute, value)
		} else {
			out.append(DescriptionAttr, commentLine)
		}
	}

	return out
}
