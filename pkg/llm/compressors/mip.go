package compressors

import (
	"github.com/greenboxal/aip/pkg/llm"
	"github.com/greenboxal/aip/pkg/llm/documents"
	"github.com/greenboxal/aip/pkg/llm/tokenizers"
)

type MipDocumentLevel struct {
	Level        int
	TokenCount   int
	DocumentLink documents.DocumentPointer
}

type MipManifest struct {
	documents.DocumentBase[*MipManifest]

	Levels []MipDocumentLevel
}

type MipCompressor struct {
	BaseCompressor DocumentCompressor

	Tokenizer tokenizers.BasicTokenizer

	MaxMipLevels int

	MinOutputTokens int
	MaxOutputTokens int
}

func (mc *MipCompressor) MaxTokens() int {
	return mc.MaxMipLevels
}

func (mc *MipCompressor) Compress(ctx llm.ChainContext, document documents.Document, options ...CompressionOption) (documents.Document, error) {
	opts := NewCompressionOptions(options...)

	mctx := &mipCompressorContext{
		MipCompressor: mc,
		ctx:           ctx,
		options:       opts,

		inputDocument:  document,
		outputDocument: &MipManifest{},
	}

	return mctx.Reduce()
}

type mipCompressorContext struct {
	*MipCompressor

	ctx     llm.ChainContext
	options CompressionOptions

	currentDocument   documents.Document
	currentTokenCount int
	currentDepth      int

	inputDocument  documents.Document
	outputDocument *MipManifest
}

func (ctx *mipCompressorContext) shouldReduce() bool {
	return ctx.currentDocument != nil && ctx.currentTokenCount > ctx.MinOutputTokens && ctx.currentDepth < ctx.MaxMipLevels
}

func (ctx *mipCompressorContext) setNextAndAppend(document documents.Document) (bool, error) {
	count, err := ctx.Tokenizer.Count(ctx.currentDocument.AsText())

	if err != nil {
		return false, err
	}

	reduced, err := ctx.ctx.Documents().Put(ctx.ctx.Context(), document)

	if err != nil {
		return false, err
	}

	if len(ctx.outputDocument.Levels) > 0 {
		lastLevel := &ctx.outputDocument.Levels[len(ctx.outputDocument.Levels)-1]

		if count > lastLevel.TokenCount {
			return false, nil
		}
	}

	ctx.outputDocument.Levels = append(ctx.outputDocument.Levels, MipDocumentLevel{
		Level:        ctx.currentDepth,
		DocumentLink: reduced,
	})

	ctx.currentDocument = document
	ctx.currentTokenCount = count

	return true, nil
}

func (ctx *mipCompressorContext) reduceRound() (bool, error) {
	ctx.currentDepth++

	targetDocument := ctx.currentDocument
	targetTokenCount := ctx.currentTokenCount / 2

	if ctx.currentTokenCount > ctx.BaseCompressor.MaxTokens() {
		return false, nil
	}

	if targetTokenCount < ctx.MinOutputTokens {
		targetTokenCount = ctx.MinOutputTokens
	}

	if targetTokenCount > ctx.MaxOutputTokens {
		targetTokenCount = ctx.MaxOutputTokens
	}

	result, err := ctx.BaseCompressor.CompressDocument(
		ctx.ctx.Context(),
		targetDocument,
		WithMaxTokens(targetTokenCount),
	)

	if err != nil {
		return false, err
	}

	return ctx.setNextAndAppend(result)
}

func (ctx *mipCompressorContext) Reduce() (documents.Document, error) {
	ctx.currentDepth = 0

	if ok, err := ctx.setNextAndAppend(ctx.inputDocument); err != nil {
		return nil, err
	} else if !ok {
		panic("shouldn't happen")
	}

	for ctx.shouldReduce() {
		shouldContinue, err := ctx.reduceRound()

		if err != nil {
			return nil, err
		}

		if !shouldContinue {
			break
		}
	}

	return ctx.outputDocument, nil
}
