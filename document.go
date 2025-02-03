// Package llm implements a very small subset of the langchain project in Go.
package llm

// Document structure used in LLM applications.
type Document struct {
	PageContent string
	Metadata    map[string]any
	Score       float32
}

// SplitDocuments splits documents using a textsplitter.
func SplitDocuments(textSplitter TextSplitter, documents []Document) ([]Document, error) {
	texts := make([]string, 0)
	metadatas := make([]map[string]any, 0)

	for _, document := range documents {
		texts = append(texts, document.PageContent)
		metadatas = append(metadatas, document.Metadata)
	}

	return CreateDocuments(textSplitter, texts, metadatas)
}

// CreateDocuments creates documents from texts and metadatas with a text splitter. If
// the length of the metadatas is zero, the result documents will contain no metadata.
// Otherwise, the numbers of texts and metadatas must match.
func CreateDocuments(textSplitter TextSplitter, texts []string, metadatas []map[string]any) ([]Document, error) {
	if len(metadatas) == 0 {
		metadatas = make([]map[string]any, len(texts))
	}

	if len(texts) != len(metadatas) {
		return nil, ErrMismatchMetadatasAndText
	}

	documents := make([]Document, 0)

	for i := 0; i < len(texts); i++ {
		chunks, err := textSplitter.SplitText(texts[i])
		if err != nil {
			return nil, err
		}

		for _, chunk := range chunks {
			// Copy the document metadata
			curMetadata := make(map[string]any, len(metadatas[i]))

			for key, value := range metadatas[i] {
				curMetadata[key] = value
			}

			documents = append(documents, Document{
				PageContent: chunk,
				Metadata:    curMetadata,
			})
		}
	}

	return documents, nil
}
