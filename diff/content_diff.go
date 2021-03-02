package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ContentDiff is the diff between two content objects each containing a media type object: https://swagger.io/specification/#media-type-object
type ContentDiff struct {
	MediaTypeAdded   bool            `json:"mediaTypeAdded,omitempty"`
	MediaTypeDeleted bool            `json:"mediaTypeDeleted,omitempty"`
	MediaTypeDiff    bool            `json:"mediaType,omitempty"`
	ExtensionProps   *ExtensionsDiff `json:"extensions,omitempty"`
	SchemaDiff       *SchemaDiff     `json:"schema,omitempty"`
	ExampleDiff      *ValueDiff      `json:"example,omitempty"`
	// Examples
	EncodingsDiff *EncodingsDiff `json:"encoding,omitempty"`
}

func newContentDiff() ContentDiff {
	return ContentDiff{}
}

func (contentDiff ContentDiff) empty() bool {
	return contentDiff == newContentDiff()
}

func getContentDiff(config *Config, content1, content2 openapi3.Content) ContentDiff {

	result := newContentDiff()

	if len(content1) == 0 && len(content2) == 0 {
		return result
	}

	if len(content1) == 0 && len(content2) != 0 {
		result.MediaTypeAdded = true
		return result
	}

	if len(content1) != 0 && len(content2) == 0 {
		result.MediaTypeDeleted = true
		return result
	}

	mediaType1, mediaTypeValue1, err := getMediaType(content1)
	if err != nil {
		return result
	}

	mediaType2, mediaTypeValue2, err := getMediaType(content2)
	if err != nil {
		return result
	}

	if mediaType1 != mediaType2 {
		result.MediaTypeDiff = true
		return result
	}

	if diff := getExtensionsDiff(mediaTypeValue1.ExtensionProps, mediaTypeValue2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	if diff := getSchemaDiff(config, mediaTypeValue1.Schema, mediaTypeValue2.Schema); !diff.empty() {
		result.SchemaDiff = &diff
	}

	if config.Examples {
		result.ExampleDiff = getValueDiff(mediaTypeValue1.Example, mediaTypeValue1.Example)
	}

	if diff := getEncodingsDiff(config, mediaTypeValue1.Encoding, mediaTypeValue2.Encoding); !diff.empty() {
		result.EncodingsDiff = diff
	}

	return result
}

// getMediaType returns the single MediaType entry in the content map
func getMediaType(content openapi3.Content) (string, *openapi3.MediaType, error) {

	var mediaType string
	var mediaTypeValue *openapi3.MediaType

	if len(content) != 1 {
		return "", nil, fmt.Errorf("content map has more than one value - this shouldn't happen. %+v", content)
	}

	for k, v := range content {
		mediaType = k
		mediaTypeValue = v
	}

	return mediaType, mediaTypeValue, nil
}
