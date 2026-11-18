# function-comment-linter

`funccomment` is a module plugin for `golangci-lint` that requires every Go
function and method to have a doc comment whose first line starts with the
function or method name and contains a Chinese description. `init` functions
are ignored.

## Requirements

- Go 1.25.4 or newer
- `golangci-lint` 2.6.2 (the version in `.custom-gcl.yml`)

Use the same Go version for the build and lint runs. If the local Go cache was
created by another Go version, use a fresh writable cache, for example:

```bash
GOCACHE=/private/tmp/function-comment-linter-gocache golangci-lint custom -v
```

The command writes the custom binary as `./custom-gcl`; it is ignored by Git.

## Build and run

Build the custom binary:

```bash
golangci-lint custom -v
```

Copy `.golangci.example.yml` to `.golangci.yml`, then run:

```bash
./custom-gcl run --config .golangci.yml ./...
```

The `custom` entry in the configuration is required so the generated binary
knows how to register the module plugin.

## Examples

Valid comments include:

```go
// ProcessData 加载数据。
func ProcessData() {}

/*
processData 加载数据。
*/
func processData() {}
```

These are reported:

```go
func processData() {}

// ProcessData loads the data.
func ProcessData() {}
```

The description must contain at least one Chinese character. A comment that
only contains the function name or only English text is reported.

Run the unit and analyzer tests with:

```bash
GOCACHE=/private/tmp/function-comment-linter-gocache go test ./...
```
