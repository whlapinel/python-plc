package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const docsDirName = "documents"
const outDirName = "static/documents"

func main() {
	dirEntries, err := os.ReadDir(docsDirName)
	if err != nil {
		log.Fatal(err)
	}

	// Filter out non-Markdown files
	var mdFiles []os.DirEntry
	for _, entry := range dirEntries {
		if filepath.Ext(entry.Name()) == ".md" {
			mdFiles = append(mdFiles, entry)
		}
	}

	for _, entry := range mdFiles {
		if entry.IsDir() {
			continue
		}

		inputPath := filepath.Join(docsDirName, entry.Name())

		// Read entire file at once
		contents, err := os.ReadFile(inputPath)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", inputPath, err)
		}

		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		)

		var buf bytes.Buffer
		if err := md.Convert(contents, &buf); err != nil {
			log.Fatalf("Failed to convert Markdown: %v", err)
		}

		outputPath := filepath.Join(outDirName, strings.TrimSuffix(entry.Name(), ".md")+".html")
		output, err := os.Create(outputPath)
		if err != nil {
			log.Fatalf("Failed to create output file %s: %v", outputPath, err)
		}
		defer output.Close()

		log.Println("Writing:", outputPath)
		Layout(buf.String()).Render(context.Background(), output)
	}
}
