package charsplitter

import (
	"strings"

	"github.com/peterhellberg/llm"
)

var _ llm.TextSplitter = Splitter{}

// Splitter is a text splitter that will split texts recursively by different characters.
type Splitter struct {
	Separators    []string
	ChunkSize     int
	ChunkOverlap  int
	LenFunc       func(string) int
	KeepSeparator bool
}

// New creates a new recursive character splitter with default values. By default,
// the separators used are "\n\n", "\n", " " and "". The chunk size is set to 512
// and chunk overlap is set to 100.
func New(opts ...Option) Splitter {
	options := defaultOptions()

	for _, o := range opts {
		o(&options)
	}

	s := Splitter{
		Separators:    options.Separators,
		ChunkSize:     options.ChunkSize,
		ChunkOverlap:  options.ChunkOverlap,
		LenFunc:       options.LenFunc,
		KeepSeparator: options.KeepSeparator,
	}

	return s
}

// SplitText splits a text into multiple text.
func (s Splitter) SplitText(text string) ([]string, error) {
	return s.splitText(text, s.Separators)
}

// addSeparatorInSplits adds the separator in each of splits.
func (s Splitter) addSeparatorInSplits(splits []string, separator string) []string {
	splitsWithSeparator := make([]string, 0, len(splits))
	for i, s := range splits {
		if i > 0 {
			s = separator + s
		}
		splitsWithSeparator = append(splitsWithSeparator, s)
	}
	return splitsWithSeparator
}

func (s Splitter) splitText(text string, separators []string) ([]string, error) {
	finalChunks := make([]string, 0)

	// Find the appropriate separator.
	separator := separators[len(separators)-1]

	newSeparators := []string{}

	for i, c := range separators {
		if c == "" || strings.Contains(text, c) {
			separator = c
			newSeparators = separators[i+1:]

			break
		}
	}

	splits := strings.Split(text, separator)

	if s.KeepSeparator {
		splits = s.addSeparatorInSplits(splits, separator)
		separator = ""
	}

	goodSplits := make([]string, 0)

	// Merge the splits, recursively splitting larger texts.
	for _, split := range splits {
		if s.LenFunc(split) < s.ChunkSize {
			goodSplits = append(goodSplits, split)

			continue
		}

		if len(goodSplits) > 0 {
			mergedText := mergeSplits(goodSplits, separator, s.ChunkSize, s.ChunkOverlap, s.LenFunc)

			finalChunks = append(finalChunks, mergedText...)
			goodSplits = make([]string, 0)
		}

		if len(newSeparators) == 0 {
			finalChunks = append(finalChunks, split)
		} else {
			otherInfo, err := s.splitText(split, newSeparators)
			if err != nil {
				return nil, err
			}

			finalChunks = append(finalChunks, otherInfo...)
		}
	}

	if len(goodSplits) > 0 {
		mergedText := mergeSplits(goodSplits, separator, s.ChunkSize, s.ChunkOverlap, s.LenFunc)
		finalChunks = append(finalChunks, mergedText...)
	}

	return finalChunks, nil
}

// joinDocs comines two documents with the separator used to split them.
func joinDocs(docs []string, separator string) string {
	return strings.TrimSpace(strings.Join(docs, separator))
}

// mergeSplits merges smaller splits into splits that are closer to the chunkSize.
func mergeSplits(splits []string, separator string, chunkSize int, chunkOverlap int, lenFunc func(string) int) []string {
	var (
		docs       = make([]string, 0)
		currentDoc = make([]string, 0)
		total      = 0
	)

	for _, split := range splits {
		totalWithSplit := total + lenFunc(split)

		if len(currentDoc) != 0 {
			totalWithSplit += lenFunc(separator)
		}

		if totalWithSplit > chunkSize && len(currentDoc) > 0 {
			doc := joinDocs(currentDoc, separator)
			if doc != "" {
				docs = append(docs, doc)
			}

			for shouldPop(chunkOverlap, chunkSize, total, lenFunc(split), lenFunc(separator), len(currentDoc)) {
				total -= lenFunc(currentDoc[0])
				if len(currentDoc) > 1 {
					total -= lenFunc(separator)
				}

				currentDoc = currentDoc[1:]
			}
		}

		currentDoc = append(currentDoc, split)

		total += lenFunc(split)

		if len(currentDoc) > 1 {
			total += lenFunc(separator)
		}
	}

	doc := joinDocs(currentDoc, separator)

	if doc != "" {
		docs = append(docs, doc)
	}

	return docs
}

// Keep popping if:
//   - the chunk is larger than the chunk overlap
//   - or if there are any chunks and the length is long
func shouldPop(chunkOverlap, chunkSize, total, splitLen, separatorLen, currentDocLen int) bool {
	docsNeededToAddSep := 2

	if currentDocLen < docsNeededToAddSep {
		separatorLen = 0
	}

	return currentDocLen > 0 && (total > chunkOverlap || (total+splitLen+separatorLen > chunkSize && total > 0))
}
