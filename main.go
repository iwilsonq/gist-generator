package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Snippet - code snippet with references to filename and index
type Snippet struct {
	Language string
	Filename string
	Content  string
	Start    int
	End      int
}

// Languages maps full programming languages to their file extensions
var Languages = map[string]string{
	"javascript": "js",
	"go":         "go",
}

// GistSnippet composed of a Gist and a Snippet
type GistSnippet struct {
	Gist    *github.Gist
	Snippet *Snippet
}

// File type with references to file name and contents
type File struct {
	Name     string
	Contents string
}

func main() {
	path := flag.String("f", "", "specify the path to the markdown file")
	accessToken := flag.String("token", "", "the personal access token from your Github account")
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	filename := strings.Split(*path, ".")[0]
	fileString, err := getFileString(*path)
	if err != nil {
		panic(err)
	}
	file := &File{Name: filename, Contents: fileString}

	snippets, err := getSnippetsFromMarkdownFile(file)
	if err != nil {
		panic(err)
	}

	gistSnippets := []*GistSnippet{}
	for index, snippet := range snippets {
		gist, err := createGist(client, snippet, index)
		if err != nil {
			log.Fatal(err)
		}

		gs := &GistSnippet{
			Gist:    gist,
			Snippet: snippet,
		}

		gistSnippets = append(gistSnippets, gs)
	}

	replaceSnippetsWithURLs(file, gistSnippets)
}
