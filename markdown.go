package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getFileString(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	bs := make([]byte, info.Size())
	_, err = f.Read(bs)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// SnippetLines represents the range of lines a code snippet extends
type SnippetLines struct {
	Start int
	End   int
}

func getSnippetsFromMarkdownFile(path, language string) ([]*Snippet, error) {
	filename := strings.Split(path, ".")[0]
	fileString, err := getFileString(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(fileString, "\n")

	snippetLinesList := []SnippetLines{}
	start := 0
	end := 0
	counting := false

	for index, line := range lines {
		lineNumber := index + 1
		if line == "```javascript" {
			start = lineNumber
			counting = true
		}

		if counting && line == "```" {
			end = lineNumber
			snippetLinesList = append(snippetLinesList, SnippetLines{Start: start, End: end})
			counting = false
		}
	}

	languageSeparator := fmt.Sprintf("```%v", language)
	splitFile := strings.Split(fileString, languageSeparator)[1:]

	javascriptSnippets := []*Snippet{}
	for index, substr := range splitFile {
		snippet := &Snippet{
			Language: language,
			Content:  strings.Split(substr, "```")[0],
			Filename: filename,
			Start:    snippetLinesList[index].Start,
			End:      snippetLinesList[index].End,
		}
		javascriptSnippets = append(javascriptSnippets, snippet)
	}

	return javascriptSnippets, nil
}

func replaceSnippetsWithURLs(path string, gistSnippets []*GistSnippet) error {
	fileString, err := getFileString(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(fileString, "\n")

	for _, gs := range gistSnippets {
		lines[gs.Snippet.Start-1] = gs.Gist.GetHTMLURL()
	}

	newLines := []string{}
	for index, gs := range gistSnippets {
		startIndex := gs.Snippet.Start
		if index == 0 {
			newLines = append(newLines, lines[:startIndex]...)
		} else {
			previousGistSnippetEndIndex := gistSnippets[index-1].Snippet.End
			newLines = append(newLines, lines[previousGistSnippetEndIndex:startIndex]...)
		}
	}

	file, err := os.Create("react-testing-medium.md")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range newLines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
