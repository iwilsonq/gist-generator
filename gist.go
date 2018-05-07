package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

func createGist(client *github.Client, snippet *Snippet, index int) (*github.Gist, error) {
	ctx := context.Background()
	description := "A Go script generated gist."
	public := true

	files := make(map[github.GistFilename]github.GistFile)
	gf := github.GistFile{
		Content: &snippet.Content,
	}

	// name vs unnamed type conversion
	filename := fmt.Sprintf("%v-%v.%v", snippet.Filename, index, snippet.Language)
	files[github.GistFilename(filename)] = gf

	gist := github.Gist{
		Description: &description,
		Public:      &public,
		Files:       files,
	}

	g, _, err := client.Gists.Create(ctx, &gist)
	if err != nil {
		log.Fatal(err)
	}

	return g, nil
}

func listGistsByUser(client *github.Client, user string) {
	ctx := context.Background()
	opt := &github.GistListOptions{}
	gists, _, err := client.Gists.List(ctx, user, opt)
	if err != nil {
		log.Fatal(err)
	}
	printGists(gists)
}

func printGists(gists []*github.Gist) {
	for index, gist := range gists {
		fmt.Printf("%v: %v\n", index, gist.Files)
	}
}
