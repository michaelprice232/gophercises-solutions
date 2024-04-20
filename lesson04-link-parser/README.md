# Link Parser

My implementation of [exercise 4](https://github.com/gophercises/link) on [Gophercises](https://gophercises.com/).

Parses an HTML document from file for any links and grabs the text from them.

I found this to be a tough exercise and my solution does not work for all the cases that the instructors solution works for (see commented out unit test case), but I have kept this code as-is, to serve as a learning record.  

## How to run

```shell
go run cmd/link-parser/main.go --file-path <path-to-html-file>

# Sample files are in ./testdata/
go run cmd/link-parser/main.go --file-path ./testdata/ex1.html
```