package main

import (
	"bufio"
	"bytes"
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

func getSnippetsFromMarkdownFile(file *File, language string) ([]*Snippet, error) {
	lines := strings.Split(file.Contents, "\n")

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
	splitFile := strings.Split(file.Contents, languageSeparator)[1:]

	javascriptSnippets := []*Snippet{}
	for index, substr := range splitFile {
		content := strings.Split(substr, "```")[0]

		content = removeInitialNewline(content)
		content = replaceTabsWithSpaces(content)

		snippet := &Snippet{
			Language: language,
			Content:  content,
			Filename: file.Name,
			Start:    snippetLinesList[index].Start,
			End:      snippetLinesList[index].End,
		}
		javascriptSnippets = append(javascriptSnippets, snippet)
	}

	return javascriptSnippets, nil
}

// UNICODE characters
const (
	TAB     = 9
	SPACE   = 32
	NEWLINE = 10
)

func removeInitialNewline(content string) string {
	if content[0] == NEWLINE {
		content = content[1:]
	}
	return content
}

func replaceTabsWithSpaces(content string) string {
	var buffer bytes.Buffer
	for _, c := range content {
		if c == TAB {
			buffer.WriteRune(SPACE)
			buffer.WriteRune(SPACE)
		} else {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}

func replaceSnippetsWithURLs(file *File, gistSnippets []*GistSnippet) error {
	lines := strings.Split(file.Contents, "\n")

	for _, gs := range gistSnippets {
		lines[gs.Snippet.Start-1] = gs.Gist.GetHTMLURL()
	}

	newLines := []string{}
	for index, gs := range gistSnippets {
		startIndex := gs.Snippet.Start
		if index == 0 {
			// append lines all the way up to the start index of the first snippet
			newLines = append(newLines, lines[:startIndex]...)

		} else {
			previousGistSnippetEndIndex := gistSnippets[index-1].Snippet.End
			newLines = append(newLines, lines[previousGistSnippetEndIndex:startIndex]...)
		}
	}

	// append the rest of the lines
	newLines = append(newLines, lines[gistSnippets[len(gistSnippets)-1].Snippet.End:]...)

	mediumPath := strings.Replace(file.Name, ".md", ".medium.md", 1)

	f, err := os.Create(mediumPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range newLines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
