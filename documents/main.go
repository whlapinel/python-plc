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
const outDirName = "static"

func main() {
	dirEntries, err := os.ReadDir(docsDirName)
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range dirEntries {
		if strings.Split(file.Name(), ".")[1] == "md" {
			continue
		} else {
			dirEntries = append(dirEntries[:i], dirEntries[i+1:]...)
		}
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open(entry.Name())
		if err != nil {
			log.Fatal(err)
		}
		contents := make([]byte, 1<<20)
		size, err := file.Read(contents)
		log.Println("bytes read: ", size)
		if err != nil {
			log.Fatal(err)
		}
		contents = contents[:size]
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
			panic(err)
		}
		outFilePath := filepath.Join(outDirName, file.Name())
		output, err := os.Create(outFilePath)
		if err != nil {
			log.Println(err)
		}
		log.Println("bytes written:", size)
		Layout(buf.String()).Render(context.Background(), output)
	}

}
