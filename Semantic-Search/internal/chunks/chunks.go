package chunks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Chunks struct {
	Id            string
	Index         int
	Content       string
	StartPosition int // starting word index, inclusive
	EndPosition   int // ending word index, exclusive
	ChunkEmbed    []float32
	DocId         string
	MetaData      map[string]any
}

type WordChunker struct {
	ChunkSize int
	Overlap   int
}

func NewWordChunker(chunkSize, overlap int) (*WordChunker, error) {
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must be greater than zero")
	}

	if overlap < 0 {
		return nil, errors.New("overlap cannot be negative")
	}

	if overlap >= chunkSize {
		return nil, errors.New("overlap must be smaller than chunk size")
	}

	return &WordChunker{
		ChunkSize: chunkSize,
		Overlap:   overlap,
	}, nil
}

func (wc *WordChunker) Split(ctx context.Context, content string) ([]Chunks, error) {
	words := strings.Fields(content)

	if len(words) == 0 {
		return []Chunks{}, fmt.Errorf("Empty content")
	}

	stepSize := wc.ChunkSize - wc.Overlap

	chunks := make([]Chunks, 0)

	for start := 0; start < len(words); start += stepSize {
		end := start + wc.ChunkSize
		if end > len(words) {
			end = len(words)
		}

		chunks = append(chunks, Chunks{
			Id:            uuid.NewString(),
			Index:         len(chunks),
			Content:       strings.Join(words[start:end], " "),
			StartPosition: start,
			EndPosition:   end,
		})

		if end == len(words) {
			break
		}
	}

	return chunks, nil
}
