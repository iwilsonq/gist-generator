# Gist Generator

This Go script which can parse markdown files for their code snippets and create gists for them. It's main dependency is the [go-github](https://github.com/google/go-github) package.

It can write to a new markdown file with snippets replaced by said gist links, making it ideal for preparing a technical article on Medium.

Someday, I hope to have this project well-covered with tests, but for now it's just a tech demo.

## Supported Languages

Though it is rather trivial to add support for other languages, this project currently supports:

*   JavaScript
*   Go
