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
	Language string
	Start    int
	End      int
}

func getLanguage(line string) string {

	languageDelimiter := strings.Split(line, "```")
	language := ""
	if len(languageDelimiter) > 1 {
		language = languageDelimiter[1]
	}
	l := Languages[language]
	if l == "" {
		log.Fatalf("Error: unsupported language: %v", language)
	}

	return l
}

var languages = []string{"javascript", "go"}

func hasLanguageDelimiter(line string) bool {
	if !strings.Contains(line, "```") {
		return false
	}

	for _, lang := range languages {
		if strings.Contains(line, lang) {
			return true
		}
	}
	return false
}

func getSnippetsFromMarkdownFile(file *File) ([]*Snippet, error) {
	lines := strings.Split(file.Contents, "\n")

	snippetLinesList := []SnippetLines{}
	start := 0
	end := 0
	counting := false
	language := ""

	for index, line := range lines {
		lineNumber := index + 1
		if hasLanguageDelimiter(line) {
			language = getLanguage(line)
			start = lineNumber
			counting = true
		}

		if counting && line == "```" {
			end = lineNumber
			snippetLinesList = append(snippetLinesList, SnippetLines{Start: start, End: end, Language: language})
			counting = false
			language = ""
		}
	}

	snippets := []*Snippet{}
	for _, snip := range snippetLinesList {
		content := strings.Join(lines[snip.Start:snip.End-1], "\n")
		content = removeInitialNewline(content)
		content = replaceTabsWithSpaces(content)

		snippet := &Snippet{
			Content:  content,
			Filename: file.Name,
			Language: snip.Language,
			Start:    snip.Start,
			End:      snip.End,
		}
		snippets = append(snippets, snippet)
	}

	return snippets, nil
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
	mediumPath := fmt.Sprintf("%s.medium.md", file.Name)

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
