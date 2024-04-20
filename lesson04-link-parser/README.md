# Link Parser

My implementation of [exercise 4](https://github.com/gophercises/link) on [Gophercises](https://gophercises.com/).

Parses a HTML document from file for any links and grabs the text from them.

## How to run

```shell
go run cmd/link-parser/main.go --file-path <path-to-html-file>

# Sample files are in ./testdata/
go run cmd/link-parser/main.go --file-path ./testdata/ex1.html
```