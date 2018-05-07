package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"
	"testing"
)

var markdownTestString = "What is a good way to go about testing my /graphql endpoint for all of my different queries and mutations?\n\n" +
	"```javascript\n" +
	"function sayHello() {\n" +
	"	return function(name) {\n" +
	"		console.log('Hello ', name)\n" +
	"	}\n" +
	"}\n" +
	"```\n\n" +
	"Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTPS, or can be used to authenticate to the API over Basic Authentication."

func createFile(filename string) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(markdownTestString)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func TestMain(m *testing.M) {
	// f, err := createFile("test.md")
	// if err != nil {
	// 	fmt.Printf("could not create file: %v", err)
	// }

	build := exec.Command("go build")
	err = build.Run()
	if err != nil {
		fmt.Printf("could not make binary: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCliArgs(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(path.Join(dir, binaryName), tt.args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if *update {
		writeFixture(t, tt.fixture, output)
	}

	actual := string(output)

	expected := loadFixture(t, tt.fixture)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual = %s, expected = %s", actual, expected)
	}
}
