package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"reflect"

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
}

// GistSnippet composed of a Gist and a Snippet
type GistSnippet struct {
	Gist    *github.Gist
	Snippet *Snippet
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "bfd47d0ac12be422fc15b6c3d2e0872af2745ddb"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	path := "react-testing.md"

	snippets, err := getSnippetsFromMarkdownFile(path, "javascript")
	if err != nil {
		log.Fatal(err)
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

	replaceSnippetsWithURLs(path, gistSnippets)
}

func stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
